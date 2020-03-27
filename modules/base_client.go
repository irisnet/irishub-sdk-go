package modules

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/irishub-sdk-go/adapter"
	"github.com/irisnet/irishub-sdk-go/tools/cache"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type baseClient struct {
	sdk.TmClient
	sdk.KeyManager
	localAccount
	localToken

	logger *log.Logger
	cfg    sdk.SDKConfig
	cdc    sdk.Codec

	l *locker
}

func NewBaseClient(cdc sdk.Codec, cfg sdk.SDKConfig, logger *log.Logger) *baseClient {
	base := baseClient{
		KeyManager: adapter.NewDAOAdapter(cfg.KeyDAO, cfg.StoreType),
		TmClient:   NewRPCClient(cfg.NodeURI, cdc, logger),
		logger:     logger,
		cfg:        cfg,
		cdc:        cdc,
		l:          NewLocker(16),
	}

	c := cache.NewLRU(100)
	base.localAccount = localAccount{
		Query:      base,
		Logger:     base.Logger(),
		Cache:      c,
		keyManager: base.KeyManager,
		expiration: 1 * time.Minute,
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
}

func (base *baseClient) Logger() *log.Logger {
	return base.logger
}

func (base *baseClient) BuildAndSend(msg []sdk.Msg, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	defer sdk.CatchPanic(func(errMsg string) {
		base.Logger().Error().
			Msgf("broadcast msg failed:%s", errMsg)
	})
	//validate msg
	for _, m := range msg {
		if err := m.ValidateBasic(); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}
	base.Logger().Info().Msg("validate msg success")

	//lock the account
	base.l.Lock(baseTx.From)
	defer base.l.Unlock(baseTx.From)

retry:
	ctx, err := base.prepare(baseTx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	tx, err := ctx.BuildAndSign(baseTx.From, msg)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	base.Logger().Info().
		Strs("data", tx.GetSignBytes()).
		Msg("sign transaction success")

	txByte, err := base.cdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	res, e := base.broadcastTx(txByte, ctx.Mode())
	if e != nil {
		if sdk.Code(e.Code()) == sdk.InvalidSequence &&
			base.enabled {
			_, _ = base.Refresh(ctx.Address())
			goto retry
		}
		base.Logger().Err(e).Msg("broadcastTx transaction failed")
		return sdk.ResultTx{}, sdk.Wrap(e)
	}
	base.Logger().Info().
		Str("txHash", res.Hash).
		Msg("broadcastTx transaction success")

	return res, nil
}

func (base *baseClient) SendMsgBatch(batch int, msgs []sdk.Msg, baseTx sdk.BaseTx) (rs []sdk.ResultTx, err sdk.Error) {
	splitMsgs := func(batch int, msgs []sdk.Msg) (segments [][]sdk.Msg) {
		max := len(msgs)
		if max < batch {
			return [][]sdk.Msg{msgs}
		}

		quantity := max / batch
		for i := 1; i <= batch; i++ {
			start := (i - 1) * quantity
			end := i * quantity
			if i != batch {
				segments = append(segments, msgs[start:end])
			} else {
				segments = append(segments, msgs[start:])
			}
		}
		return segments
	}

	if msgs == nil || len(msgs) == 0 {
		return rs, sdk.Wrapf("must have at least one message in list")
	}

	baseTx.Mode = sdk.Commit
	for _, ms := range splitMsgs(batch, msgs) {
		res, err := base.BuildAndSend(ms, baseTx)
		if err != nil {
			return rs, sdk.Wrap(err)
		}
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
	path := fmt.Sprintf("/store/%s/%s", storeName, "subspace")
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

// QueryTx returns the tx info
func (base baseClient) QueryTx(hash string) (sdk.ResultQueryTx, error) {
	tx, err := hex.DecodeString(hash)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	res, err := base.Tx(tx, true)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	resBlocks, err := base.getResultBlocks([]*ctypes.ResultTx{res})
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}
	return base.parseTxResult(res, resBlocks[res.Height])
}

func (base baseClient) QueryTxs(builder *sdk.EventQueryBuilder, page, size int) (sdk.ResultSearchTxs, error) {

	query := builder.Build()
	if len(query) == 0 {
		return sdk.ResultSearchTxs{}, errors.New("must declare at least one tag to search")
	}

	res, err := base.TxSearch(query, true, page, size)
	if err != nil {
		return sdk.ResultSearchTxs{}, err
	}

	resBlocks, err := base.getResultBlocks(res.Txs)
	if err != nil {
		return sdk.ResultSearchTxs{}, err
	}

	var txs []sdk.ResultQueryTx
	for i, tx := range res.Txs {
		txInfo, err := base.parseTxResult(tx, resBlocks[res.Txs[i].Height])
		if err != nil {
			return sdk.ResultSearchTxs{}, err
		}
		txs = append(txs, txInfo)
	}

	return sdk.ResultSearchTxs{
		Total: res.TotalCount,
		Txs:   txs,
	}, nil
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

func (base baseClient) broadcastTx(txBytes []byte, mode sdk.BroadcastMode) (sdk.ResultTx, sdk.Error) {
	switch mode {
	case sdk.Commit:
		return base.broadcastTxCommit(txBytes)
	case sdk.Async:
		return base.broadcastTxAsync(txBytes)
	case sdk.Sync:
		return base.broadcastTxSync(txBytes)

	}
	return sdk.ResultTx{}, sdk.Wrapf("commit mode(%s) not supported", base.cfg.Mode)
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (base baseClient) broadcastTxCommit(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxCommit(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if !res.CheckTx.IsOK() {
		return sdk.ResultTx{}, sdk.GetError(res.CheckTx.Codespace,
			res.CheckTx.Code, res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return sdk.ResultTx{}, sdk.GetError(res.DeliverTx.Codespace,
			res.DeliverTx.Code, res.DeliverTx.Log)
	}

	return sdk.ResultTx{
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Tags:      sdk.ParseTags(res.DeliverTx.Tags),
		Hash:      res.Hash.String(),
		Height:    res.Height,
	}, nil
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously.
func (base baseClient) broadcastTxSync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxSync(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if res.Code != 0 {
		return sdk.ResultTx{}, sdk.GetError(sdk.RootCodespace,
			res.Code, res.Log)
	}

	return sdk.ResultTx{
		Hash: res.Hash.String(),
	}, nil
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (base baseClient) broadcastTxAsync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxAsync(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return sdk.ResultTx{
		Hash: res.Hash.String(),
	}, nil
}

func (base baseClient) getResultBlocks(resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	resBlocks := make(map[int64]*ctypes.ResultBlock)
	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := base.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}
	return resBlocks, nil
}

func (base baseClient) parseTxResult(res *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (sdk.ResultQueryTx, error) {

	var tx sdk.StdTx
	err := base.cdc.UnmarshalBinaryLengthPrefixed(res.Tx, &tx)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	return sdk.ResultQueryTx{
		Hash:   res.Hash.String(),
		Height: res.Height,
		Tx:     tx,
		Result: sdk.TxResult{
			Code:      res.TxResult.Code,
			Log:       res.TxResult.Log,
			GasWanted: res.TxResult.GasWanted,
			GasUsed:   res.TxResult.GasUsed,
			Tags:      sdk.ParseTags(res.TxResult.Tags),
		},
		Timestamp: resBlock.Block.Time.Format(time.RFC3339),
	}, nil
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
	index := uint(indexFor(key)) % uint(l.size)
	return l.shards[index]
}

func indexFor(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
