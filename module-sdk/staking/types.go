package staking

import (
	"bytes"

	"github.com/irisnet/core-sdk-go/common/codec"
	codectypes "github.com/irisnet/core-sdk-go/common/codec/types"
	crypto "github.com/irisnet/core-sdk-go/common/crypto/types"
	sdk "github.com/irisnet/core-sdk-go/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	ModuleName = "staking"
)

var (
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateValidator)(nil)
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ sdk.Msg                            = &MsgEditValidator{}
	_ sdk.Msg                            = &MsgDelegate{}
	_ sdk.Msg                            = &MsgUndelegate{}
	_ sdk.Msg                            = &MsgBeginRedelegate{}
)

// DelegationI delegation bond for a delegated proof of stake system
type DelegationI interface {
	GetDelegatorAddr() sdk.AccAddress // delegator sdk.AccAddress for the bond
	GetValidatorAddr() sdk.ValAddress // validator operator address
	GetShares() sdk.Dec               // amount of validator's shares held in this delegation
}

// ValidatorI expected validator functions
type ValidatorI interface {
	IsJailed() bool                                         // whether the validator is jailed
	GetMoniker() string                                     // moniker of the validator
	GetStatus() BondStatus                                  // status of the validator
	IsBonded() bool                                         // check if has a bonded status
	IsUnbonded() bool                                       // check if has status unbonded
	IsUnbonding() bool                                      // check if has status unbonding
	GetOperator() sdk.ValAddress                            // operator address to receive/return validators coins
	TmConsPubKey() (crypto.PubKey, error)                   // validation consensus pubkey
	GetConsAddr() (sdk.ConsAddress, error)                  // validation consensus address
	GetTokens() sdk.Int                                     // validation tokens
	GetBondedTokens() sdk.Int                               // validator bonded tokens
	GetConsensusPower() int64                               // validation power in tendermint
	GetCommission() sdk.Dec                                 // validator commission rate
	GetMinSelfDelegation() sdk.Int                          // validator minimum self delegation
	GetDelegatorShares() sdk.Dec                            // total outstanding delegator shares
	TokensFromShares(sdk.Dec) sdk.Dec                       // token worth of provided delegator shares
	TokensFromSharesTruncated(sdk.Dec) sdk.Dec              // token worth of provided delegator shares, truncated
	TokensFromSharesRoundUp(sdk.Dec) sdk.Dec                // token worth of provided delegator shares, rounded up
	SharesFromTokens(amt sdk.Int) (sdk.Dec, error)          // shares worth of delegator's bond
	SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, error) // truncated shares worth of delegator's bond
}

func (msg MsgCreateValidator) Route() string { return ModuleName }

func (msg MsgCreateValidator) Type() string { return "create_validator" }

func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(delAddr.Bytes(), addr.Bytes()) {
		addrs = append(addrs, sdk.AccAddress(addr))
	}

	return addrs
}

func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}
	if delAddr.Empty() {
		return sdk.Wrapf("missing delegatorAddr")
	}

	if msg.ValidatorAddress == "" {
		return sdk.Wrapf("missing validatorAddr")
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdk.Wrap(err)
	}
	if !sdk.AccAddress(valAddr).Equals(delAddr) {
		return sdk.Wrapf("validatorAddr must equal delegatorAddr, validatorAddr:[%s], delegatorAddr:[%s]", valAddr, delAddr)
	}

	if msg.Pubkey == nil {
		return sdk.Wrapf("missing validatorPubKey")
	}

	if msg.Description == (Description{}) {
		return sdk.Wrapf("missing description")
	}

	if msg.Commission == (CommissionRates{}) {
		return sdk.Wrapf("missing commission")
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return sdk.Wrapf("minSelfDelegation isn't positive")
	}

	return nil
}

func (msg MsgCreateValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey crypto.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}

func (msg MsgEditValidator) Route() string { return ModuleName }

func (msg MsgEditValidator) Type() string { return "edit_validator" }

func (msg MsgEditValidator) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{valAddr.Bytes()}
}

func (msg MsgEditValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditValidator) ValidateBasic() error {
	if msg.ValidatorAddress == "" {
		return sdk.Wrapf("missing validatorAddress")
	}

	if msg.Description == (Description{}) {
		return sdk.Wrapf("missing description")
	}

	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return sdk.Wrapf("minSelfDelegation isn't positive")
	}

	return nil
}

func (msg MsgDelegate) Route() string { return ModuleName }

func (msg MsgDelegate) Type() string { return "delegate" }

func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return sdk.Wrapf("missing delegatorAddress")
	}

	if msg.ValidatorAddress == "" {
		return sdk.Wrapf("missing errEmptyValidatorAddr")
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdk.Wrapf("amount isn't positive or valid")
	}

	return nil
}

func (msg MsgBeginRedelegate) Route() string { return ModuleName }

func (msg MsgBeginRedelegate) Type() string { return "begin_redelegate" }

func (msg MsgBeginRedelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBeginRedelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return sdk.Wrapf("missing delegatorAddress")
	}

	if msg.ValidatorSrcAddress == "" {
		return sdk.Wrapf("missing validatorSrcAddress")
	}

	if msg.ValidatorDstAddress == "" {
		return sdk.Wrapf("missing validatorDstAddress")
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdk.Wrapf("amount isn't positive or valid")
	}

	return nil
}

func (msg MsgUndelegate) Route() string { return ModuleName }

func (msg MsgUndelegate) Type() string { return "begin_unbonding" }

func (msg MsgUndelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUndelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return sdk.Wrapf("missing delegatorAddr")
	}

	if msg.ValidatorAddress == "" {
		return sdk.Wrapf("missing validatorAddress")
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdk.Wrapf("amount isn't positive or valid")
	}

	return nil
}

func (q QueryValidatorsResponse) Convert(cdc codec.Marshaler) interface{} {
	var validatorResps []QueryValidatorResp
	for _, v := range q.Validators {
		validatorResps = append(validatorResps, v.Convert(cdc).(QueryValidatorResp))
	}

	return QueryValidatorsResp{
		Validators: validatorResps,
		Total:      q.Pagination.Total,
	}
}

func (v Validator) Convert(cdc codec.Marshaler) interface{} {
	pubKey, _ := v.GetPubKey(cdc)
	return QueryValidatorResp{
		OperatorAddress: v.OperatorAddress,
		ConsensusPubkey: pubKey.String(),
		Jailed:          v.Jailed,
		Status:          BondStatus_name[int32(v.Status)],
		Tokens:          v.Tokens,
		DelegatorShares: v.DelegatorShares,
		Description: description{
			Moniker:         v.Description.Moniker,
			Identity:        v.Description.Identity,
			Website:         v.Description.Website,
			SecurityContact: v.Description.SecurityContact,
			Details:         v.Description.Details,
		},
		UnbondingHeight: v.UnbondingHeight,
		UnbondingTime:   v.UnbondingTime,
		Commission: commission{
			commissionRates: commissionRates{
				Rate:          v.Commission.Rate,
				MaxRate:       v.Commission.MaxRate,
				MaxChangeRate: v.Commission.MaxChangeRate,
			},
			UpdateTime: v.Commission.UpdateTime,
		},
		MinSelfDelegation: v.MinSelfDelegation,
	}
}

// GetPubKey - Implements Validator.
func (v Validator) GetPubKey(unpacker codectypes.AnyUnpacker) (pk crypto.PubKey, err error) {
	if v.ConsensusPubkey == nil {
		return nil, nil
	}

	var pubKey crypto.PubKey
	if err = unpacker.UnpackAny(v.ConsensusPubkey, &pubKey); err != nil {
		return nil, err
	}
	return pubKey, nil
}

func (q QueryValidatorDelegationsResponse) Convert() interface{} {
	var delegationResps []QueryDelegationResp
	for _, v := range q.DelegationResponses {
		delegationResps = append(delegationResps, v.Convert().(QueryDelegationResp))
	}

	return QueryValidatorDelegationsResp{
		DelegationResponses: delegationResps,
		Total:               q.Pagination.Total,
	}
}

func (d DelegationResponse) Convert() interface{} {
	return QueryDelegationResp{
		Delegation: delegation{
			DelegatorAddress: d.Delegation.DelegatorAddress,
			ValidatorAddress: d.Delegation.ValidatorAddress,
			Shares:           d.Delegation.Shares,
		},
		Balance: sdk.Coin{
			Denom:  d.Balance.Denom,
			Amount: d.Balance.Amount,
		},
	}
}

func (q QueryValidatorUnbondingDelegationsResponse) Convert() interface{} {
	var unbondingDelegations []QueryUnbondingDelegationResp
	for _, v := range q.UnbondingResponses {
		unbondingDelegations = append(unbondingDelegations, v.Convert().(QueryUnbondingDelegationResp))
	}

	return QueryValidatorUnbondingDelegationsResp{
		UnbondingResponses: unbondingDelegations,
		Total:              q.Pagination.Total,
	}
}

func (u UnbondingDelegation) Convert() interface{} {
	var entries []unbondingDelegationEntry
	for _, v := range u.Entries {
		entries = append(entries, v.Convert().(unbondingDelegationEntry))
	}

	return QueryUnbondingDelegationResp{
		DelegatorAddress: u.DelegatorAddress,
		ValidatorAddress: u.ValidatorAddress,
		Entries:          entries,
	}
}

func (u UnbondingDelegationEntry) Convert() interface{} {
	return unbondingDelegationEntry{
		CreationHeight: u.CreationHeight,
		CompletionTime: u.CompletionTime,
		InitialBalance: u.InitialBalance,
		Balance:        u.Balance,
	}
}

func (q QueryDelegatorDelegationsResponse) Convert() interface{} {
	var delegationResps []QueryDelegationResp
	for _, v := range q.DelegationResponses {
		delegationResps = append(delegationResps, v.Convert().(QueryDelegationResp))
	}

	return QueryDelegatorDelegationsResp{
		DelegationResponses: delegationResps,
		Total:               0,
	}
}

func (q QueryDelegatorUnbondingDelegationsResponse) Convert() interface{} {
	var unbondingDelegations []QueryUnbondingDelegationResp
	for _, v := range q.UnbondingResponses {
		unbondingDelegations = append(unbondingDelegations, v.Convert().(QueryUnbondingDelegationResp))
	}
	return QueryDelegatorUnbondingDelegationsResp{
		UnbondingDelegations: unbondingDelegations,
		Total:                q.Pagination.Total,
	}
}

func (q QueryRedelegationsResponse) Convert() interface{} {
	var redelegationResps []RedelegationResp
	for _, v := range q.RedelegationResponses {
		redelegationResps = append(redelegationResps, v.Convert().(RedelegationResp))
	}

	return QueryRedelegationsResp{
		RedelegationResponses: redelegationResps,
		Total:                 q.Pagination.Total,
	}
}

func (r RedelegationResponse) Convert() interface{} {
	var outerEntries []redelegationEntryResponse
	for _, v := range r.Entries {
		outerEntries = append(outerEntries, v.Convert().(redelegationEntryResponse))
	}

	var innerEntries []redelegationEntry
	for _, v := range r.Redelegation.Entries {
		innerEntries = append(innerEntries, v.Convert().(redelegationEntry))
	}

	return RedelegationResp{
		Redelegation: redelegation{
			DelegatorAddress:    r.Redelegation.DelegatorAddress,
			ValidatorSrcAddress: r.Redelegation.ValidatorSrcAddress,
			ValidatorDstAddress: r.Redelegation.ValidatorDstAddress,
			Entries:             innerEntries,
		},
		Entries: outerEntries,
	}
}

func (r RedelegationEntry) Convert() interface{} {
	return redelegationEntry{
		CreationHeight: r.CreationHeight,
		CompletionTime: r.CompletionTime,
		InitialBalance: r.InitialBalance,
		SharesDst:      r.SharesDst,
	}
}

func (r RedelegationEntryResponse) Convert() interface{} {
	return redelegationEntryResponse{
		RedelegationEntry: redelegationEntry{
			CreationHeight: r.RedelegationEntry.CreationHeight,
			CompletionTime: r.RedelegationEntry.CompletionTime,
			InitialBalance: r.RedelegationEntry.InitialBalance,
			SharesDst:      r.RedelegationEntry.SharesDst,
		},
		Balance: r.Balance,
	}
}

func (q QueryDelegatorValidatorsResponse) Convert(cdc codec.Marshaler) interface{} {
	var validators []QueryValidatorResp
	for _, v := range q.Validators {
		validators = append(validators, v.Convert(cdc).(QueryValidatorResp))
	}

	return QueryDelegatorValidatorsResp{
		Validator: validators,
		Total:     q.Pagination.Total,
	}
}

func (q QueryHistoricalInfoResponse) Convert(cdc codec.Marshaler) interface{} {
	var valset []QueryValidatorResp
	for _, v := range q.Hist.Valset {
		valset = append(valset, v.Convert(cdc).(QueryValidatorResp))
	}

	header := q.Hist.Header
	lastBlockId := q.Hist.Header.LastBlockId.PartSetHeader
	partSetHeader := q.Hist.Header.LastBlockId.PartSetHeader
	return QueryHistoricalInfoResp{
		Header: sdk.Header{
			Version: header.Version,
			ChainID: header.ChainID,
			Height:  header.Height,
			Time:    header.Time,
			LastBlockID: tmtypes.BlockID{
				Hash: lastBlockId.Hash,
				PartSetHeader: tmtypes.PartSetHeader{
					Total: partSetHeader.Total,
					Hash:  partSetHeader.Hash,
				},
			},
			LastCommitHash:     header.LastCommitHash,
			DataHash:           header.DataHash,
			ValidatorsHash:     header.ValidatorsHash,
			NextValidatorsHash: header.NextValidatorsHash,
			ConsensusHash:      header.ConsensusHash,
			AppHash:            header.AppHash,
			LastResultsHash:    header.LastResultsHash,
			EvidenceHash:       header.EvidenceHash,
			ProposerAddress:    header.ProposerAddress,
		},
		Valset: valset,
	}
}

func (q QueryParamsResponse) Convert() interface{} {
	return QueryParamsResp{
		UnbondingTime:     q.Params.UnbondingTime,
		MaxValidators:     q.Params.MaxValidators,
		MaxEntries:        q.Params.MaxEntries,
		HistoricalEntries: q.Params.HistoricalEntries,
		BondDenom:         q.Params.BondDenom,
	}
}
