//**
// Package staking provides staking functionalities for validators and delegators
//
// [More Details](https://www.irisnet.org/docs/features/stake.html)
//
package staking

import (
	"errors"
	"strings"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type stakingClient struct {
	sdk.AbstractClient
	*log.Logger
}

func (s stakingClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (s stakingClient) Name() string {
	return ModuleName
}

func New(ac sdk.AbstractClient) rpc.Staking {
	return stakingClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
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
func (s stakingClient) QueryDelegation(delegatorAddr, validatorAddr string) (rpc.Delegation, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Delegation{}, err
	}

	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.Delegation{}, err
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
		ValidatorAddr sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	var delegation delegation
	if err = s.QueryWithResponse("custom/stake/delegation", param, &delegation); err != nil {
		return rpc.Delegation{}, err
	}
	return delegation.Convert().(rpc.Delegation), err
}

// QueryDelegations return the specified delegations by delegatorAddr
func (s stakingClient) QueryDelegations(delegatorAddr string) (rpc.Delegations, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Delegations{}, err
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var ds Delegations
	if err = s.QueryWithResponse("custom/stake/delegatorDelegations", param, &ds); err != nil {
		return rpc.Delegations{}, err
	}
	return ds.Convert().(rpc.Delegations), err
}

// QueryUnbondingDelegation return the specified unbonding delegation by delegatorAddr and validatorAddr
func (s stakingClient) QueryUnbondingDelegation(delegatorAddr, validatorAddr string) (rpc.UnbondingDelegation, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.UnbondingDelegation{}, err
	}

	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.UnbondingDelegation{}, err
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
		ValidatorAddr sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	var ubd unbondingDelegation
	if err = s.QueryWithResponse("custom/stake/unbondingDelegation", param, &ubd); err != nil {
		return rpc.UnbondingDelegation{}, err
	}
	return ubd.Convert().(rpc.UnbondingDelegation), err
}

// QueryUnbondingDelegations return the specified unbonding delegations by delegatorAddr
func (s stakingClient) QueryUnbondingDelegations(delegatorAddr string) (rpc.UnbondingDelegations, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.UnbondingDelegations{}, err
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var unds UnbondingDelegations
	err = s.QueryWithResponse("custom/stake/delegatorUnbondingDelegations", param, &unds)
	if err != nil {
		return rpc.UnbondingDelegations{}, err
	}
	return unds.Convert().(rpc.UnbondingDelegations), err
}

// QueryRedelegation return the specified redelegation by delegatorAddr,srcValidatorAddr,dstValidatorAddr
func (s stakingClient) QueryRedelegation(delegatorAddr, srcValidatorAddr, dstValidatorAddr string) (rpc.Redelegation, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Redelegation{}, err
	}

	srcVarAddr, err := sdk.ValAddressFromBech32(srcValidatorAddr)
	if err != nil {
		return rpc.Redelegation{}, err
	}

	dstVarAddr, err := sdk.ValAddressFromBech32(dstValidatorAddr)
	if err != nil {
		return rpc.Redelegation{}, err
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

	var rd redelegation
	if err = s.QueryWithResponse("custom/stake/redelegation", param, &rd); err != nil {
		return rpc.Redelegation{}, err
	}
	return rd.Convert().(rpc.Redelegation), nil
}

// QueryRedelegations return the specified redelegations by delegatorAddr
func (s stakingClient) QueryRedelegations(delegatorAddr string) (rpc.Redelegations, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Redelegations{}, err
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var rds Redelegations
	if err = s.QueryWithResponse("custom/stake/delegatorRedelegations", param, &rds); err != nil {
		return rpc.Redelegations{}, err
	}
	return rds.Convert().(rpc.Redelegations), nil
}

// QueryDelegationsTo return the specified delegations by validatorAddr
func (s stakingClient) QueryDelegationsTo(validatorAddr string) (rpc.Delegations, error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.Delegations{}, err
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var ds Delegations
	if err = s.QueryWithResponse("custom/stake/validatorDelegations", param, &ds); err != nil {
		return rpc.Delegations{}, err
	}
	return ds.Convert().(rpc.Delegations), nil
}

// QueryUnbondingDelegationsFrom return the specified unbonding delegations by validatorAddr
func (s stakingClient) QueryUnbondingDelegationsFrom(validatorAddr string) (rpc.UnbondingDelegations, error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.UnbondingDelegations{}, err
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var ubds UnbondingDelegations
	if err = s.QueryWithResponse("custom/stake/validatorUnbondingDelegations", param, &ubds); err != nil {
		return rpc.UnbondingDelegations{}, err
	}
	return ubds.Convert().(rpc.UnbondingDelegations), nil
}

// QueryRedelegationsFrom return the specified redelegations by validatorAddr
func (s stakingClient) QueryRedelegationsFrom(validatorAddr string) (rpc.Redelegations, error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.Redelegations{}, err
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var rds Redelegations
	if err = s.QueryWithResponse("custom/stake/validatorRedelegations", param, &rds); err != nil {
		return rpc.Redelegations{}, err
	}
	return rds.Convert().(rpc.Redelegations), nil
}

// QueryValidator return the specified validator by validator address
func (s stakingClient) QueryValidator(address string) (rpc.Validator, error) {
	varAddr, err := sdk.ValAddressFromBech32(address)
	if err != nil {
		return rpc.Validator{}, err
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var validator validator
	if err = s.QueryWithResponse("custom/stake/validator", param, &validator); err != nil {
		return rpc.Validator{}, err
	}
	return validator.Convert().(rpc.Validator), nil
}

// QueryValidators return the specified validators by page and size
func (s stakingClient) QueryValidators(page uint64, size uint16) (rpc.Validators, error) {
	param := struct {
		Page uint64
		Size uint16
	}{
		Page: page,
		Size: size,
	}

	var validators Validators
	if err := s.QueryWithResponse("custom/stake/validators", param, &validators); err != nil {
		return rpc.Validators{}, err
	}
	return validators.Convert().(rpc.Validators), nil
}

// QueryValidators return the staking pool status
func (s stakingClient) QueryPool() (rpc.StakePool, error) {
	var pool Pool
	if err := s.QueryWithResponse("custom/stake/pool", nil, &pool); err != nil {
		return rpc.StakePool{}, err
	}
	return pool.Convert().(rpc.StakePool), nil
}

// QueryValidators return the staking gov params
func (s stakingClient) QueryParams() (rpc.StakeParams, error) {
	var params params
	if err := s.QueryWithResponse("custom/stake/parameters", nil, &params); err != nil {
		return rpc.StakeParams{}, err
	}
	return params.Convert().(rpc.StakeParams), nil
}

//
func (s stakingClient) SubscribeValidatorInfoUpdates(validator string,
	callback func(data rpc.EventDataMsgEditValidator)) sdk.Subscription {
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
				data := rpc.EventDataMsgEditValidator{
					Height: tx.Height,
					Hash:   tx.Hash,
					Description: rpc.Description{
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
