package staking

import (
	"errors"

	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type stakingClient struct {
	sdk.AbstractClient
	*log.Logger
}

func New(ac sdk.AbstractClient) sdk.Staking {
	return stakingClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With("staking"),
	}
}

//Delegate is responsible for delegating liquid tokens to an validator
func (s stakingClient) Delegate(valAddr string, amount sdk.Coin, baseTx sdk.BaseTx) (sdk.Result, error) {
	delegator, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	varAddr, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return nil, err
	}

	msg := MsgDelegate{
		DelegatorAddr: delegator,
		ValidatorAddr: varAddr,
		Delegation:    amount,
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//Undelegate is responsible for undelegating from a validator
func (s stakingClient) Undelegate(valAddr string, amount sdk.Coin, baseTx sdk.BaseTx) (sdk.Result, error) {
	delegator, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	val, err := s.QueryValidator(valAddr)
	if err != nil {
		return nil, err
	}
	exRate := val.DelegatorShareExRate()
	if exRate.IsZero() {
		return nil, errors.New("zero exRate should not happen")
	}
	amountDec := sdk.NewDecFromInt(amount.Amount)
	share := amountDec.Quo(exRate)

	varAddr, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return nil, err
	}

	msg := MsgUndelegate{
		DelegatorAddr: delegator,
		ValidatorAddr: varAddr,
		SharesAmount:  share,
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//Redelegate is responsible for redelegating illiquid tokens from one validator to another
func (s stakingClient) Redelegate(srcValidatorAddr,
	dstValidatorAddr string, amount sdk.Coin, baseTx sdk.BaseTx) (sdk.Result, error) {
	delAddr, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	srcValAddr, err := sdk.ValAddressFromBech32(srcValidatorAddr)
	if err != nil {
		return nil, err
	}

	dstValAddr, err := sdk.ValAddressFromBech32(dstValidatorAddr)
	if err != nil {
		return nil, err
	}

	val, err := s.QueryValidator(srcValidatorAddr)
	if err != nil {
		return nil, err
	}

	exRate := val.DelegatorShareExRate()
	if exRate.IsZero() {
		return nil, errors.New("zero exRate should not happen")
	}
	amountDec := sdk.NewDecFromInt(amount.Amount)
	share := amountDec.Quo(exRate)

	msg := MsgBeginRedelegate{
		DelegatorAddr:    delAddr,
		ValidatorSrcAddr: srcValAddr,
		ValidatorDstAddr: dstValAddr,
		SharesAmount:     share,
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// QueryDelegation return the specified delegation by delegatorAddr and validatorAddr
func (s stakingClient) QueryDelegation(delegatorAddr, validatorAddr string) (delegation sdk.Delegation, err error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return delegation, err
	}

	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return delegation, err
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
		ValidatorAddr sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	var del Delegation
	err = s.Query("custom/stake/delegation", param, &del)
	if err != nil {
		return delegation, err
	}
	return del.ToSDKResponse(), err
}

// QueryDelegations return the specified delegations by delegatorAddr
func (s stakingClient) QueryDelegations(delegatorAddr string) (delegations sdk.Delegations, err error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return delegations, err
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var del Delegations
	err = s.Query("custom/stake/delegatorDelegations", param, &del)
	if err != nil {
		return delegations, err
	}
	return del.ToSDKResponse(), err
}

// QueryUnbondingDelegation return the specified unbonding delegation by delegatorAddr and validatorAddr
func (s stakingClient) QueryUnbondingDelegation(delegatorAddr, validatorAddr string) (ubd sdk.UnbondingDelegation, err error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return ubd, err
	}
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return ubd, err
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
		ValidatorAddr sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	var unbonding UnbondingDelegation
	err = s.Query("custom/stake/unbondingDelegation", param, &unbonding)
	if err != nil {
		return ubd, err
	}
	return unbonding.ToSDKResponse(), err
}

// QueryUnbondingDelegations return the specified unbonding delegations by delegatorAddr
func (s stakingClient) QueryUnbondingDelegations(delegatorAddr string) (ubds sdk.UnbondingDelegations, err error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return ubds, err
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var unbondings UnbondingDelegations
	err = s.Query("custom/stake/delegatorUnbondingDelegations", param, &unbondings)
	if err != nil {
		return ubds, err
	}
	return unbondings.ToSDKResponse(), err
}

// QueryRedelegation return the specified redelegation by delegatorAddr,srcValidatorAddr,dstValidatorAddr
func (s stakingClient) QueryRedelegation(delegatorAddr, srcValidatorAddr, dstValidatorAddr string) (rd sdk.Redelegation, err error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rd, err
	}
	srcVarAddr, err := sdk.ValAddressFromBech32(srcValidatorAddr)
	if err != nil {
		return rd, err
	}

	dstVarAddr, err := sdk.ValAddressFromBech32(dstValidatorAddr)
	if err != nil {
		return rd, err
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
		ValSrcAddr    sdk.ValAddress
		ValDstAddr    sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValSrcAddr:    srcVarAddr,
		ValDstAddr:    dstVarAddr,
	}

	var redelegation Redelegation
	err = s.Query("custom/stake/redelegation", param, &redelegation)
	if err != nil {
		return rd, err
	}
	return redelegation.ToSDKResponse(), nil
}

// QueryRedelegations return the specified redelegations by delegatorAddr
func (s stakingClient) QueryRedelegations(delegatorAddr string) (sdk.Redelegations, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return sdk.Redelegations{}, err
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var rds Redelegations
	err = s.Query("custom/stake/delegatorRedelegations", param, &rds)
	if err != nil {
		return sdk.Redelegations{}, err
	}
	return rds.ToSDKResponse(), nil
}

// QueryDelegationsTo return the specified delegations by validatorAddr
func (s stakingClient) QueryDelegationsTo(validatorAddr string) (delegations sdk.Delegations, err error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return delegations, err
	}
	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var ds Delegations
	err = s.Query("custom/stake/validatorDelegations", param, &ds)
	if err != nil {
		return delegations, err
	}
	return ds.ToSDKResponse(), nil
}

// QueryUnbondingDelegationsFrom return the specified unbonding delegations by validatorAddr
func (s stakingClient) QueryUnbondingDelegationsFrom(validatorAddr string) (sdk.UnbondingDelegations, error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return sdk.UnbondingDelegations{}, err
	}
	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var ubds UnbondingDelegations
	err = s.Query("custom/stake/validatorUnbondingDelegations", param, &ubds)
	if err != nil {
		return sdk.UnbondingDelegations{}, err
	}
	return ubds.ToSDKResponse(), nil
}

// QueryRedelegationsFrom return the specified redelegations by validatorAddr
func (s stakingClient) QueryRedelegationsFrom(validatorAddr string) (sdk.Redelegations, error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return sdk.Redelegations{}, err
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var rds Redelegations
	err = s.Query("custom/stake/validatorRedelegations", param, &rds)
	if err != nil {
		return sdk.Redelegations{}, err
	}
	return rds.ToSDKResponse(), nil
}

// QueryValidator return the specified validator by validator address
func (s stakingClient) QueryValidator(address string) (sdk.Validator, error) {
	varAddr, err := sdk.ValAddressFromBech32(address)
	if err != nil {
		return sdk.Validator{}, err
	}
	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var validator Validator
	err = s.Query("custom/stake/validator", param, &validator)
	if err != nil {
		return sdk.Validator{}, err
	}
	return validator.ToSDKResponse(), nil
}

// QueryValidators return the specified validators by page and size
func (s stakingClient) QueryValidators(page uint64, size uint16) (sdk.Validators, error) {
	param := struct {
		Page uint64
		Size uint16
	}{
		Page: page,
		Size: size,
	}

	var validators Validators
	err := s.Query("custom/stake/validators", param, &validators)
	if err != nil {
		return sdk.Validators{}, err
	}
	return validators.ToSDKResponse(), nil
}

// QueryValidators return the staking pool status
func (s stakingClient) QueryPool() (sdk.StakePool, error) {
	var pool Pool
	err := s.Query("custom/stake/pool", nil, &pool)
	if err != nil {
		return sdk.StakePool{}, err
	}
	return pool.ToSDKResponse(), nil
}

// QueryValidators return the staking gov params
func (s stakingClient) QueryParams() (sdk.StakeParams, error) {
	var params Params
	err := s.Query("custom/stake/parameters", nil, &params)
	if err != nil {
		return sdk.StakeParams{}, err
	}
	return params.ToSDKResponse(), nil
}
