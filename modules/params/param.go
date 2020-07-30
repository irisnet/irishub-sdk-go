package params

import (
	"fmt"
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

// QueryParams By subspace and key this query need certain args
func (p paramClient) QueryParamsBySubAndKey(subspace, key string) (rpc.SubspaceParamsResponse, sdk.Error) {
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

// QueryParams By moduleName , if the moduleName isn't pass will return params of all module
// Support Module: auth,distribution,gov,mint,service,slashing,staking,token
func (p paramClient) QueryParams(moduleName string) ([]byte, sdk.Error) {
	//if len(module) > 0 {
	//router := fmt.Sprintf("custom/%s/parameters", moduleName)
	router := getRouter(moduleName)
	if len(router) == 0 {
		return []byte{}, sdk.Wrapf("the params of module isn't exist")
	}
	res, err := p.Query(router, nil)
	if err != nil {
		return []byte{}, sdk.Wrap(err)
	}
	return res, nil
}

func getRouter(moduleName string) string {
	var router string
	switch moduleName {
	case AUTH:
		fallthrough
	case DISTRIBUTION:
		fallthrough
	case GOV:
		fallthrough
	case TOKEN:
		router = fmt.Sprintf("custom/%s/params", moduleName)
	case MINT:
		fallthrough
	case SERVICE:
		fallthrough
	case SLASHING:
		fallthrough
	case STAKING:
		router = fmt.Sprintf("custom/%s/params", moduleName)
	default:
		router = ""
	}
	return router
}
