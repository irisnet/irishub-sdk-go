// Package asset allows individuals and companies to create and issue their own tokens.
//
// [More Details](https://www.irisnet.org/docs/features/asset.html)
package asset

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type assetClient struct {
	sdk.BaseClient
	*log.Logger
}

func Create(ac sdk.BaseClient) rpc.Asset {
	return assetClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (a assetClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (a assetClient) Name() string {
	return ModuleName
}

func (a assetClient) QueryTokens() (sdk.Tokens, error) {
	param := struct {
		Owner string
	}{
		Owner: "",
	}

	var tokens sdk.Tokens
	if err := a.QueryWithResponse("custom/token/tokens", param, &tokens); err != nil {
		return sdk.Tokens{}, err
	}
	return tokens, nil
}
