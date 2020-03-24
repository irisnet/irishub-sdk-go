package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Tendermint interface {
	sdk.Module
	QueryBlock(height int64) (sdk.Block, sdk.Error)
	QueryBlockResult(height int64) (sdk.BlockResult, sdk.Error)
	QueryTx(hash string) (sdk.TxDetail, sdk.Error)
	QueryTxs(builder *sdk.EventQueryBuilder, page, size int) (sdk.TxSearch, sdk.Error)
	QueryValidatorSet(height int64) (sdk.ResultValidators, sdk.Error)
}
