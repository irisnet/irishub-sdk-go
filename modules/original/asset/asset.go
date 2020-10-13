// Package asset allows individuals and companies to create and issue their own tokens.
//
// [More Details](https://www.irisnet.org/docs/features/asset.html)
package asset

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type assetClient struct {
	original.BaseClient
	*log.Logger
}

func Create(ac original.BaseClient) rpc.Asset {
	return assetClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (a assetClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (a assetClient) Name() string {
	return ModuleName
}

func (a assetClient) QueryTokens() (original.Tokens, error) {
	param := struct {
		Owner string
	}{
		Owner: "",
	}

	var tokens original.Tokens
	if err := a.QueryWithResponse("custom/token/tokens", param, &tokens); err != nil {
		return original.Tokens{}, err
	}
	return tokens, nil
}

func (a assetClient) QueryTokenDenom(denom string) (original.TokenData, error) {
	uri := fmt.Sprintf("custom/%s/token", ModuleName)
	param := struct {
		Denom string
	}{
		Denom: denom,
	}
	var tokendata original.TokenData
	if err := a.QueryWithResponse(uri, param, &tokendata); err != nil {
		return original.TokenData{}, err
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
