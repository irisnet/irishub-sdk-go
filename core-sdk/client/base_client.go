// Package modules is to warpped the API provided by each module of IRITA
//
package client

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/gogo/protobuf/proto"
	"github.com/irisnet/irishub-sdk-go/common"
	commoncache "github.com/irisnet/irishub-sdk-go/common/cache"
	commoncodec "github.com/irisnet/irishub-sdk-go/common/codec"
	sdklog "github.com/irisnet/irishub-sdk-go/common/log"
	sdktypes "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/tx"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"strings"
	"time"
)

const (
	concurrency       = 16
	cacheCapacity     = 100
	cacheExpirePeriod = 1 * time.Minute
	tryThreshold      = 3
	maxBatch          = 100
)

type baseClient struct {
	sdktypes.TmClient
	sdktypes.GRPCClient
	sdktypes.KeyManager
	sdktypes.TokenManager
	logger         log.Logger
	cfg            *sdktypes.ClientConfig
	encodingConfig sdktypes.EncodingConfig
	l              *locker

	accountQuery
}

// NewBaseClient return the baseClient for every sub modules
func NewBaseClient(cfg sdktypes.ClientConfig, encodingConfig sdktypes.EncodingConfig, logger log.Logger) sdktypes.BaseClient {
	// create logger
	if logger == nil {
		logger = sdklog.NewLogger(sdklog.Config{
			Format: sdklog.FormatText,
			Level:  cfg.Level,
		})
	}
	base := baseClient{
		TmClient:       NewRPCClient(cfg.NodeURI, encodingConfig.Amino, encodingConfig.TxConfig.TxDecoder(), logger, cfg.Timeout),
		GRPCClient:     NewGRPCClient(cfg.GRPCAddr),
		logger:         logger,
		cfg:            &cfg,
		encodingConfig: encodingConfig,
		l:              NewLocker(concurrency).setLogger(logger),
	}
	base.TokenManager = cfg.TokenManager
	base.KeyManager = keyManager{
		keyDAO: cfg.KeyDAO,
		algo:   cfg.Algo,
	}

	c := commoncache.NewCache(cacheCapacity, cfg.Cached)
	base.accountQuery = accountQuery{
		Queries:    base,
		GRPCClient: base.GRPCClient,
		Logger:     base.Logger(),
		Cache:      c,
		cdc:        encodingConfig.Marshaler,
		km:         base.KeyManager,
		expiration: cacheExpirePeriod,
	}
	return &base
}

func (base *baseClient) Logger() log.Logger {
	return base.logger
}

func (base *baseClient) SetLogger(logger log.Logger) {
	base.logger = logger
}

// Codec returns codec.
func (base *baseClient) Marshaler() commoncodec.Marshaler {
	return base.encodingConfig.Marshaler
}

func (base *baseClient) BuildTxHash(msg []sdktypes.Msg, baseTx sdktypes.BaseTx) (string, sdktypes.Error) {
	txByte, _, err := base.buildTx(msg, baseTx)
	if err != nil {
		return "", sdktypes.Wrap(err)
	}
	return strings.ToUpper(hex.EncodeToString(tmhash.Sum(txByte))), nil
}

func (base *baseClient) BuildAndSign(msg []sdktypes.Msg, baseTx sdktypes.BaseTx) ([]byte, sdktypes.Error) {
	builder, err := base.prepare(baseTx)
	if err != nil {
		return nil, sdktypes.Wrap(err)
	}

	txByte, err := builder.BuildAndSign(baseTx.From, msg, true)
	if err != nil {
		return nil, sdktypes.Wrap(err)
	}

	base.Logger().Debug("sign transaction success")
	return txByte, nil
}

func (base *baseClient) BuildAndSendWithAccount(addr string, accountNumber, sequence uint64, msg []sdktypes.Msg, baseTx sdktypes.BaseTx) (sdktypes.ResultTx, sdktypes.Error) {
	txByte, ctx, err := base.buildTxWithAccount(addr, accountNumber, sequence, msg, baseTx)
	if err != nil {
		return sdktypes.ResultTx{}, err
	}

	valid, err := base.ValidateTxSize(len(txByte), msg)
	if err != nil {
		return sdktypes.ResultTx{}, err
	}
	if !valid {
		base.Logger().Debug("tx is too large")
		// filter out transactions that have been sent
		// reset the maximum number of msg in each transaction
		//batch = batch / 2
		return sdktypes.ResultTx{}, sdktypes.GetError(sdktypes.RootCodespace, uint32(sdktypes.TxTooLarge))
	}
	return base.broadcastTx(txByte, ctx.Mode())
}

func (base *baseClient) BuildAndSend(msg []sdktypes.Msg, baseTx sdktypes.BaseTx) (sdktypes.ResultTx, sdktypes.Error) {
	var res sdktypes.ResultTx
	var address string

	// lock the account
	base.l.Lock(baseTx.From)
	defer base.l.Unlock(baseTx.From)

	retryableFunc := func() error {
		txByte, ctx, e := base.buildTx(msg, baseTx)
		if e != nil {
			return e
		}

		if res, e = base.broadcastTx(txByte, ctx.Mode()); e != nil {
			address = ctx.Address()
			return e
		}
		return nil
	}

	retryIfFunc := func(err error) bool {
		e, ok := err.(sdktypes.Error)
		if ok && sdktypes.Code(e.Code()) == sdktypes.WrongSequence {
			return true
		}
		return false
	}

	onRetryFunc := func(n uint, err error) {
		_ = base.removeCache(address)
		base.Logger().Error("wrong sequence, will retry",
			"address", address, "attempts", n, "err", err.Error())
	}

	err := retry.Do(retryableFunc,
		retry.Attempts(tryThreshold),
		retry.RetryIf(retryIfFunc),
		retry.OnRetry(onRetryFunc),
	)

	if err != nil {
		return res, sdktypes.Wrap(err)
	}
	return res, nil
}

func (base *baseClient) SendBatch(msgs sdktypes.Msgs, baseTx sdktypes.BaseTx) (rs []sdktypes.ResultTx, err sdktypes.Error) {
	if msgs == nil || len(msgs) == 0 {
		return rs, sdktypes.Wrapf("must have at least one message in list")
	}

	defer sdktypes.CatchPanic(func(errMsg string) {
		base.Logger().Error("broadcast msg failed", "errMsg", errMsg)
	})
	// validate msg
	for _, m := range msgs {
		if err := m.ValidateBasic(); err != nil {
			return rs, sdktypes.Wrap(err)
		}
	}
	base.Logger().Debug("validate msg success")

	// lock the account
	base.l.Lock(baseTx.From)
	defer base.l.Unlock(baseTx.From)

	var address string
	var batch = maxBatch

	retryableFunc := func() error {
		for i, ms := range common.SubArray(batch, msgs) {
			mss := ms.(sdktypes.Msgs)
			txByte, ctx, err := base.buildTx(mss, baseTx)
			if err != nil {
				return err
			}

			valid, err := base.ValidateTxSize(len(txByte), mss)
			if err != nil {
				return err
			}
			if !valid {
				base.Logger().Debug("tx is too large", "msgsLength", batch)
				// filter out transactions that have been sent
				msgs = msgs[i*batch:]
				// reset the maximum number of msg in each transaction
				batch = batch / 2
				return sdktypes.GetError(sdktypes.RootCodespace, uint32(sdktypes.TxTooLarge))
			}
			res, err := base.broadcastTx(txByte, ctx.Mode())
			if err != nil {
				address = ctx.Address()
				return err
			}
			rs = append(rs, res)
		}
		return nil
	}

	retryIf := func(err error) bool {
		e, ok := err.(sdktypes.Error)
		if ok && (sdktypes.Code(e.Code()) == sdktypes.InvalidSequence || sdktypes.Code(e.Code()) == sdktypes.TxTooLarge) {
			return true
		}
		return false
	}

	onRetry := func(n uint, err error) {
		_ = base.removeCache(address)
		base.Logger().Error("wrong sequence, will retry",
			"address", address, "attempts", n, "err", err.Error())
	}

	_ = retry.Do(retryableFunc,
		retry.Attempts(tryThreshold),
		retry.RetryIf(retryIf),
		retry.OnRetry(onRetry),
	)
	return rs, nil
}

func (base baseClient) QueryWithResponse(path string, data interface{}, result sdktypes.Response) error {
	res, err := base.Query(path, data)
	if err != nil {
		return err
	}

	if err := base.encodingConfig.Marshaler.UnmarshalJSON(res, result.(proto.Message)); err != nil {
		return err
	}

	return nil
}

func (base baseClient) Query(path string, data interface{}) ([]byte, error) {
	var bz []byte
	var err error
	if data != nil {
		bz, err = base.encodingConfig.Marshaler.MarshalJSON(data.(proto.Message))
		if err != nil {
			return nil, err
		}
	}

	opts := rpcclient.ABCIQueryOptions{
		// Height: cliCtx.Height,
		Prove: false,
	}
	result, err := base.ABCIQueryWithOptions(context.Background(), path, bz, opts)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	return resp.Value, nil
}

func (base baseClient) QueryStore(key sdktypes.HexBytes, storeName string, height int64, prove bool) (res abci.ResponseQuery, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "key")
	opts := rpcclient.ABCIQueryOptions{
		Prove:  prove,
		Height: height,
	}

	result, err := base.ABCIQueryWithOptions(context.Background(), path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.New(resp.Log)
	}
	return resp, nil
}

func (base *baseClient) prepare(baseTx sdktypes.BaseTx) (*sdktypes.Factory, error) {
	factory := sdktypes.NewFactory().
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.KeyManager).
		WithMode(base.cfg.Mode).
		WithSimulateAndExecute(baseTx.SimulateAndExecute).
		WithGas(base.cfg.Gas).
		WithGasAdjustment(base.cfg.GasAdjustment).
		WithSignModeHandler(tx.MakeSignModeHandler(tx.DefaultSignModes)).
		WithTxConfig(base.encodingConfig.TxConfig).
		WithQueryFunc(base.QueryWithData)

	addr, err := base.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	factory.WithAddress(addr.String())

	account, err := base.QueryAndRefreshAccount(addr.String())
	if err != nil {
		return nil, err
	}
	factory.WithAccountNumber(account.AccountNumber).
		WithSequence(account.Sequence).
		WithPassword(baseTx.Password)

	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := base.TokenManager.ToMinCoin(baseTx.Fee...)
		if err != nil {
			return nil, err
		}
		factory.WithFee(fees)
	} else {
		fees, err := base.TokenManager.ToMinCoin(base.cfg.Fee...)
		if err != nil {
			panic(err)
		}
		factory.WithFee(fees)
	}

	if len(baseTx.Mode) > 0 {
		factory.WithMode(baseTx.Mode)
	}

	if baseTx.Gas > 0 {
		factory.WithGas(baseTx.Gas)
	}

	if baseTx.GasAdjustment > 0 {
		factory.WithGasAdjustment(baseTx.GasAdjustment)
	}

	if len(baseTx.Memo) > 0 {
		factory.WithMemo(baseTx.Memo)
	}
	return factory, nil
}

// TODO
func (base *baseClient) prepareTemp(addr string, accountNumber, sequence uint64, baseTx sdktypes.BaseTx) (*sdktypes.Factory, error) {
	factory := sdktypes.NewFactory().
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.KeyManager).
		WithMode(base.cfg.Mode).
		WithSimulateAndExecute(baseTx.SimulateAndExecute).
		WithGas(base.cfg.Gas).
		WithSignModeHandler(tx.MakeSignModeHandler(tx.DefaultSignModes)).
		WithTxConfig(base.encodingConfig.TxConfig)

	factory.WithAddress(addr).
		WithAccountNumber(accountNumber).
		WithSequence(sequence).
		WithPassword(baseTx.Password)

	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := base.TokenManager.ToMinCoin(baseTx.Fee...)
		if err != nil {
			return nil, err
		}
		factory.WithFee(fees)
	} else {
		fees, err := base.TokenManager.ToMinCoin(base.cfg.Fee...)
		if err != nil {
			panic(err)
		}
		factory.WithFee(fees)
	}

	if len(baseTx.Mode) > 0 {
		factory.WithMode(baseTx.Mode)
	}

	if baseTx.Gas > 0 {
		factory.WithGas(baseTx.Gas)
	}

	if len(baseTx.Memo) > 0 {
		factory.WithMemo(baseTx.Memo)
	}
	return factory, nil
}

func (base *baseClient) ValidateTxSize(txSize int, msgs []sdktypes.Msg) (bool, sdktypes.Error) {
	if uint64(txSize) > base.cfg.TxSizeLimit {
		return false, nil
	}
	return true, nil
}

type locker struct {
	shards []chan int
	size   int
	logger log.Logger
}

//NewLocker implement the function of lock, can lock resources according to conditions
func NewLocker(size int) *locker {
	shards := make([]chan int, size)
	for i := 0; i < size; i++ {
		shards[i] = make(chan int, 1)
	}
	return &locker{
		shards: shards,
		size:   size,
	}
}

func (l *locker) setLogger(logger log.Logger) *locker {
	l.logger = logger
	return l
}

func (l *locker) Lock(key string) {
	ch := l.getShard(key)
	ch <- 1
}

func (l *locker) Unlock(key string) {
	ch := l.getShard(key)
	<-ch
}

func (l *locker) getShard(key string) chan int {
	index := uint(l.indexFor(key)) % uint(l.size)
	return l.shards[index]
}

func (l *locker) indexFor(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
