// Package asset allows individuals and companies to create and issue their own tokens.
//
// [More Details](https://www.irisnet.org/docs/features/asset.html)
package asset

import (
	"fmt"
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

func (a assetClient) QueryTokenDenom(denom string) (sdk.TokenData, error) {
	uri := fmt.Sprintf("custom/%s/token", ModuleName)
	param := struct {
		Denom string
	}{
		Denom: denom,
	}
	var tokendata sdk.TokenData
	if err := a.QueryWithResponse(uri, param, &tokendata); err != nil {
		return sdk.TokenData{}, err
	}
	return tokendata, nil
}

func (a assetClient) QueryFees(symbol string) (rpc.TokenFees, error) {
	param := struct {
		Symbol string
	}{
		Symbol: symbol,
	}
	uri := fmt.Sprintf("custom/%s/fees/tokens", ModuleName)

	var tokens tokenFees
	if err := a.QueryWithResponse(uri, param, &tokens); err != nil {
		return rpc.TokenFees{}, err
	}
	return tokens.Convert().(rpc.TokenFees), nil
}
