package params

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type paramClient struct {
	sdk.BaseClient
	*log.Logger
}

func (p paramClient) RegisterCodec(cdc sdk.Codec) {

}

func (p paramClient) Name() string {
	return ModuleName
}

func Create(ac sdk.BaseClient) rpc.Params {
	return paramClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (p paramClient) QueryParams() (string, error) {
	return "", nil
}
