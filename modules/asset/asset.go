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

func (a assetClient) QueryTokens(owner string) (sdk.Tokens, error) {
	param := struct {
		Symbol string
		Owner  string
	}{
		Owner: owner,
	}

	var tokens sdk.Tokens
	if err := a.QueryWithResponse("custom/asset/tokens", param, &tokens); err != nil {
		return sdk.Tokens{}, err
	}
	sdk.CacheTokens(tokens...)
	return tokens, nil
}

func (a assetClient) QueryFees(symbol string) (rpc.TokenFees, error) {
	param := struct {
		Symbol string
	}{
		Symbol: symbol,
	}

	var tokens tokenFees
	if err := a.QueryWithResponse("custom/asset/fees", param, &tokens); err != nil {
		return rpc.TokenFees{}, err
	}
	return tokens.Convert().(rpc.TokenFees), nil
}
