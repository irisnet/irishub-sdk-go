package coinswap

import (
	sdk "github.com/irisnet/core-sdk-go/types"
)

const (
	ModuleName = "coinswap"

	eventTypeTransfer = "transfer"
	eventTypeSwap     = "swap"

	attributeKeyAmount = "amount"
)

var (
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgRemoveLiquidity{}
	_ sdk.Msg = &MsgSwapOrder{}
)

type totalSupply = func() (sdk.Coins, sdk.Error)

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
	if !(msg.MaxToken.IsValid() && msg.MaxToken.IsPositive()) {
		return sdk.Wrapf("invalid MaxToken: %s", msg.MaxToken.String())
	}

	if !msg.ExactStandardAmt.IsPositive() {
		return sdk.Wrapf("standard token amount must be positive")
	}

	if msg.MinLiquidity.IsNegative() {
		return sdk.Wrapf("minimum liquidity can not be negative")
	}

	if msg.Deadline <= 0 {
		return sdk.Wrapf("deadline %d must be greater than 0", msg.Deadline)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdk.Wrap(err)
	}
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
	if msg.MinToken.IsNegative() {
		return sdk.Wrapf("minimum token amount can not be negative")
	}
	if !msg.WithdrawLiquidity.IsValid() || !msg.WithdrawLiquidity.IsPositive() {
		return sdk.Wrapf("invalid withdrawLiquidity (%s)", msg.WithdrawLiquidity.String())
	}
	if msg.MinStandardAmt.IsNegative() {
		return sdk.Wrapf("minimum standard token amount %s can not be negative", msg.MinStandardAmt.String())
	}
	if msg.Deadline <= 0 {
		return sdk.Wrapf("deadline %d must be greater than 0", msg.Deadline)
	}
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
	if !(msg.Input.Coin.IsValid() && msg.Input.Coin.IsPositive()) {
		return sdk.Wrapf("invalid input (%s)", msg.Input.Coin.String())
	}

	if _, err := sdk.AccAddressFromBech32(msg.Input.Address); err != nil {
		return sdk.Wrap(err)
	}

	if !(msg.Output.Coin.IsValid() && msg.Output.Coin.IsPositive()) {
		return sdk.Wrapf("invalid output (%s)", msg.Output.Coin.String())
	}

	if _, err := sdk.AccAddressFromBech32(msg.Output.Address); err != nil {
		return sdk.Wrap(err)
	}

	if msg.Input.Coin.Denom == msg.Output.Coin.Denom {
		return sdk.Wrapf("invalid swap")
	}

	if msg.Deadline <= 0 {
		return sdk.Wrapf("deadline %d must be greater than 0", msg.Deadline)
	}
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

func (m QueryLiquidityResponse) Convert() interface{} {
	return &QueryPoolResponse{
		BaseCoin:  m.Standard,
		TokenCoin: m.Token,
		Liquidity: m.Liquidity,
		Fee:       m.Fee,
	}
}
