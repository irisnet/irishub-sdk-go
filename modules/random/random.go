package random

import (
	"context"
	"github.com/irisnet/irishub-sdk-go/codec"
	cdctypes "github.com/irisnet/irishub-sdk-go/codec/types"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type randomClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(baseClient sdk.BaseClient, marshaler codec.Marshaler) *randomClient {
	return &randomClient{
		BaseClient: baseClient,
		Marshaler:  marshaler,
	}
}

func (rc randomClient) Name() string {
	return ModuleName
}

func (rc randomClient) RegisterInterfaceTypes(registry cdctypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (rc randomClient) QueryRandom(reqID string) (QueryRandomResp, sdk.Error) {
	if len(reqID) == 0 {
		return QueryRandomResp{}, sdk.Wrapf("reqId is required")
	}

	conn, err := rc.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryRandomResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Random(
		context.Background(),
		&QueryRandomRequest{ReqId: reqID},
	)
	if err != nil {
		return QueryRandomResp{}, sdk.Wrap(err)
	}
	return res.Random.Convert().(QueryRandomResp), nil
}

func (rc randomClient) QueryRandomRequestQueue(height int64) ([]QueryRandomRequestQueueResp, sdk.Error) {
	if height == 0 {
		return []QueryRandomRequestQueueResp{}, nil
	}

	conn, err := rc.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return []QueryRandomRequestQueueResp{}, sdk.Wrap(err)
	}
	res, err := NewQueryClient(conn).RandomRequestQueue(
		context.Background(),
		&QueryRandomRequestQueueRequest{Height: height},
	)
	if err != nil {
		return []QueryRandomRequestQueueResp{}, sdk.Wrap(err)
	}
	return Requests(res.Requests).Convert().([]QueryRandomRequestQueueResp), nil
}
