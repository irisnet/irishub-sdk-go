package random

import (
	"context"
	"strconv"

	"github.com/irisnet/core-sdk-go/common/codec"
	cdctypes "github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
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

func (rc randomClient) RequestRandom(request RequestRandomRequest, basTx sdk.BaseTx) (RequestRandomResp, sdk.ResultTx, sdk.Error) {
	author, err := rc.QueryAddress(basTx.From, basTx.Password)
	if err != nil {
		return RequestRandomResp{}, sdk.ResultTx{}, nil
	}

	msg := &MsgRequestRandom{
		BlockInterval: request.BlockInterval,
		Consumer:      author.String(),
		Oracle:        request.Oracle,
		ServiceFeeCap: request.ServiceFeeCap,
	}
	result, err := rc.BuildAndSend([]sdk.Msg{msg}, basTx)
	if err != nil {
		return RequestRandomResp{}, sdk.ResultTx{}, err
	}

	reqID, e := result.Events.GetValue(eventTypeRequestRequestRandom, attributeKeyRequestID)
	if e != nil {
		return RequestRandomResp{}, result, sdk.Wrap(e)
	}
	generateHeight, e := result.Events.GetValue(eventTypeRequestRequestRandom, attributeKeyGenerateHeight)
	if e != nil {
		return RequestRandomResp{}, result, sdk.Wrap(e)
	}
	height, e := strconv.Atoi(generateHeight)
	if e != nil {
		return RequestRandomResp{}, result, sdk.Wrap(e)
	}

	res := RequestRandomResp{
		Height: int64(height),
		ReqID:  reqID,
	}
	return res, result, nil
}

func (rc randomClient) QueryRandom(reqID string) (QueryRandomResp, sdk.Error) {
	if len(reqID) == 0 {
		return QueryRandomResp{}, sdk.Wrapf("reqId is required")
	}

	conn, err := rc.GenConn()

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
