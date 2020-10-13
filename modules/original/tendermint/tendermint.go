//Package tendermint provides tendermint rpc queriers implementation
//
//
package tendermint

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
)

const (
	ModuleName = "tendermint"
)

type tmClient struct {
	original.BaseClient
	cdc original.Codec
}

func Create(ac original.BaseClient, cdc original.Codec) rpc.Tendermint {
	return tmClient{
		BaseClient: ac,
		cdc:        cdc,
	}
}

func (t tmClient) RegisterCodec(cdc original.Codec) {
	//nothing
}

func (t tmClient) Name() string {
	return ModuleName
}

func (t tmClient) QueryBlock(height int64) (original.Block, original.Error) {
	block, err := t.Block(&height)
	if err != nil {
		return original.Block{}, original.Wrap(err)
	}
	return original.ParseBlock(t.cdc, block.Block), nil
}

func (t tmClient) QueryBlockLatest() (original.Block, original.Error) {
	status, err := t.Status()
	if err != nil {
		return original.Block{}, original.Wrap(err)
	}
	return t.QueryBlock(status.SyncInfo.LatestBlockHeight)
}

func (t tmClient) QueryBlockResult(height int64) (original.BlockResult, original.Error) {
	blockResult, err := t.BlockResults(&height)
	if err != nil {
		return original.BlockResult{}, original.Wrap(err)
	}
	return original.ParseBlockResult(blockResult), nil
}

func (t tmClient) QueryTx(hash string) (original.ResultQueryTx, original.Error) {
	tx, err := t.BaseClient.QueryTx(hash)
	if err != nil {
		return original.ResultQueryTx{}, original.Wrap(err)
	}
	return tx, nil
}

func (t tmClient) SearchTxs(builder *original.EventQueryBuilder, page, size int) (original.ResultSearchTxs, original.Error) {
	txs, err := t.BaseClient.QueryTxs(builder, page, size)
	if err != nil {
		return original.ResultSearchTxs{}, original.Wrap(err)
	}
	return txs, nil
}

func (t tmClient) QueryValidators(height int64) (rpc.ResultValidators, original.Error) {
	rs, err := t.Validators(&height, 0, 100)
	if err != nil {
		return rpc.ResultValidators{}, original.Wrap(err)
	}
	return rpc.ResultValidators{
		BlockHeight: rs.BlockHeight,
		Validators:  original.ParseValidators(rs.Validators),
		Count:       rs.Count,
		Total:       rs.Total,
	}, nil
}

func (t tmClient) QueryValidatorsLatest() (rpc.ResultValidators, original.Error) {
	status, err := t.Status()
	if err != nil {
		return rpc.ResultValidators{}, original.Wrap(err)
	}
	return t.QueryValidators(status.SyncInfo.LatestBlockHeight)
}

func (t tmClient) QueryNodeInfo() (original.ResultStatus, original.Error) {
	status, err := t.Status()
	if err != nil {
		return original.ResultStatus{}, original.Wrap(err)
	}
	return original.ParseNodeStatus(status), nil
}

func (t tmClient) QueryNodeVersion() (string, original.Error) {
	var version string
	bz, err := t.Query("/app/version", nil)
	if err != nil {
		return "", original.Wrap(err)
	}
	version = string(bz)
	return version, nil
}

func (t tmClient) QueryGenesis() (original.GenesisDoc, original.Error) {
	genesis, err := t.Genesis()
	if err != nil {
		return original.GenesisDoc{}, original.Wrap(err)
	}
	return original.ParseGenesis(genesis.Genesis), nil
}
