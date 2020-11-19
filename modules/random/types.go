package random

import (
	"errors"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "random"

	eventTypeRequestRequestRandom = "request_random"
	attributeKeyRequestID         = "request_id"
	attributeKeyGenerateHeight    = "generate_height"
)

var (
	_ sdk.Msg = &MsgRequestRandom{}
)

// Route implements Msg.
func (msg MsgRequestRandom) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgRequestRandom) Type() string { return "request_rand" }

// ValidateBasic implements Msg.
func (msg MsgRequestRandom) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("the consumer address must be specified")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgRequestRandom) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg.
func (msg MsgRequestRandom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

func (m Random) Convert() interface{} {
	return QueryRandomResp{
		RequestTxHash: m.RequestTxHash.String(),
		Height:        m.Height,
		Value:         m.Value,
	}
}

type Requests []Request

func (m Requests) Convert() interface{} {
	var res []QueryRandomRequestQueueResp

	for _, request := range m {
		q := QueryRandomRequestQueueResp{
			Height:           request.Height,
			Consumer:         request.Consumer.String(),
			TxHash:           request.TxHash.String(),
			Oracle:           request.Oracle,
			ServiceFeeCap:    request.ServiceFeeCap,
			ServiceContextID: request.ServiceContextID.String(),
		}
		res = append(res, q)
	}
	return res
}
