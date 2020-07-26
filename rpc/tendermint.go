package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Tendermint interface {
	sdk.Module
	QueryBlock(height int64) (sdk.Block, sdk.Error)
	QueryBlockLatest() (sdk.Block, sdk.Error)
	QueryBlockResult(height int64) (sdk.BlockResult, sdk.Error)
	QueryTx(hash string) (sdk.ResultQueryTx, sdk.Error)
	SearchTxs(builder *sdk.EventQueryBuilder, page, size int) (sdk.ResultSearchTxs, sdk.Error)
	QueryValidators(height int64) (ResultValidators, sdk.Error)
	QueryValidatorsLatest() (ResultValidators, sdk.Error)
	QueryNodeInfo() (sdk.ResultStatus, sdk.Error)
	QueryNodeVersion() (string, sdk.Error)
	QueryGenesis() (sdk.GenesisDoc, sdk.Error)
}

// Validators for a height.
type ResultValidators struct {
	BlockHeight int64           `json:"block_height"`
	Validators  []sdk.Validator `json:"validators"`
	Count       int             `json:"count"`
	Total       int             `json:"total"`
}
