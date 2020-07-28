package staking

import (
	"encoding/json"
	"strings"

	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type stakingClient struct {
	sdk.BaseClient
	*log.Logger
}

func (s stakingClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (s stakingClient) Name() string {
	return ModuleName
}

func Create(ac sdk.BaseClient) rpc.Staking {
	return stakingClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

//Delegate is responsible for delegating liquid tokens to an validator
func (s stakingClient) Delegate(valAddr string, amount sdk.DecCoin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegator, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	validator, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	//amt, err := s.ToMinCoin(amount)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgDelegate{
		DelegatorAddr: delegator,
		ValidatorAddr: validator,
		//Delegation:    amt[0],
	}

	s.Info().Str("delegator", delegator.String()).
		Str("validator", validator.String()).
		Str("amount", amount.String()).
		Msg("execute delegate transaction")
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//Undelegate is responsible for undelegating from a validator
func (s stakingClient) Undelegate(valAddr string, amount sdk.DecCoin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegator, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	val, err := s.QueryValidator(valAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	//amt, err := s.ToMinCoin(amount)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	exRate := val.DelegatorShareExRate()
	if exRate.IsZero() {
		return sdk.ResultTx{}, sdk.Wrapf("zero exRate should not happen")
	}
	//amountDec := sdk.NewDecFromInt(amt[0].Amount)
	//share := amountDec.Quo(exRate)

	varAddr, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgUndelegate{
		DelegatorAddr: delegator,
		ValidatorAddr: varAddr,
		//SharesAmount:  share,
	}

	s.Info().Str("delegator", delegator.String()).
		Str("validator", valAddr).
		Str("amount", amount.String()).
		Msg("execute undelegate transaction")
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//Redelegate is responsible for redelegating illiquid tokens from one validator to another
func (s stakingClient) Redelegate(srcValidatorAddr,
	dstValidatorAddr string, amount sdk.DecCoin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delAddr, err := s.QueryAddress(baseTx.From)
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

	//amt, err := s.ToMinCoin(amount)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	exRate := val.DelegatorShareExRate()
	if exRate.IsZero() {
		return sdk.ResultTx{}, sdk.Wrapf("zero exRate should not happen")
	}
	//amountDec := sdk.NewDecFromInt(amt[0].Amount)
	//share := amountDec.Quo(exRate)

	msg := MsgBeginRedelegate{
		DelegatorAddr:    delAddr,
		ValidatorSrcAddr: srcValAddr,
		ValidatorDstAddr: dstValAddr,
		//SharesAmount:     share,
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

	var delegationResponse delegationResponse
	if err := s.QueryWithResponse("custom/staking/delegation", param, &delegationResponse); err != nil {
		return rpc.Delegation{}, sdk.Wrap(err)
	}
	return delegationResponse.Convert().(rpc.Delegation), nil
}

// QueryDelegations return the specified delegations by delegatorAddr
func (s stakingClient) QueryDelegations(delegatorAddr string) (rpc.DelegationResponses, sdk.Error) {
	delAddr, err := sdk.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rpc.DelegationResponses{}, sdk.Wrap(err)
	}

	param := struct {
		DelegatorAddr sdk.AccAddress
	}{
		DelegatorAddr: delAddr,
	}

	var ds delegationResponses
	if err := s.QueryWithResponse("custom/staking/delegatorDelegations", param, &ds); err != nil {
		return rpc.DelegationResponses{}, sdk.Wrap(err)
	}
	return ds.Convert().(rpc.DelegationResponses), nil
}

// QueryDelegationsTo return the specified delegations by validatorAddr
func (s stakingClient) QueryDelegationsTo(validatorAddr string) (rpc.DelegationResponses, sdk.Error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.DelegationResponses{}, sdk.Wrap(err)
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
		Page          int
	}{
		ValidatorAddr: varAddr,
		Page:          defaultPage, // A page number must be passed in
	}

	var ds delegationResponses
	if err := s.QueryWithResponse("custom/staking/validatorDelegations", param, &ds); err != nil {
		return rpc.DelegationResponses{}, sdk.Wrap(err)
	}
	return ds.Convert().(rpc.DelegationResponses), nil
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
	if err := s.QueryWithResponse("custom/staking/delegatorUnbondingDelegations", param, &unds); err != nil {
		return rpc.UnbondingDelegations{}, sdk.Wrap(err)
	}
	return unds.Convert().(rpc.UnbondingDelegations), nil
}

// QueryUnbondingDelegationsFrom return the specified unbonding delegations by validatorAddr
func (s stakingClient) QueryUnbondingDelegationsFrom(validatorAddr string) (rpc.UnbondingDelegations, sdk.Error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.UnbondingDelegations{}, sdk.Wrap(err)
	}

	param := struct {
		ValidatorAddr sdk.ValAddress
		Page          int
	}{
		ValidatorAddr: varAddr,
		Page:          1, // A page number must be passed in
	}

	var ubds unbondingDelegations
	if err := s.QueryWithResponse("custom/staking/validatorUnbondingDelegations", param, &ubds); err != nil {
		return rpc.UnbondingDelegations{}, sdk.Wrap(err)
	}
	return ubds.Convert().(rpc.UnbondingDelegations), nil
}

// QueryRedelegationsFrom return the specified redelegations by validatorAddr
func (s stakingClient) QueryRedelegationsFrom(validatorAddr string) (rpc.RedelegationResponses, sdk.Error) {
	varAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rpc.RedelegationResponses{}, sdk.Wrap(err)
	}

	param := struct {
		SrcValidatorAddr sdk.ValAddress
	}{
		SrcValidatorAddr: varAddr,
	}

	var rds redelegationResponses
	if err := s.QueryWithResponse("custom/staking/redelegations", param, &rds); err != nil {
		return rpc.RedelegationResponses{}, sdk.Wrap(err)
	}
	return rds.Convert().(rpc.RedelegationResponses), nil
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
	res, err1 := s.Query("custom/staking/validator", param)
	if err1 != nil {
		return rpc.Validator{}, sdk.Wrap(err)
	}
	if err1 = json.Unmarshal(res, &validator); err1 != nil {
		return rpc.Validator{}, sdk.Wrap(err)
	}
	return validator.Convert().(rpc.Validator), nil
}

// QueryValidators return the specified validators by page and size
func (s stakingClient) QueryValidators(page, size int) (rpc.Validators, sdk.Error) {
	var statuses = []string{BondStatusUnbonded, BondStatusUnbonding, BondStatusBonded}
	var result rpc.Validators
	for _, status := range statuses {
		validators, err := s.queryValidators(page, size, status)
		if err != nil {
			return rpc.Validators{}, sdk.Wrap(err)
		}
		result = append(result, validators...)
	}
	return result, nil
}

// queryValidators return the specified validators by status
func (s stakingClient) queryValidators(page, size int, status string) (rpc.Validators, sdk.Error) {
	param := struct {
		Page, Limit int
		Status      string
	}{
		Page:   page,
		Limit:  size,
		Status: status,
	}

	var validators validators
	if err := s.QueryWithResponse("custom/staking/validators", param, &validators); err != nil {
		return rpc.Validators{}, sdk.Wrap(err)
	}
	return validators.Convert().(rpc.Validators), nil
}

// QueryValidators return the staking pool status
func (s stakingClient) QueryPool() (rpc.StakePool, sdk.Error) {
	var pool Pool
	res, err := s.Query("custom/staking/pool", nil)
	if err != nil {
		return rpc.StakePool{}, sdk.Wrap(err)
	}
	if err := json.Unmarshal(res, &pool); err != nil {
		return rpc.StakePool{}, sdk.Wrap(err)
	}
	return pool.Convert().(rpc.StakePool), nil
}

// QueryValidators return the staking gov params
func (s stakingClient) QueryParams() (rpc.StakeParams, sdk.Error) {
	var params params
	if err := s.BaseClient.QueryParams(s.Name(), &params); err != nil {
		return rpc.StakeParams{}, sdk.Wrap(err)
	}
	return params.Convert().(rpc.StakeParams), nil
}

//
func (s stakingClient) SubscribeValidatorInfoUpdates(validator string,
	callback func(data rpc.EventDataMsgEditValidator)) (sdk.Subscription, sdk.Error) {
	var builder = sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(sdk.ActionKey).EQ("edit_validator"))

	s.Info().Str("validator", validator).Msg("subscribe validator update event")
	validator = strings.TrimSpace(validator)
	if len(validator) != 0 {
		builder.AddCondition(sdk.Cond("destination-validator").EQ(sdk.EventValue(validator)))
	}
	return s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
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
}
