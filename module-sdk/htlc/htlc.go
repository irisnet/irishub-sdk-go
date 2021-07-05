package htlc

import (
	"context"

	"github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

type htlcClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(baseClient sdk.BaseClient, marshaler codec.Marshaler) Client {
	return htlcClient{
		BaseClient: baseClient,
		Marshaler:  marshaler,
	}
}

func (hc htlcClient) Name() string {
	return ModuleName
}

func (hc htlcClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (hc htlcClient) CreateHTLC(request CreateHTLCRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := hc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if request.TimeLock == 0 {
		request.TimeLock = MinTimeLock
	}

	amount, err := hc.ToMinCoin(request.Amount...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgCreateHTLC{
		Sender:               sender.String(),
		To:                   request.To,
		ReceiverOnOtherChain: request.ReceiverOnOtherChain,
		SenderOnOtherChain:   request.SenderOnOtherChain,
		Amount:               amount,
		HashLock:             request.HashLock,
		Timestamp:            request.Timestamp,
		TimeLock:             request.TimeLock,
		Transfer:             request.Transfer,
	}
	return hc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (hc htlcClient) ClaimHTLC(hashLockId string, secret string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := hc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgClaimHTLC{
		Sender: sender.String(),
		Id:     hashLockId,
		Secret: secret,
	}
	return hc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (hc htlcClient) QueryHTLC(hashLockId string) (QueryHTLCResp, sdk.Error) {
	if len(hashLockId) == 0 {
		return QueryHTLCResp{}, sdk.Wrapf("hashLock id is required")
	}

	conn, err := hc.GenConn()

	if err != nil {
		return QueryHTLCResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).HTLC(
		context.Background(),
		&QueryHTLCRequest{
			Id: hashLockId,
		})
	if err != nil {
		return QueryHTLCResp{}, sdk.Wrap(err)
	}
	return res.Htlc.Convert().(QueryHTLCResp), nil
}

func (hc htlcClient) QueryParams() (QueryParamsResp, sdk.Error) {

	conn, err := hc.GenConn()

	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{})
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}
	return res.Params.Convert().(QueryParamsResp), nil
}
