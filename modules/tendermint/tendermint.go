package tendermint

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "tendermint"
)

type tmClient struct {
	sdk.AbstractClient
	cdc sdk.Codec
}

func Create(ac sdk.AbstractClient) rpc.Tendermint {
	return tmClient{
		AbstractClient: ac,
	}
}

func (t tmClient) RegisterCodec(cdc sdk.Codec) {
	//nothing
}

func (t tmClient) Name() string {
	return ModuleName
}

func (t tmClient) QueryBlock(height int64) (sdk.Block, sdk.Error) {
	block, err := t.Block(&height)
	if err != nil {
		return sdk.Block{}, sdk.Wrap(err)
	}
	return sdk.ParseBlock(t.cdc, block.Block), nil
}

func (t tmClient) QueryBlockResult(height int64) (sdk.BlockResult, sdk.Error) {
	blockResult, err := t.BlockResults(&height)
	if err != nil {
		return sdk.BlockResult{}, sdk.Wrap(err)
	}
	return sdk.ParseBlockResult(blockResult), nil
}

func (t tmClient) QueryTx(hash string) (sdk.ResultQueryTx, sdk.Error) {
	tx, err := t.AbstractClient.QueryTx(hash)
	if err != nil {
		return sdk.ResultQueryTx{}, sdk.Wrap(err)
	}
	return tx, nil
}

func (t tmClient) SearchTxs(builder *sdk.EventQueryBuilder, page, size int) (sdk.ResultSearchTxs, sdk.Error) {
	txs, err := t.AbstractClient.QueryTxs(builder, page, size)
	if err != nil {
		return sdk.ResultSearchTxs{}, sdk.Wrap(err)
	}
	return txs, nil
}

func (t tmClient) QueryValidators(height int64) (rpc.ResultQueryValidators, sdk.Error) {
	rs, err := t.Validators(&height)
	if err != nil {
		return rpc.ResultQueryValidators{}, sdk.Wrap(err)
	}
	return rpc.ResultQueryValidators{
		BlockHeight: rs.BlockHeight,
		Validators:  sdk.ParseValidators(rs.Validators),
	}, nil
}
