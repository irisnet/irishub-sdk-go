package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

type Params interface {
	original.Module
	QueryParams(moduleName string) (ParamsResponses, original.Error)
}

type ParamsResponses []ParamsResponse
type ParamsResponse struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
