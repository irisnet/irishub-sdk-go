package rpc

import sdk "github.com/irisnet/irishub-sdk-go/types"

type Params interface {
	sdk.Module
	QueryParamsBySubAndKey(subspace, key string) (SubspaceParamsResponse, sdk.Error)
	QueryParams(moduleName string) ([]byte, sdk.Error)
}

type SubspaceParamsResponse struct {
	Subspace string
	Key      string
	Value    string
}
