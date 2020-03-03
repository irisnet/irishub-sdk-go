package staking

import (
	"errors"
	"strings"

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

	validator, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return nil, err
	}

	msg := MsgDelegate{
		DelegatorAddr: delegator,
		ValidatorAddr: validator,
		Delegation:    amount,
	}

	s.Info().Str("delegator", delegator.String()).
		Str("validator", validator.String()).
		Str("amount", amount.String()).
		Msg("execute delegate transaction")
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

	s.Info().Str("delegator", delegator.String()).
		Str("validator", valAddr).
		Str("amount", amount.String()).
		Msg("execute undelegate transaction")
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

	s.Info().Str("delegator", delAddr.String()).
		Str("srcValidatorAddr", srcValidatorAddr).
		Str("dstValidatorAddr", dstValidatorAddr).
		Str("amount", amount.String()).
		Msg("execute redelegate transaction")
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// QueryDelegation return the specified delegation by delegatorAddr and validatorAddr
func (s stakingClient) QueryDelegation(delegatorAddr, validatorAddr string) (sdk.Delegation, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return sdk.Delegation{}, err
	}

	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return sdk.Delegation{}, err
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
		ValidatorAddr sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	var delegation Delegation
	err = s.Query("custom/stake/delegation", param, &delegation)
	if err != nil {
		return sdk.Delegation{}, err
	}
	return delegation.ToSDKResponse(), err
}

// QueryDelegations return the specified delegations by delegatorAddr
func (s stakingClient) QueryDelegations(delegatorAddr string) (sdk.Delegations, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return sdk.Delegations{}, err
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var delegations Delegations
	err = s.Query("custom/stake/delegatorDelegations", param, &delegations)
	if err != nil {
		return sdk.Delegations{}, err
	}
	return delegations.ToSDKResponse(), err
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
func (s stakingClient) QueryUnbondingDelegations(delegatorAddr string) (sdk.UnbondingDelegations, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return sdk.UnbondingDelegations{}, err
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var unbondings UnbondingDelegations
	err = s.Query("custom/stake/delegatorUnbondingDelegations", param, &unbondings)
	if err != nil {
		return sdk.UnbondingDelegations{}, err
	}
	return unbondings.ToSDKResponse(), err
}

// QueryRedelegation return the specified redelegation by delegatorAddr,srcValidatorAddr,dstValidatorAddr
func (s stakingClient) QueryRedelegation(delegatorAddr, srcValidatorAddr, dstValidatorAddr string) (sdk.Redelegation, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return sdk.Redelegation{}, err
	}

	srcVarAddr, err := sdk.ValAddressFromBech32(srcValidatorAddr)
	if err != nil {
		return sdk.Redelegation{}, err
	}

	dstVarAddr, err := sdk.ValAddressFromBech32(dstValidatorAddr)
	if err != nil {
		return sdk.Redelegation{}, err
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
		return sdk.Redelegation{}, err
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
func (s stakingClient) QueryDelegationsTo(validatorAddr string) (sdk.Delegations, error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return sdk.Delegations{}, err
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var delegations Delegations
	err = s.Query("custom/stake/validatorDelegations", param, &delegations)
	if err != nil {
		return sdk.Delegations{}, err
	}
	return delegations.ToSDKResponse(), nil
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

//
func (s stakingClient) SubscribeValidatorInfoUpdates(validator string,
	callback func(data sdk.EventDataMsgEditValidator)) sdk.Subscription {
	var builder = sdk.NewEventQueryBuilder().AddCondition(sdk.ActionKey,
		"edit_validator")

	s.Info().Str("validator", validator).Msg("subscribe validator update event")
	validator = strings.TrimSpace(validator)
	if len(validator) != 0 {
		builder.AddCondition("destination-validator",
			sdk.EventValue(validator))
	}
	subscription, err := s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		for _, msg := range tx.Tx.Msgs {
			msg, ok := msg.(MsgEditValidator)
			if ok && validator == msg.ValidatorAddr.String() {
				data := sdk.EventDataMsgEditValidator{
					Height: tx.Height,
					Hash:   tx.Hash,
					Description: sdk.Description{
						Moniker:  msg.Moniker,
						Identity: msg.Identity,
						Website:  msg.Website,
						Details:  msg.Details,
					},
					Address:        msg.ValidatorAddr.String(),
					CommissionRate: msg.CommissionRate.String(),
				}
				callback(data)
			}
		}
	})
	s.Err(err).Msg("subscribe validator update event failed")
	return subscription
}
