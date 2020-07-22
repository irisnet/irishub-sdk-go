//Package tendermint provides tendermint rpc queriers implementation
//
//
package tendermint

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/p2p"
)

const (
	ModuleName = "tendermint"
)

type tmClient struct {
	sdk.BaseClient
	cdc sdk.Codec
}

func Create(ac sdk.BaseClient) rpc.Tendermint {
	return tmClient{
		BaseClient: ac,
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
	tx, err := t.BaseClient.QueryTx(hash)
	if err != nil {
		return sdk.ResultQueryTx{}, sdk.Wrap(err)
	}
	return tx, nil
}

func (t tmClient) SearchTxs(builder *sdk.EventQueryBuilder, page, size int) (sdk.ResultSearchTxs, sdk.Error) {
	txs, err := t.BaseClient.QueryTxs(builder, page, size)
	if err != nil {
		return sdk.ResultSearchTxs{}, sdk.Wrap(err)
	}
	return txs, nil
}

func (t tmClient) QueryValidators(height int64) (rpc.ResultQueryValidators, sdk.Error) {
	rs, err := t.Validators(&height, 0, 100)
	if err != nil {
		return rpc.ResultQueryValidators{}, sdk.Wrap(err)
	}
	return rpc.ResultQueryValidators{
		BlockHeight: rs.BlockHeight,
		Validators:  sdk.ParseValidators(rs.Validators),
	}, nil
}

func (t tmClient) QueryNodeInfo() (p2p.DefaultNodeInfo, sdk.Error) {
	status, err := t.Status()
	if err != nil {
		return p2p.DefaultNodeInfo{}, sdk.Wrap(err)
	}
	return status.NodeInfo, nil
}

func (t tmClient) QueryNodeVersion() (string, sdk.Error) {
	var version string
	bz, err := t.Query("/app/version", nil)
	if err != nil {
		return "", sdk.Wrap(err)
	}
	version = string(bz)
	return version, nil
}
