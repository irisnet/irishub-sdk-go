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

// QueryParams By moduleName , if the moduleName isn't pass will return params of all module
// Support Module: auth,distribution,gov,mint,service,slashing,staking,token
func (p paramClient) QueryParams(moduleName string) (rpc.ParamsResponses, sdk.Error) {
	var prs paramsResponses
	if len(moduleName) == 0 {
		for _, v := range AllModule {
			res, err := p.queryParams(v)
			if err != nil {
				return rpc.ParamsResponses{}, sdk.Wrap(err)
			}
			prs = append(prs, res)
		}
	} else {
		res, err := p.queryParams(moduleName)
		if err != nil {
			return rpc.ParamsResponses{}, sdk.Wrap(err)
		}
		prs = append(prs, res)
	}
	return prs.Convert().(rpc.ParamsResponses), nil
}

// QueryParams By moduleName (moduleName must be have)
func (p paramClient) queryParams(moduleName string) (paramsResponse, sdk.Error) {
	routers := getRouters(moduleName)
	if len(routers) == 0 {
		return paramsResponse{}, sdk.Wrapf("the params of module isn't exist")
	}

	var ps paramsResponse
	for _, router := range routers {
		res, err := p.Query(router, nil)
		if err != nil {
			return paramsResponse{}, sdk.Wrap(err)
		}
		ps.Type = moduleName + "/Params"
		ps.Value += string(res)
	}
	return ps, nil
}

// getRouters By moduleName , GOV  will return 3 router ,the others is 1
func getRouters(moduleName string) []string {
	var router []string
	switch moduleName {
	case AUTH:
		fallthrough
	case DISTRIBUTION:
		fallthrough
	case TOKEN:
		router = append(router, fmt.Sprintf("custom/%s/params", moduleName))
	case MINT:
		fallthrough
	case SERVICE:
		fallthrough
	case SLASHING:
		fallthrough
	case STAKING:
		router = append(router, fmt.Sprintf("custom/%s/parameters", moduleName))
	case GOV:
		router = append(router, fmt.Sprintf("custom/%s/params/deposit", moduleName))
		router = append(router, fmt.Sprintf("custom/%s/params/voting", moduleName))
		router = append(router, fmt.Sprintf("custom/%s/params/tallying", moduleName))
	default:
	}
	return router
}
