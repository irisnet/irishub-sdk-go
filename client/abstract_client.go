package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/irisnet/irishub-sdk-go/adapter"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type abstractClient struct {
	sdk.TmClient
	sdk.KeyManager

	logger *log.Logger
	cfg    sdk.SDKConfig
	cdc    sdk.Codec
}

func createAbstractClient(cdc sdk.Codec, cfg sdk.SDKConfig, logger *log.Logger) *abstractClient {
	ac := abstractClient{
		KeyManager: adapter.NewDAOAdapter(cfg.KeyDAO, cfg.StoreType),
		TmClient:   NewRPCClient(cfg.NodeURI, cdc, logger),
		logger:     logger,
		cfg:        cfg,
		cdc:        cdc,
	}

	ac.init()
	return &ac
}

func (ac *abstractClient) init() {
	fees, err := ac.ToMinCoin(ac.cfg.Fee...)
	if err != nil {
		panic(err)
	}
	ac.cfg.Fee = sdk.NewDecCoinsFromCoins(fees...)
}

func (ac abstractClient) Logger() *log.Logger {
	return ac.logger
}

func (ac *abstractClient) BuildAndSend(msg []sdk.Msg, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	defer sdk.CatchPanic(func(errMsg string) {
		ac.Logger().Error().
			Msgf("broadcast msg failed:%s", errMsg)
	})
	//validate msg
	for _, m := range msg {
		if err := m.ValidateBasic(); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}
	ac.Logger().Info().Msg("validate msg success")

	ctx, err := ac.prepare(baseTx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	tx, err := ctx.BuildAndSign(baseTx.From, msg)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	ac.Logger().Info().
		Strs("data", tx.GetSignBytes()).
		Msg("sign transaction success")

	txByte, err := ac.cdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	res, err := ac.broadcastTx(txByte, ctx.Mode)
	if err != nil {
		ac.Logger().Err(err).Msg("broadcastTx transaction failed")
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	ac.Logger().Info().
		Str("txHash", res.Hash).
		Msg("broadcastTx transaction success")
	return res, nil
}

func (ac *abstractClient) SendMsgBatch(batch int, msgs []sdk.Msg, baseTx sdk.BaseTx) (rs []sdk.ResultTx, err sdk.Error) {
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
		res, err := ac.BuildAndSend(ms, baseTx)
		if err != nil {
			return rs, sdk.Wrap(err)
		}
		rs = append(rs, res)
	}
	return rs, nil
}

func (ac abstractClient) Broadcast(signedTx sdk.StdTx, mode sdk.BroadcastMode) (sdk.ResultTx, sdk.Error) {
	txByte, err := ac.cdc.MarshalBinaryLengthPrefixed(signedTx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return ac.broadcastTx(txByte, mode)
}

func (ac abstractClient) QueryWithResponse(path string, data interface{}, result sdk.Response) error {
	res, err := ac.Query(path, data)
	if err != nil {
		return err
	}

	if err := ac.cdc.UnmarshalJSON(res, result); err != nil {
		return err
	}

	return nil
}

func (ac abstractClient) Query(path string, data interface{}) ([]byte, error) {
	var bz []byte
	var err error
	if data != nil {
		bz, err = ac.cdc.MarshalJSON(data)
		if err != nil {
			return nil, err
		}
	}

	opts := rpcclient.ABCIQueryOptions{
		//Height: cliCtx.Height,
		Prove: false,
	}
	result, err := ac.ABCIQueryWithOptions(path, bz, opts)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	return resp.Value, nil
}

func (ac abstractClient) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "subspace")
	opts := rpcclient.ABCIQueryOptions{
		//Height: cliCtx.Height,
		Prove: false,
	}

	result, err := ac.TmClient.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.New(resp.Log)
	}
	return resp.Value, nil
}

func (ac abstractClient) QueryAccount(address string) (sdk.BaseAccount, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.BaseAccount{}, err
	}

	param := struct {
		Address sdk.AccAddress
	}{
		Address: addr,
	}

	var account sdk.BaseAccount
	if err := ac.QueryWithResponse("custom/acc/account", param, &account); err != nil {
		return sdk.BaseAccount{}, err
	}
	return account, nil
}

func (ac abstractClient) QueryToken(symbol string) (sdk.Token, error) {
	if token, ok := sdk.GetToken(symbol); ok {
		return token, nil
	}
	param := struct {
		Symbol string
	}{
		Symbol: symbol,
	}

	//TODO will remove from v1.0
	symbol = strings.TrimSuffix(symbol, "-min")
	var token sdk.Token
	if err := ac.QueryWithResponse("custom/asset/token", param, &token); err != nil {
		return sdk.Token{}, err
	}
	sdk.CacheTokens(token)
	return token, nil
}

func (ac abstractClient) ToMinCoin(coins ...sdk.DecCoin) (dstCoins sdk.Coins, err sdk.Error) {
	for _, coin := range coins {
		token, err := ac.QueryToken(coin.Denom)
		if err != nil {
			return nil, sdk.Wrap(err)
		}

		minCoin, err := token.GetCoinType().ConvertToMinCoin(coin)
		if err != nil {
			return nil, sdk.Wrap(err)
		}
		dstCoins = append(dstCoins, minCoin)
	}
	return dstCoins.Sort(), nil
}

func (ac abstractClient) ToMainCoin(coins ...sdk.Coin) (dstCoins sdk.DecCoins, err sdk.Error) {
	for _, coin := range coins {
		token, err := ac.QueryToken(coin.Denom)
		if err != nil {
			return dstCoins, sdk.Wrap(err)
		}

		mainCoin, err := token.GetCoinType().ConvertToMainCoin(coin)
		if err != nil {
			return dstCoins, sdk.Wrap(err)
		}
		dstCoins = append(dstCoins, mainCoin)
	}
	return dstCoins.Sort(), nil
}

func (ac abstractClient) QueryAddress(name string) (sdk.AccAddress, error) {
	return ac.KeyManager.Query(name)
}

// QueryTx returns the tx info
func (ac abstractClient) QueryTx(hash string) (sdk.ResultQueryTx, error) {
	tx, err := hex.DecodeString(hash)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	res, err := ac.Tx(tx, true)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	resBlocks, err := ac.getBlocksForTxResults([]*ctypes.ResultTx{res})
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}
	return ac.formatTxResult(res, resBlocks[res.Height])
}

func (ac abstractClient) QueryTxs(builder *sdk.EventQueryBuilder, page, size int) (sdk.ResultSearchTxs, error) {

	query := builder.Build()
	if len(query) == 0 {
		return sdk.ResultSearchTxs{}, errors.New("must declare at least one tag to search")
	}

	res, err := ac.TxSearch(query, true, page, size)
	if err != nil {
		return sdk.ResultSearchTxs{}, err
	}

	resBlocks, err := ac.getBlocksForTxResults(res.Txs)
	if err != nil {
		return sdk.ResultSearchTxs{}, err
	}

	var txs []sdk.ResultQueryTx
	for i, tx := range res.Txs {
		txInfo, err := ac.formatTxResult(tx, resBlocks[res.Txs[i].Height])
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

func (ac *abstractClient) prepare(baseTx sdk.BaseTx) (*sdk.TxContext, error) {
	fees, _ := ac.cfg.Fee.TruncateDecimal()
	ctx := &sdk.TxContext{
		Codec:      ac.cdc,
		ChainID:    ac.cfg.ChainID,
		KeyManager: ac.KeyManager,
		Network:    ac.cfg.Network,
		Fee:        fees,
		Mode:       ac.cfg.Mode,
		Simulate:   false,
		Gas:        ac.cfg.Gas,
	}

	addr, err := ac.QueryAddress(baseTx.From)
	if err != nil {
		return nil, err
	}

	account, err := ac.QueryAccount(addr.String())
	if err != nil {
		return nil, err
	}
	ctx.WithAccountNumber(account.AccountNumber).
		WithSequence(account.Sequence).
		WithPassword(baseTx.Password)

	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := ac.ToMinCoin(baseTx.Fee...)
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

func (ac abstractClient) broadcastTx(txBytes []byte, mode sdk.BroadcastMode) (sdk.ResultTx, sdk.Error) {
	switch mode {
	case sdk.Commit:
		return ac.broadcastTxCommit(txBytes)
	case sdk.Async:
		return ac.broadcastTxAsync(txBytes)
	case sdk.Sync:
		return ac.broadcastTxSync(txBytes)

	}
	return sdk.ResultTx{}, sdk.Wrapf("commit mode(%s) not supported", ac.cfg.Mode)
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (ac abstractClient) broadcastTxCommit(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := ac.BroadcastTxCommit(tx)
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
func (ac abstractClient) broadcastTxSync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := ac.BroadcastTxSync(tx)
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
func (ac abstractClient) broadcastTxAsync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := ac.BroadcastTxAsync(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return sdk.ResultTx{
		Hash: res.Hash.String(),
	}, nil
}

func (ac abstractClient) getBlocksForTxResults(resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	resBlocks := make(map[int64]*ctypes.ResultBlock)
	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := ac.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}
	return resBlocks, nil
}

func (ac abstractClient) formatTxResult(res *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (sdk.ResultQueryTx, error) {

	var tx sdk.StdTx
	err := ac.cdc.UnmarshalBinaryLengthPrefixed(res.Tx, &tx)
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
