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

// SubspaceParamsResponse defines the response for quering parameters by subspace.
type subspaceParamsResponse struct {
	Subspace string
	Key      string
	Value    string
}

func (s subspaceParamsResponse) Convert() interface{} {
	return rpc.SubspaceParamsResponse{
		Subspace: s.Subspace,
		Key:      s.Key,
		Value:    s.Value,
	}
}
