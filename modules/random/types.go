package random

import (
	"fmt"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "random"
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
		return fmt.Errorf("the consumer address must be specified")
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

func (r Random) Convert() interface{} {
	return QueryRandomResp{
		RequestTxHash: r.RequestTxHash.String(),
		Height:        r.Height,
		Value:         r.Value,
	}
}
