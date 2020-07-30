package rpc

import sdk "github.com/irisnet/irishub-sdk-go/types"

type Params interface {
	sdk.Module
	QueryParams(moduleName string) (ParamsResponses, sdk.Error)
}

type ParamsResponses []ParamsResponse
type ParamsResponse struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
