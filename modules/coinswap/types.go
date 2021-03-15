package coinswap

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "coinswap"


	eventTypeTransfer = "transfer"
	eventTypeSwap = "swap"

	attributeKeyAmount  = "amount"
)




var (
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgRemoveLiquidity{}
	_ sdk.Msg = &MsgSwapOrder{}
)

// Route implements Msg.
func (msg MsgAddLiquidity) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgAddLiquidity) Type() string { return "add_liquidity" }

// GetSignBytes implements Msg.
func (msg MsgAddLiquidity) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgAddLiquidity) ValidateBasic() error {
	return nil
}

// GetSigners implements Msg.
func (msg MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}

// Route implements Msg.
func (msg MsgRemoveLiquidity) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgRemoveLiquidity) Type() string { return "remove_liquidity" }

// GetSignBytes implements Msg.
func (msg MsgRemoveLiquidity) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRemoveLiquidity) ValidateBasic() error {
	return nil
}

// GetSigners implements Msg.
func (msg MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}

// Route implements Msg.
func (msg MsgSwapOrder) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgSwapOrder) Type() string { return "swap_order" }

// GetSignBytes implements Msg.
func (msg MsgSwapOrder) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSwapOrder) ValidateBasic() error {
	return nil
}

// GetSigners implements Msg.
func (msg MsgSwapOrder) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Input.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
