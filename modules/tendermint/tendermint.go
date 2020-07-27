//Package tendermint provides tendermint rpc queriers implementation
//
//
package tendermint

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "tendermint"
)

type tmClient struct {
	sdk.BaseClient
	cdc sdk.Codec
}

func Create(ac sdk.BaseClient, cdc sdk.Codec) rpc.Tendermint {
	return tmClient{
		BaseClient: ac,
		cdc:        cdc,
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
	bz, _ := json.Marshal(block)
	fmt.Println(bz)
	if err != nil {
		return sdk.Block{}, sdk.Wrap(err)
	}
	return sdk.ParseBlock(t.cdc, block.Block), nil
}

func (t tmClient) QueryBlockLatest() (sdk.Block, sdk.Error) {
	status, err := t.Status()
	if err != nil {
		return sdk.Block{}, sdk.Wrap(err)
	}
	return t.QueryBlock(status.SyncInfo.LatestBlockHeight)
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

func (t tmClient) QueryValidators(height int64) (rpc.ResultValidators, sdk.Error) {
	rs, err := t.Validators(&height, 0, 100)
	if err != nil {
		return rpc.ResultValidators{}, sdk.Wrap(err)
	}
	return rpc.ResultValidators{
		BlockHeight: rs.BlockHeight,
		Validators:  sdk.ParseValidators(rs.Validators),
		Count:       rs.Count,
		Total:       rs.Total,
	}, nil
}

func (t tmClient) QueryValidatorsLatest() (rpc.ResultValidators, sdk.Error) {
	status, err := t.Status()
	if err != nil {
		return rpc.ResultValidators{}, sdk.Wrap(err)
	}
	return t.QueryValidators(status.SyncInfo.LatestBlockHeight)
}

func (t tmClient) QueryNodeInfo() (sdk.ResultStatus, sdk.Error) {
	status, err := t.Status()
	if err != nil {
		return sdk.ResultStatus{}, sdk.Wrap(err)
	}
	return sdk.ParseNodeStatus(status), nil
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

func (t tmClient) QueryGenesis() (sdk.GenesisDoc, sdk.Error) {
	genesis, err := t.Genesis()
	if err != nil {
		return sdk.GenesisDoc{}, sdk.Wrap(err)
	}
	return sdk.ParseGenesis(genesis.Genesis), nil
}
