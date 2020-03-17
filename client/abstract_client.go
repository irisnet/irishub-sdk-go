package client

import (
	"errors"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/adapter"

	"github.com/irisnet/irishub-sdk-go/tools/log"

	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type abstractClient struct {
	*sdk.TxContext
	sdk.TmClient
	logger *log.Logger
	cfg    sdk.SDKConfig
	cdc    sdk.Codec
}

func createAbstractClient(cdc sdk.Codec, cfg sdk.SDKConfig, logger *log.Logger) *abstractClient {
	ctx := sdk.TxContext{
		Codec:      cdc,
		ChainID:    cfg.ChainID,
		Online:     cfg.Online,
		KeyManager: adapter.NewDAOAdapter(cfg.KeyDAO, cfg.StoreType),
		Network:    cfg.Network,
	}
	ac := abstractClient{
		TxContext: &ctx,
		TmClient:  NewRPCClient(cfg.NodeURI, cdc, logger),
		logger:    logger,
		cfg:       cfg,
		cdc:       cdc,
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
	ac.reset()
}

func (ac abstractClient) Logger() *log.Logger {
	return ac.logger
}

func (ac *abstractClient) BuildAndSend(msg []sdk.Msg, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	//validate msg
	for _, m := range msg {
		if err := m.ValidateBasic(); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}
	ac.Logger().Info().Msg("validate msg success")

	if err := ac.prepare(baseTx); err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	tx, err := ac.BuildAndSign(baseTx.From, msg)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	ac.Logger().Info().
		Strs("data", tx.GetSignBytes()).
		Msg("sign transaction success")

	txByte, err := ac.Codec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	res, err := ac.broadcastTx(txByte)
	if err != nil {
		ac.Logger().Err(err).Msg("broadcastTx transaction failed")
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	ac.Logger().Info().
		Str("txhash", res.Hash).
		Str("tags", res.Tags.String()).
		Msg("broadcastTx transaction success")
	return res, nil
}

func (ac abstractClient) Broadcast(signedTx sdk.StdTx, mode sdk.BroadcastMode) (sdk.ResultTx, sdk.Error) {
	ac.Mode = mode
	txByte, err := ac.Codec.MarshalBinaryLengthPrefixed(signedTx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return ac.broadcastTx(txByte)
}

func (ac abstractClient) QueryWithResponse(path string, data interface{}, result sdk.Response) error {
	var bz []byte
	var err error
	if data != nil {
		bz, err = ac.Codec.MarshalJSON(data)
		if err != nil {
			return err
		}
	}

	res, err := ac.TmClient.Query(path, bz)
	if err != nil {
		return err
	}

	if err := ac.Codec.UnmarshalJSON(res, result); err != nil {
		return err
	}

	return nil
}

func (ac abstractClient) Query(path string, data interface{}) ([]byte, error) {
	var bz []byte
	var err error
	if data != nil {
		bz, err = ac.Codec.MarshalJSON(data)
		if err != nil {
			return nil, err
		}
	}

	res, err := ac.TmClient.Query(path, bz)
	if err != nil {
		return nil, err
	}
	return res, nil
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

	var token sdk.Token
	if err := ac.QueryWithResponse("custom/asset/token", param, &token); err != nil {
		return sdk.Token{}, err
	}
	sdk.CacheTokens(token)
	return token, nil
}

func (ac abstractClient) ToMinCoin(coins ...sdk.DecCoin) (dstCoins sdk.Coins, err error) {
	for _, coin := range coins {
		token, err := ac.QueryToken(coin.Denom)
		if err != nil {
			return nil, err
		}

		minCoin, err := token.GetCoinType().ConvertToMinCoin(coin)
		if err != nil {
			return nil, err
		}
		dstCoins = append(dstCoins, minCoin)
	}
	return dstCoins.Sort(), nil
}

func (ac abstractClient) ToMainCoin(coins ...sdk.Coin) (dstCoins sdk.DecCoins, err error) {
	for _, coin := range coins {
		token, err := ac.QueryToken(coin.Denom)
		if err != nil {
			return dstCoins, err
		}

		mainCoin, err := token.GetCoinType().ConvertToMainCoin(coin)
		if err != nil {
			return dstCoins, err
		}
		dstCoins = append(dstCoins, sdk.NewDecCoinFromCoin(mainCoin))
	}
	return dstCoins.Sort(), nil
}

func (ac abstractClient) QueryAddress(name string) (sdk.AccAddress, error) {
	return (*ac.TxContext).KeyManager.Query(name)
}

func (ac *abstractClient) prepare(baseTx sdk.BaseTx) error {
	//clear some params
	ac.reset()
	if ac.Online {
		addr, err := ac.QueryAddress(baseTx.From)
		if err != nil {
			return err
		}

		account, err := ac.QueryAccount(addr.String())
		if err != nil {
			return err
		}

		ac.WithAccountNumber(account.AccountNumber).
			WithSequence(account.Sequence)
	}
	ac.WithPassword(baseTx.Password)
	// first use baseTx params
	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := ac.ToMinCoin(baseTx.Fee...)
		if err != nil {
			return err
		}
		ac.WithFee(fees)
	}

	if len(baseTx.Mode) > 0 {
		ac.WithMode(baseTx.Mode)
	}

	if baseTx.Simulate {
		ac.WithSimulate(baseTx.Simulate)
	}

	if baseTx.Gas > 0 {
		ac.WithGas(baseTx.Gas)
	}

	if len(baseTx.Memo) > 0 {
		ac.WithMemo(baseTx.Memo)
	}
	return nil
}

func (ac *abstractClient) reset() {
	fees, _ := ac.cfg.Fee.TruncateDecimal()
	ac.WithAccountNumber(uint64(0)).
		WithSequence(uint64(0)).
		WithFee(fees).
		WithMode(ac.cfg.Mode).
		WithSimulate(false).
		WithGas(ac.cfg.Gas)
}

func (ac abstractClient) broadcastTx(txBytes []byte) (sdk.ResultTx, sdk.Error) {
	switch ac.Mode {
	case sdk.Commit:
		return ac.broadcastTxCommit(txBytes)
	case sdk.Async:
		return ac.broadcastTxAsync(txBytes)
	case sdk.Sync:
		return ac.broadcastTxSync(txBytes)

	}
	return sdk.ResultTx{}, sdk.Wrapf("commit mode(%s) not supported", ac.Mode)
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (ac abstractClient) broadcastTxCommit(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := ac.TmClient.BroadcastTxCommit(tx)
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
	res, err := ac.TmClient.BroadcastTxSync(tx)
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
	res, err := ac.TmClient.BroadcastTxAsync(tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return sdk.ResultTx{
		Hash: res.Hash.String(),
	}, nil
}
