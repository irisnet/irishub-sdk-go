package original

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/original/service"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"time"

	"github.com/irisnet/irishub-sdk-go/utils/cache"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type paramsQuery struct {
	original.Queries
	*log.Logger
	cache.Cache
	cdc        original.Codec
	expiration time.Duration
}

func (p paramsQuery) prefixKey(module string) string {
	return fmt.Sprintf("params:%s", module)
}

func (p paramsQuery) QueryParams(module string, res original.Response) original.Error {
	param, err := p.Get(p.prefixKey(module))
	if err == nil {
		bz := param.([]byte)
		err = p.cdc.UnmarshalJSON(bz, res)
		if err != nil {
			return original.Wrap(err)
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
		return original.Wrapf("unsupported param query")
	}

	bz, err := p.Query(path, nil)
	if err != nil {
		return original.Wrap(err)
	}

	err = p.cdc.UnmarshalJSON(bz, res)
	if err != nil {
		return original.Wrap(err)
	}

	if err := p.SetWithExpire(p.prefixKey(module), bz, p.expiration); err != nil {
		p.Warn().
			Str("module", module).
			Msg("params cache failed")
	}
	return nil
}
