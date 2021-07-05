package htlc

import (
	"encoding/hex"

	sdk "github.com/irisnet/core-sdk-go/types"
)

const (
	ModuleName = "htlc"

	SecretLength                    = 64    // length for the secret in bytes
	HashLockLength                  = 64    // length for the hash lock in bytes
	HTLCIDLength                    = 64    // HTLCIDLength is the length for the hash lock in hex string
	MaxLengthForAddressOnOtherChain = 128   // maximum length for the address on other chains
	MinTimeLock                     = 50    // minimum time span for HTLC
	MaxTimeLock                     = 25480 // maximum time span for HTLC
)

var (
	_ sdk.Msg = &MsgCreateHTLC{}
	_ sdk.Msg = &MsgClaimHTLC{}
)

func (msg MsgCreateHTLC) Route() string { return ModuleName }

func (msg MsgCreateHTLC) Type() string { return "create_htlc" }

func (msg MsgCreateHTLC) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdk.Wrapf("invalid sender address (%s)", err)
	}
	if len(msg.To) == 0 {
		return sdk.Wrapf("recipient missing")
	}
	if len(msg.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain {
		return sdk.Wrapf("length of the receiver on other chain must be between [0,%d]", MaxLengthForAddressOnOtherChain)
	}
	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return sdk.Wrapf("the transferred amount must be valid")
	}
	if _, err := hex.DecodeString(msg.HashLock); err != nil {
		return sdk.Wrapf("hash lock must be a hex encoded string")
	}
	if len(msg.HashLock) != HashLockLength {
		return sdk.Wrapf("length of the hash lock must be %d in bytes", HashLockLength)
	}
	if msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock {
		return sdk.Wrapf("the time lock must be between [%d,%d]", MinTimeLock, MaxTimeLock)
	}
	return nil
}

func (msg MsgCreateHTLC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateHTLC) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg MsgClaimHTLC) Route() string { return ModuleName }

func (msg MsgClaimHTLC) Type() string { return "claim_htlc" }

func (msg MsgClaimHTLC) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdk.Wrapf("invalid sender address (%s)", err)
	}
	if _, err := hex.DecodeString(msg.Id); err != nil {
		return sdk.Wrapf("htlc id must be a hex encoded string")
	}
	if len(msg.Id) != HTLCIDLength {
		return sdk.Wrapf("length of the htlc id must be %d in bytes", HashLockLength)
	}
	if _, err := hex.DecodeString(msg.Secret); err != nil {
		return sdk.Wrapf("secret must be a hex encoded string")
	}
	if len(msg.Secret) != SecretLength {
		return sdk.Wrapf("length of the secret must be %d in bytes", SecretLength)
	}
	return nil
}

func (msg MsgClaimHTLC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgClaimHTLC) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (h HTLC) Convert() interface{} {
	return QueryHTLCResp{
		Sender:               h.Sender,
		To:                   h.To,
		ReceiverOnOtherChain: h.ReceiverOnOtherChain,
		SenderOnOtherChain:   h.SenderOnOtherChain,
		Amount:               h.Amount,
		Secret:               h.Secret,
		Timestamp:            h.Timestamp,
		ExpirationHeight:     h.ExpirationHeight,
		State:                int32(h.State),
		Transfer:             h.Transfer,
	}
}

func (h Params) Convert() interface{} {
	var params []AssetParamDto
	for _, val := range h.AssetParams {
		params = append(params, AssetParamDto{
			Denom:         val.Denom,
			Active:        val.Active,
			DeputyAddress: val.DeputyAddress,
			FixedFee:      val.FixedFee.Uint64(),
			MinSwapAmount: val.MinSwapAmount.Uint64(),
			MinBlockLock:  val.MinBlockLock,
			MaxSwapAmount: val.MaxSwapAmount.Uint64(),
			MaxBlockLock:  val.MaxBlockLock,
			SupplyLimit: SupplyLimitDto{
				Limit:          val.SupplyLimit.Limit.Uint64(),
				TimeLimited:    val.SupplyLimit.TimeLimited,
				TimePeriod:     int64(val.SupplyLimit.TimePeriod),
				TimeBasedLimit: val.SupplyLimit.TimeBasedLimit.Uint64(),
			},
		})
	}
	return QueryParamsResp{
		AssetParams: params,
	}
}
