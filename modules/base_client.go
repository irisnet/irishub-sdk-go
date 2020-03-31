package modules

import (
	"errors"
	"fmt"
	"time"

	"github.com/irisnet/irishub-sdk-go/utils"

	"github.com/irisnet/irishub-sdk-go/adapter"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/cache"
	"github.com/irisnet/irishub-sdk-go/utils/log"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

const (
	concurrency       = 16
	cacheCapacity     = 100
	cacheExpirePeriod = 1 * time.Minute
	timeout           = 5 * time.Second
	tryThreshold      = 3
	maxMsgsCnt        = 10
)

type baseClient struct {
	sdk.TmClient
	sdk.KeyManager
	localAccount
	localToken

	logger *log.Logger
	cfg    sdk.ClientConfig
	cdc    sdk.Codec

	l *locker
}

func NewBaseClient(cdc sdk.Codec, cfg sdk.ClientConfig, logger *log.Logger) *baseClient {
	base := baseClient{
		KeyManager: adapter.NewDAOAdapter(cfg.KeyDAO, cfg.StoreType),
		TmClient:   NewRPCClient(cfg.NodeURI, cdc, logger),
		logger:     logger,
		cfg:        cfg,
		cdc:        cdc,
		l:          NewLocker(concurrency),
	}

	c := cache.NewLRU(cacheCapacity)
	base.localAccount = localAccount{
		Queries:    base,
		Logger:     base.Logger(),
		Cache:      c,
		keyManager: base.KeyManager,
		expiration: cacheExpirePeriod,
	}

	base.localToken = localToken{
		q:      base,
		Logger: base.Logger(),
		Cache:  c,
	}

	base.init()
	return &base
}

func (base *baseClient) init() {
	fees, err := base.ToMinCoin(base.cfg.Fee...)
	if err != nil {
		panic(err)
	}
	base.cfg.Fee = sdk.NewDecCoinsFromCoins(fees...)
	if base.cfg.Timeout.Nanoseconds() <= 0 {
		base.cfg.Timeout = timeout
	}
}

func (base *baseClient) Logger() *log.Logger {
	return base.logger
}

func (base *baseClient) BuildAndSend(msg []sdk.Msg, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	res, err := base.SendMsgBatch(msg, baseTx)
	if err != nil || len(res) == 0 {
		return sdk.ResultTx{}, sdk.WrapWithMessage(err, "send transaction failed")
	}

	return res[0], nil
}

func (base *baseClient) SendMsgBatch(msgs sdk.Msgs, baseTx sdk.BaseTx) (rs []sdk.ResultTx, err sdk.Error) {
	if msgs == nil || len(msgs) == 0 {
		return rs, sdk.Wrapf("must have at least one message in list")
	}

	defer sdk.CatchPanic(func(errMsg string) {
		base.Logger().Error().
			Msgf("broadcast msg failed:%s", errMsg)
	})
	//validate msg
	for _, m := range msgs {
		if err := m.ValidateBasic(); err != nil {
			return rs, sdk.Wrap(err)
		}
	}
	base.Logger().Info().Msg("validate msg success")

	//lock the account
	base.l.Lock(baseTx.From)
	defer base.l.Unlock(baseTx.From)

	batch := maxMsgsCnt
	var tryCnt = 0

resize:
	for i, ms := range utils.SplitArray(batch, msgs) {
		mss := ms.(sdk.Msgs)

	retry:
		txByte, ctx, err := base.buildTx(mss, baseTx)
		if err != nil {
			return rs, err
		}

		if len(txByte) > 1000 {
			base.Logger().Warn().
				Int("MaxMsgsLen", batch).
				Msg("the transaction content is too large and will be re-sent in batches")

			// filter out transactions that have been sent
			msgs = msgs[i*batch:]
			// reset the maximum number of msg in each transaction
			batch--
			_ = base.RemoveAccount(ctx.Address())
			goto resize
		}

		res, err := base.broadcastTx(txByte, ctx.Mode())
		if err != nil {
			if sdk.Code(err.Code()) == sdk.InvalidSequence {
				base.Logger().Warn().
					Str("address", ctx.Address()).
					Int("tryCnt", tryCnt).
					Msg("account information cached has error,will sync from chain and try to send transaction again")

				if tryCnt++; tryCnt >= tryThreshold {
					return rs, err
				}
				_ = base.RemoveAccount(ctx.Address())
				goto retry
			}
			base.Logger().
				Err(err).
				Msg("broadcastTx transaction failed")
			return rs, err
		}
		base.Logger().Info().
			Str("txHash", res.Hash).
			Msg("broadcastTx transaction success")
		rs = append(rs, res)
	}
	return rs, nil
}

func (base baseClient) Broadcast(signedTx sdk.StdTx, mode sdk.BroadcastMode) (sdk.ResultTx, sdk.Error) {
	txByte, err := base.cdc.MarshalBinaryLengthPrefixed(signedTx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return base.broadcastTx(txByte, mode)
}

func (base baseClient) QueryWithResponse(path string, data interface{}, result sdk.Response) error {
	res, err := base.Query(path, data)
	if err != nil {
		return err
	}

	if err := base.cdc.UnmarshalJSON(res, result); err != nil {
		return err
	}

	return nil
}

func (base baseClient) Query(path string, data interface{}) ([]byte, error) {
	var bz []byte
	var err error
	if data != nil {
		bz, err = base.cdc.MarshalJSON(data)
		if err != nil {
			return nil, err
		}
	}

	opts := rpcclient.ABCIQueryOptions{
		//Height: cliCtx.Height,
		Prove: false,
	}
	result, err := base.ABCIQueryWithOptions(path, bz, opts)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	return resp.Value, nil
}

func (base baseClient) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "key")
	opts := rpcclient.ABCIQueryOptions{
		//Height: cliCtx.Height,
		Prove: false,
	}

	result, err := base.TmClient.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.New(resp.Log)
	}
	return resp.Value, nil
}

func (base *baseClient) prepare(baseTx sdk.BaseTx) (*sdk.TxContext, error) {
	fees, _ := base.cfg.Fee.TruncateDecimal()
	ctx := &sdk.TxContext{}
	ctx.WithCodec(base.cdc).
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.KeyManager).
		WithNetwork(base.cfg.Network).
		WithFee(fees).
		WithMode(base.cfg.Mode).
		WithSimulate(false).
		WithGas(base.cfg.Gas)

	addr, err := base.QueryAddress(baseTx.From)
	if err != nil {
		return nil, err
	}
	ctx.WithAddress(addr.String())

	account, err := base.QueryAndRefreshAccount(addr.String())
	if err != nil {
		return nil, err
	}
	ctx.WithAccountNumber(account.AccountNumber).
		WithSequence(account.Sequence).
		WithPassword(baseTx.Password)

	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := base.ToMinCoin(baseTx.Fee...)
		if err != nil {
			return nil, err
		}
		ctx.WithFee(fees)
	}

	if len(baseTx.Mode) > 0 {
		ctx.WithMode(baseTx.Mode)
	}

	if baseTx.Simulate {
		ctx.WithSimulate(baseTx.Simulate)
	}

	if baseTx.Gas > 0 {
		ctx.WithGas(baseTx.Gas)
	}

	if len(baseTx.Memo) > 0 {
		ctx.WithMemo(baseTx.Memo)
	}
	return ctx, nil
}

type locker struct {
	shards []chan int
	size   int
}

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
