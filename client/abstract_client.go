package client

import (
	"errors"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/types"
)

type abstractClient struct {
	*types.TxContext
}

func (ac abstractClient) Broadcast(baseTx types.BaseTx, msg []types.Msg) (types.Result, error) {
	err := ac.prepareTxContext(baseTx)
	if err != nil {
		return nil, err
	}

	tx, err := ac.BuildAndSign(baseTx.From, msg)
	if err != nil {
		return nil, err
	}

	txByte, err := ac.Codec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return nil, err
	}

	return ac.broadcastTx(txByte)
}

func (ac abstractClient) BroadcastTx(signedTx types.StdTx, mode types.BroadcastMode) (types.Result, error) {
	ac.Mode = mode
	txByte, err := ac.Codec.MarshalBinaryLengthPrefixed(signedTx)
	if err != nil {
		return nil, err
	}

	return ac.broadcastTx(txByte)
}

func (ac abstractClient) Sign(stdTx types.StdTx, name string, password string, online bool) (types.StdTx, error) {
	baseTx := types.BaseTx{
		From:     name,
		Password: password,
		Gas:      stdTx.Fee.Gas,
		Fee:      stdTx.Fee.Amount.String(),
		Memo:     stdTx.Memo,
	}
	err := ac.prepareTxContext(baseTx)
	if err != nil {
		return stdTx, err
	}

	ac.WithOnline(online)
	tx, err := ac.BuildAndSign(baseTx.From, stdTx.GetMsgs())
	if err != nil {
		return stdTx, err
	}

	return tx, nil
}

func (ac abstractClient) Query(path string, data interface{}, result interface{}) error {
	var bz []byte
	var err error
	if data != nil {
		bz, err = ac.Codec.MarshalJSON(data)
		if err != nil {
			return err
		}
	}

	res, err := ac.RPC.Query(path, bz)
	if err != nil {
		return err
	}

	if err = ac.Codec.UnmarshalJSON(res, result); err != nil {
		return err
	}

	return nil
}

func (ac abstractClient) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "subspace")
	opts := rpcclient.ABCIQueryOptions{
		//Height: cliCtx.Height,
		Prove: false,
	}

	result, err := ac.RPC.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.New(resp.Log)
	}
	return resp.Value, nil
}

func (ac abstractClient) QueryAccount(address string) (baseAccount types.BaseAccount, err error) {
	addr, err := types.AccAddressFromBech32(address)
	if err != nil {
		return baseAccount, err
	}
	param := bank.QueryAccountParams{
		Address: addr,
	}
	if err = ac.Query("custom/acc/account", param, &baseAccount); err != nil {
		return baseAccount, err
	}
	return
}

func (ac abstractClient) QueryAddress(name string) types.AccAddress {
	keyDAO := (*ac.TxContext).KeyDAO
	keystore := keyDAO.Read(name)
	return types.MustAccAddressFromBech32(keystore.GetAddress())
}

func (ac abstractClient) prepareTxContext(baseTx types.BaseTx) error {
	ctx := ac.TxContext
	if ctx.Online {
		keyStore := ctx.KeyDAO.Read(baseTx.From)
		account, err := ac.QueryAccount(keyStore.GetAddress())
		if err != nil {
			return err
		}
		ctx = ctx.WithAccountNumber(account.AccountNumber).
			WithSequence(account.Sequence)
	}
	if len(baseTx.Fee) > 0 {
		fee, err := types.ParseCoins(baseTx.Fee)
		if err != nil {
			return err
		}
		ctx = ctx.WithFee(fee)
	}

	if len(baseTx.Mode) > 0 {
		ctx = ctx.WithMode(baseTx.Mode)
	}

	if baseTx.Simulate {
		ctx = ctx.WithSimulate(baseTx.Simulate)
	}

	ctx = ctx.WithGas(baseTx.Gas)
	ctx = ctx.WithMemo(baseTx.Memo)
	return nil
}
func (ac abstractClient) broadcastTx(txBytes []byte) (types.Result, error) {
	switch ac.Mode {
	case types.Commit:
		return ac.broadcastTxCommit(txBytes)
	case types.Async:
		return ac.broadcastTxAsync(txBytes)
	case types.Sync:
		return ac.broadcastTxSync(txBytes)

	}
	panic("invalid broadcast mode")
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (ac abstractClient) broadcastTxCommit(tx []byte) (result types.ResultBroadcastTxCommit, err error) {
	res, err := ac.RPC.BroadcastTxCommit(tx)
	if err != nil {
		return result, err
	}

	if !res.CheckTx.IsOK() {
		return result, errors.New(res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return result, errors.New(res.DeliverTx.Log)
	}
	return types.ResultBroadcastTxCommit{
		CheckTx:   res.CheckTx,
		DeliverTx: res.DeliverTx,
		Hash:      res.Hash,
		Height:    res.Height,
	}, err
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously.
func (ac abstractClient) broadcastTxSync(tx []byte) (result types.ResultBroadcastTxCommit, err error) {
	res, err := ac.RPC.BroadcastTxSync(tx)
	if err != nil {
		return result, err
	}

	return types.ResultBroadcastTxCommit{
		Hash: res.Hash,
		CheckTx: abci.ResponseCheckTx{
			Code: res.Code,
			Data: res.Data,
			Log:  res.Log,
		},
	}, nil
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (ac abstractClient) broadcastTxAsync(tx []byte) (result types.ResultBroadcastTx, err error) {
	res, err := ac.RPC.BroadcastTxAsync(tx)
	if err != nil {
		return result, err
	}

	return types.ResultBroadcastTx{
		Code: res.Code,
		Data: res.Data,
		Log:  res.Log,
		Hash: res.Hash,
	}, nil
}
