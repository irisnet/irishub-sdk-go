package modules

import (
	"encoding/hex"
	"errors"
	"time"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

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

func (base *baseClient) buildTx(msg []sdk.Msg, baseTx sdk.BaseTx) ([]byte, *sdk.TxContext, sdk.Error) {
	ctx, err := base.prepare(baseTx)
	if err != nil {
		return nil, ctx, sdk.Wrap(err)
	}

	tx, err := ctx.BuildAndSign(baseTx.From, msg)
	if err != nil {
		return nil, ctx, sdk.Wrap(err)
	}

	base.Logger().Debug().
		Strs("data", tx.GetSignBytes()).
		Msg("sign transaction success")

	txByte, err := base.cdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return nil, ctx, sdk.Wrap(err)
	}

	return txByte, ctx, nil
}

func (base baseClient) broadcastTx(txBytes []byte, mode sdk.BroadcastMode) (res sdk.ResultTx, err sdk.Error) {
	ch := make(chan sdk.ResultTx, 1)

	go func() {
		switch mode {
		case sdk.Commit:
			res, err = base.broadcastTxCommit(txBytes)
		case sdk.Async:
			res, err = base.broadcastTxAsync(txBytes)
		case sdk.Sync:
			res, err = base.broadcastTxSync(txBytes)
		default:
			err = sdk.Wrapf("commit mode(%s) not supported", base.cfg.Mode)
		}
		ch <- res
	}()

	select {
	case result := <-ch:
		return result, err
	case <-time.After(base.cfg.Timeout):
		return res, sdk.Wrap(errors.New("commit transaction timed out"))

	}
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
