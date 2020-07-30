package params

import "github.com/irisnet/irishub-sdk-go/rpc"

const (
	ModuleName = "params"

	AUTH         = "auth"
	DISTRIBUTION = "distribution"
	GOV          = "gov"
	MINT         = "mint"
	SERVICE      = "service"
	SLASHING     = "slashing"
	STAKING      = "staking"
	TOKEN        = "token"
)

var AllModule = []string{AUTH, DISTRIBUTION, GOV, MINT, SERVICE, SLASHING, STAKING, TOKEN}

type paramsResponses []paramsResponse
type paramsResponse struct {
	Type  string
	Value string
}

func (ps paramsResponses) Convert() interface{} {
	paramsResponse := make(rpc.ParamsResponses, len(ps))
	for i, v := range ps {
		paramsResponse[i] = v.Convert().(rpc.ParamsResponse)
	}
	return paramsResponse
}

func (p paramsResponse) Convert() interface{} {
	return rpc.ParamsResponse{
		Type:  p.Type,
		Value: p.Value,
	}
}
