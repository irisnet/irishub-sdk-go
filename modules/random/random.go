package random

import (
	"context"

	"github.com/irisnet/irishub-sdk-go/codec"
	"github.com/irisnet/irishub-sdk-go/codec/types"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type randomClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) RandomI {
	return randomClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (r randomClient) Name() string {
	return ModuleName
}

func (r randomClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (r randomClient) QueryRandom(requestID string) (QueryRandomResp, sdk.Error) {
	conn, err := r.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryRandomResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Random(
		context.Background(),
		&QueryRandomRequest{ReqId: requestID},
	)
	if err != nil {
		return QueryRandomResp{}, sdk.Wrap(err)
	}

	return res.Random.Convert().(QueryRandomResp), nil
}
