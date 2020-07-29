package rpc

import sdk "github.com/irisnet/irishub-sdk-go/types"

type Params interface {
	sdk.Module
	QueryParams(subspace, key string) (SubspaceParamsResponse, sdk.Error)
	//QueryParams(module string) (SubspaceParamsResponse, sdk.Error)
}

type SubspaceParamsResponse struct {
	Subspace string
	Key      string
	Value    string
}
