package original

import (
	"encoding/hex"
	"errors"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// QueryTx returns the tx info
func (base baseClient) QueryTx(hash string) (original.ResultQueryTx, error) {
	tx, err := hex.DecodeString(hash)
	if err != nil {
		return original.ResultQueryTx{}, err
	}

	res, err := base.Tx(tx, true)
	if err != nil {
		return original.ResultQueryTx{}, err
	}

	resBlocks, err := base.getResultBlocks([]*ctypes.ResultTx{res})
	if err != nil {
		return original.ResultQueryTx{}, err
	}
	return base.parseTxResult(res, resBlocks[res.Height])
}

func (base baseClient) QueryTxs(builder *original.EventQueryBuilder, page, size int) (original.ResultSearchTxs, error) {

	query := builder.Build()
	if len(query) == 0 {
		return original.ResultSearchTxs{}, errors.New("must declare at least one tag to search")
	}

	res, err := base.TxSearch(query, true, page, size, "asc")
	if err != nil {
		return original.ResultSearchTxs{}, err
	}

	resBlocks, err := base.getResultBlocks(res.Txs)
	if err != nil {
		return original.ResultSearchTxs{}, err
	}

	var txs []original.ResultQueryTx
	for i, tx := range res.Txs {
		txInfo, err := base.parseTxResult(tx, resBlocks[res.Txs[i].Height])
		if err != nil {
			return original.ResultSearchTxs{}, err
		}
		txs = append(txs, txInfo)
	}

	return original.ResultSearchTxs{
		Total: res.TotalCount,
		Txs:   txs,
	}, nil
}

func (base *baseClient) buildTx(msg []original.Msg, baseTx original.BaseTx) ([]byte, *original.TxContext, original.Error) {
	ctx, err := base.prepare(baseTx)
	if err != nil {
		return nil, ctx, original.Wrap(err)
	}

	tx, err := ctx.BuildAndSign(baseTx.From, msg)
	if err != nil {
		return nil, ctx, original.Wrap(err)
	}

	base.Logger().Debug().
		Strs("data", tx.GetSignBytes()).
		Msg("sign transaction success")

	// TODO is this correct??
	txByte, err := base.cdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return nil, ctx, original.Wrap(err)
	}

	return txByte, ctx, nil
}

func (base baseClient) broadcastTx(txBytes []byte, mode original.BroadcastMode) (res original.ResultTx, err original.Error) {
	ch := make(chan original.ResultTx, 1)

	go func() {
		switch mode {
		case original.Commit:
			res, err = base.broadcastTxCommit(txBytes)
		case original.Async:
			res, err = base.broadcastTxAsync(txBytes)
		case original.Sync:
			res, err = base.broadcastTxSync(txBytes)
		default:
			err = original.Wrapf("commit mode(%s) not supported", base.cfg.Mode)
		}
		ch <- res
	}()

	select {
	case result := <-ch:
		return result, err
	case <-time.After(base.cfg.Timeout):
		return res, original.Wrap(errors.New("commit transaction timed out"))

	}
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (base baseClient) broadcastTxCommit(tx []byte) (original.ResultTx, original.Error) {
	res, err := base.BroadcastTxCommit(tx)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	if !res.CheckTx.IsOK() {
		return original.ResultTx{}, original.GetError(res.CheckTx.Codespace,
			res.CheckTx.Code, res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return original.ResultTx{}, original.GetError(res.DeliverTx.Codespace,
			res.DeliverTx.Code, res.DeliverTx.Log)
	}

	return original.ResultTx{
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Events:    original.ParseEvents(res.DeliverTx.Events),
		Hash:      res.Hash.String(),
		Height:    res.Height,
	}, nil
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously.
func (base baseClient) broadcastTxSync(tx []byte) (original.ResultTx, original.Error) {
	res, err := base.BroadcastTxSync(tx)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	if res.Code != 0 {
		return original.ResultTx{}, original.GetError(original.RootCodespace,
			res.Code, res.Log)
	}

	return original.ResultTx{
		Hash: res.Hash.String(),
	}, nil
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (base baseClient) broadcastTxAsync(tx []byte) (original.ResultTx, original.Error) {
	res, err := base.BroadcastTxAsync(tx)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	return original.ResultTx{
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

func (base baseClient) parseTxResult(res *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (original.ResultQueryTx, error) {

	var tx original.StdTx
	err := base.cdc.UnmarshalBinaryBare(res.Tx, &tx)
	if err != nil {
		return original.ResultQueryTx{}, err
	}

	return original.ResultQueryTx{
		Hash:   res.Hash.String(),
		Height: res.Height,
		Tx:     tx,
		Result: original.TxResult{
			Code:      res.TxResult.Code,
			Log:       res.TxResult.Log,
			GasWanted: res.TxResult.GasWanted,
			GasUsed:   res.TxResult.GasUsed,
		},
		Timestamp: resBlock.Block.Time.Format(time.RFC3339),
	}, nil
}
