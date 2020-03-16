package asset

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type assetClient struct {
	sdk.AbstractClient
	*log.Logger
}

func Create(ac sdk.AbstractClient) rpc.Asset {
	return assetClient{
		AbstractClient: ac,
		Logger:         ac.Logger(),
	}
}

func (a assetClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (a assetClient) Name() string {
	return ModuleName
}

func (a assetClient) QueryToken(symbol string) (sdk.Token, error) {
	return a.AbstractClient.QueryToken(symbol)
}
