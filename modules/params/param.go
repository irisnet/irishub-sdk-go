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

func (p paramClient) QueryParams(subspace, key string) (rpc.SubspaceParamsResponse, sdk.Error) {
	param := struct {
		Subspace, Key string
	}{
		Subspace: subspace,
		Key:      key,
	}

	var sr subspaceParamsResponse
	if err := p.QueryWithResponse("custom/params/params", param, &sr); err != nil {
		return rpc.SubspaceParamsResponse{}, sdk.Wrap(err)
	}
	return sr.Convert().(rpc.SubspaceParamsResponse), nil
}
