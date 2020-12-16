package htlc

import (
	"context"
	"github.com/irisnet/irishub-sdk-go/codec"
	"github.com/irisnet/irishub-sdk-go/codec/types"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type htlcClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(baseClient sdk.BaseClient, marshaler codec.Marshaler) *htlcClient {
	return &htlcClient{
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
		Amount:               amount,
		HashLock:             request.HashLock,
		Timestamp:            request.Timestamp,
		TimeLock:             request.TimeLock,
	}
	return hc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (hc htlcClient) ClaimHTLC(hashLock string, secret string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := hc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgClaimHTLC{
		Sender:   sender.String(),
		HashLock: hashLock,
		Secret:   secret,
	}
	return hc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (hc htlcClient) RefundHTLC(hashLock string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	if len(hashLock) == 0 {
		return sdk.ResultTx{}, sdk.Wrapf("hashLock is required")
	}

	sender, err := hc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgRefundHTLC{
		Sender:   sender.String(),
		HashLock: hashLock,
	}
	return hc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (hc htlcClient) QueryHTLC(hashLock string) (QueryHTLCResp, sdk.Error) {
	if len(hashLock) == 0 {
		return QueryHTLCResp{}, sdk.Wrapf("hashLock is required")
	}

	conn, err := hc.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryHTLCResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).HTLC(
		context.Background(),
		&QueryHTLCRequest{
			HashLock: hashLock,
		})
	if err != nil {
		return QueryHTLCResp{}, sdk.Wrap(err)
	}
	return res.Htlc.Convert().(QueryHTLCResp), nil
}
