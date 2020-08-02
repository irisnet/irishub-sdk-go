package modules

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/service"
	"time"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/cache"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type paramsQuery struct {
	sdk.Queries
	*log.Logger
	cache.Cache
	cdc        sdk.Codec
	expiration time.Duration
}

func (p paramsQuery) prefixKey(module string) string {
	return fmt.Sprintf("params:%s", module)
}

func (p paramsQuery) QueryParams(module string, res sdk.Response) sdk.Error {
	param, err := p.Get(p.prefixKey(module))
	if err == nil {
		bz := param.([]byte)
		err = p.cdc.UnmarshalJSON(bz, res)
		if err != nil {
			return sdk.Wrap(err)
		}
		return nil
	}

	var path string
	switch module {
	case service.ModuleName:
		path = fmt.Sprintf("custom/%s/parameters", module)
	case "auth":
		path = fmt.Sprintf("custom/%s/params", "auth")
	default:
		return sdk.Wrapf("unsupported param query")
	}

	bz, err := p.Query(path, nil)
	if err != nil {
		return sdk.Wrap(err)
	}

	err = p.cdc.UnmarshalJSON(bz, res)
	if err != nil {
		return sdk.Wrap(err)
	}

	if err := p.SetWithExpire(p.prefixKey(module), bz, p.expiration); err != nil {
		p.Warn().
			Str("module", module).
			Msg("params cache failed")
	}
	return nil
}
