package params

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type paramsClient struct {
	sdk.BaseClient
	*log.Logger
}

func Create(ac sdk.BaseClient) rpc.Htlc {
	return paramsClient{}
}
