//**
// Package staking provides staking functionalities for validators and delegators
//
// [More Details](https://www.irisnet.org/docs/features/stake.html)
//
package staking

import (
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

func Create(ac sdk.AbstractClient) rpc.Staking {
	return stakingClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

//Delegate is responsible for delegating liquid tokens to an validator
func (s stakingClient) Delegate(valAddr string, amount sdk.Coin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegator, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	validator, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
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
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//Undelegate is responsible for undelegating from a validator
func (s stakingClient) Undelegate(valAddr string, amount sdk.Coin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegator, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	val, err := s.QueryValidator(valAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	exRate := val.DelegatorShareExRate()
	if exRate.IsZero() {
		return sdk.ResultTx{}, sdk.Wrapf("zero exRate should not happen")
	}
	amountDec := sdk.NewDecFromInt(amount.Amount)
	share := amountDec.Quo(exRate)

	varAddr, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
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
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//Redelegate is responsible for redelegating illiquid tokens from one validator to another
func (s stakingClient) Redelegate(srcValidatorAddr,
	dstValidatorAddr string, amount sdk.Coin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delAddr, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	srcValAddr, err := sdk.ValAddressFromBech32(srcValidatorAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	dstValAddr, err := sdk.ValAddressFromBech32(dstValidatorAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	val, err := s.QueryValidator(srcValidatorAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	exRate := val.DelegatorShareExRate()
	if exRate.IsZero() {
		return sdk.ResultTx{}, sdk.Wrapf("zero exRate should not happen")
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
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// QueryDelegation return the specified delegation by delegatorAddr and validatorAddr
func (s stakingClient) QueryDelegation(delegatorAddr, validatorAddr string) (rpc.Delegation, sdk.Error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Delegation{}, sdk.Wrap(err)
	}

	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.Delegation{}, sdk.Wrap(err)
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
		ValidatorAddr sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	var delegation delegation
	if err := s.QueryWithResponse("custom/stake/delegation", param, &delegation); err != nil {
		return rpc.Delegation{}, sdk.Wrap(err)
	}
	return delegation.Convert().(rpc.Delegation), nil
}

// QueryDelegations return the specified delegations by delegatorAddr
func (s stakingClient) QueryDelegations(delegatorAddr string) (rpc.Delegations, sdk.Error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Delegations{}, sdk.Wrap(err)
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var ds delegations
	if err := s.QueryWithResponse("custom/stake/delegatorDelegations", param, &ds); err != nil {
		return rpc.Delegations{}, sdk.Wrap(err)
	}
	return ds.Convert().(rpc.Delegations), nil
}

// QueryUnbondingDelegation return the specified unbonding delegation by delegatorAddr and validatorAddr
func (s stakingClient) QueryUnbondingDelegation(delegatorAddr, validatorAddr string) (rpc.UnbondingDelegation, sdk.Error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.UnbondingDelegation{}, sdk.Wrap(err)
	}

	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.UnbondingDelegation{}, sdk.Wrap(err)
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
		ValidatorAddr sdk.ValAddress
	}{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	var ubd unbondingDelegation
	if err := s.QueryWithResponse("custom/stake/unbondingDelegation", param, &ubd); err != nil {
		return rpc.UnbondingDelegation{}, sdk.Wrap(err)
	}
	return ubd.Convert().(rpc.UnbondingDelegation), nil
}

// QueryUnbondingDelegations return the specified unbonding delegations by delegatorAddr
func (s stakingClient) QueryUnbondingDelegations(delegatorAddr string) (rpc.UnbondingDelegations, sdk.Error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.UnbondingDelegations{}, sdk.Wrap(err)
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var unds unbondingDelegations
	if err := s.QueryWithResponse("custom/stake/delegatorUnbondingDelegations", param, &unds); err != nil {
		return rpc.UnbondingDelegations{}, sdk.Wrap(err)
	}
	return unds.Convert().(rpc.UnbondingDelegations), nil
}

// QueryRedelegation return the specified redelegation by delegatorAddr,srcValidatorAddr,dstValidatorAddr
func (s stakingClient) QueryRedelegation(delegatorAddr, srcValidatorAddr, dstValidatorAddr string) (rpc.Redelegation, sdk.Error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Redelegation{}, sdk.Wrap(err)
	}

	srcVarAddr, err := sdk.ValAddressFromBech32(srcValidatorAddr)
	if err != nil {
		return rpc.Redelegation{}, sdk.Wrap(err)
	}

	dstVarAddr, err := sdk.ValAddressFromBech32(dstValidatorAddr)
	if err != nil {
		return rpc.Redelegation{}, sdk.Wrap(err)
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
	if err := s.QueryWithResponse("custom/stake/redelegation", param, &rd); err != nil {
		return rpc.Redelegation{}, sdk.Wrap(err)
	}
	return rd.Convert().(rpc.Redelegation), nil
}

// QueryRedelegations return the specified redelegations by delegatorAddr
func (s stakingClient) QueryRedelegations(delegatorAddr string) (rpc.Redelegations, sdk.Error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.Redelegations{}, sdk.Wrap(err)
	}
	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var rds redelegations
	if err := s.QueryWithResponse("custom/stake/delegatorRedelegations", param, &rds); err != nil {
		return rpc.Redelegations{}, sdk.Wrap(err)
	}
	return rds.Convert().(rpc.Redelegations), nil
}

// QueryDelegationsTo return the specified delegations by validatorAddr
func (s stakingClient) QueryDelegationsTo(validatorAddr string) (rpc.Delegations, sdk.Error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.Delegations{}, sdk.Wrap(err)
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var ds delegations
	if err := s.QueryWithResponse("custom/stake/validatorDelegations", param, &ds); err != nil {
		return rpc.Delegations{}, sdk.Wrap(err)
	}
	return ds.Convert().(rpc.Delegations), nil
}

// QueryUnbondingDelegationsFrom return the specified unbonding delegations by validatorAddr
func (s stakingClient) QueryUnbondingDelegationsFrom(validatorAddr string) (rpc.UnbondingDelegations, sdk.Error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.UnbondingDelegations{}, sdk.Wrap(err)
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var ubds unbondingDelegations
	if err := s.QueryWithResponse("custom/stake/validatorUnbondingDelegations", param, &ubds); err != nil {
		return rpc.UnbondingDelegations{}, sdk.Wrap(err)
	}
	return ubds.Convert().(rpc.UnbondingDelegations), nil
}

// QueryRedelegationsFrom return the specified redelegations by validatorAddr
func (s stakingClient) QueryRedelegationsFrom(validatorAddr string) (rpc.Redelegations, sdk.Error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.Redelegations{}, sdk.Wrap(err)
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var rds redelegations
	if err := s.QueryWithResponse("custom/stake/validatorRedelegations", param, &rds); err != nil {
		return rpc.Redelegations{}, sdk.Wrap(err)
	}
	return rds.Convert().(rpc.Redelegations), nil
}

// QueryValidator return the specified validator by validator address
func (s stakingClient) QueryValidator(address string) (rpc.Validator, sdk.Error) {
	varAddr, err := sdk.ValAddressFromBech32(address)
	if err != nil {
		return rpc.Validator{}, sdk.Wrap(err)
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
	}{
		ValidatorAddr: varAddr,
	}

	var validator validator
	if err := s.QueryWithResponse("custom/stake/validator", param, &validator); err != nil {
		return rpc.Validator{}, sdk.Wrap(err)
	}
	return validator.Convert().(rpc.Validator), nil
}

// QueryValidators return the specified validators by page and size
func (s stakingClient) QueryValidators(page uint64, size uint16) (rpc.Validators, sdk.Error) {
	param := struct {
		Page uint64
		Size uint16
	}{
		Page: page,
		Size: size,
	}

	var validators validators
	if err := s.QueryWithResponse("custom/stake/validators", param, &validators); err != nil {
		return rpc.Validators{}, sdk.Wrap(err)
	}
	return validators.Convert().(rpc.Validators), nil
}

// QueryValidators return the staking pool status
func (s stakingClient) QueryPool() (rpc.StakePool, sdk.Error) {
	var pool Pool
	if err := s.QueryWithResponse("custom/stake/pool", nil, &pool); err != nil {
		return rpc.StakePool{}, sdk.Wrap(err)
	}
	return pool.Convert().(rpc.StakePool), nil
}

// QueryValidators return the staking gov params
func (s stakingClient) QueryParams() (rpc.StakeParams, sdk.Error) {
	var params params
	if err := s.QueryWithResponse("custom/stake/parameters", nil, &params); err != nil {
		return rpc.StakeParams{}, sdk.Wrap(err)
	}
	return params.Convert().(rpc.StakeParams), nil
}

//
func (s stakingClient) SubscribeValidatorInfoUpdates(validator string,
	callback func(data rpc.EventDataMsgEditValidator)) sdk.Subscription {
	var builder = sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(sdk.ActionKey).Equal("edit_validator"))

	s.Info().Str("validator", validator).Msg("subscribe validator update event")
	validator = strings.TrimSpace(validator)
	if len(validator) != 0 {
		builder.AddCondition(sdk.Cond("destination-validator").Equal(sdk.EventValue(validator)))
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
