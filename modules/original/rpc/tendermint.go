package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

type Tendermint interface {
	original.Module
	QueryBlock(height int64) (original.Block, original.Error)
	QueryBlockLatest() (original.Block, original.Error)
	QueryBlockResult(height int64) (original.BlockResult, original.Error)
	QueryTx(hash string) (original.ResultQueryTx, original.Error)
	SearchTxs(builder *original.EventQueryBuilder, page, size int) (original.ResultSearchTxs, original.Error)
	QueryValidators(height int64) (ResultValidators, original.Error)
	QueryValidatorsLatest() (ResultValidators, original.Error)
	QueryNodeInfo() (original.ResultStatus, original.Error)
	QueryNodeVersion() (string, original.Error)
	QueryGenesis() (original.GenesisDoc, original.Error)
}

// Validators for a height.
type ResultValidators struct {
	BlockHeight int64                `json:"block_height"`
	Validators  []original.Validator `json:"validators"`
	Count       int                  `json:"count"`
	Total       int                  `json:"total"`
}
