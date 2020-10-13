package distribution

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/irisnet/irishub-sdk-go/utils/bech32"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type distributionClient struct {
	original.BaseClient
	*log.Logger
}

func (d distributionClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (d distributionClient) Name() string {
	return ModuleName
}

func Create(ac original.BaseClient) rpc.Distribution {
	return distributionClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (d distributionClient) QueryRewards(delAddrOrValAddr string) (rpc.Rewards, original.Error) {
	_, bz, err := bech32.DecodeAndConvert(delAddrOrValAddr)
	if err != nil {
		return rpc.Rewards{}, original.Wrap(err)
	}

	param := struct {
		DelegatorAddress original.AccAddress `json:"delegator_address"`
	}{
		DelegatorAddress: original.AccAddress(bz),
	}

	var rewards rewardsResponse
	if err := d.QueryWithResponse("custom/distribution/delegator_total_rewards", param, &rewards); err != nil {
		return rpc.Rewards{}, original.Wrap(err)
	}
	return rewards.Convert().(rpc.Rewards), nil
}

func (d distributionClient) QueryWithdrawAddr(delAddrOrValAddr string) (string, original.Error) {
	_, address, err := bech32.DecodeAndConvert(delAddrOrValAddr)
	if err != nil {
		return "", original.Wrap(err)
	}

	param := struct {
		DelegatorAddress original.AccAddress `json:"delegator_address"`
	}{
		DelegatorAddress: original.AccAddress(address),
	}

	res, newErr := d.Query("custom/distribution/withdraw_addr", param)
	if newErr != nil {
		return "", original.Wrap(err)
	}
	return string(res), nil
}

func (d distributionClient) QueryCommission(validator string) (rpc.ValidatorAccumulatedCommission, original.Error) {
	address, err := original.ValAddressFromBech32(validator)
	if err != nil {
		return rpc.ValidatorAccumulatedCommission{}, original.Wrap(err)
	}

	param := struct {
		ValidatorAddress original.ValAddress `json:"validator_address"`
	}{
		ValidatorAddress: address,
	}

	var commission validatorAccumulatedCommission
	if err := d.QueryWithResponse("custom/distribution/validator_commission", param, &commission); err != nil {
		return rpc.ValidatorAccumulatedCommission{}, original.Wrap(err)
	}

	return commission.Convert().(rpc.ValidatorAccumulatedCommission), nil
}

func (d distributionClient) SetWithdrawAddr(withdrawAddr string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	delegator, err := d.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	withdraw, err := original.AccAddressFromBech32(withdrawAddr)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgSetWithdrawAddress{
		DelegatorAddr: delegator,
		WithdrawAddr:  withdraw,
	}
	d.Info().Str("delegator", delegator.String()).
		Str("withdrawAddr", withdrawAddr).
		Msg("execute setWithdrawAddr transaction")
	return d.BuildAndSend([]original.Msg{msg}, baseTx)
}

func (d distributionClient) WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	delegator, err := d.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	var msgs []original.Msg
	switch {
	case isValidator:
		msgs = append(msgs, MsgWithdrawValidatorCommission{
			ValidatorAddr: original.ValAddress(delegator.Bytes()),
		})

		d.Info().Str("delegator", delegator.String()).
			Msg("execute withdrawValidatorRewardsAll transaction")
		break
	case onlyFromValidator != "":
		valAddr, err := original.ValAddressFromBech32(onlyFromValidator)
		if err != nil {
			return original.ResultTx{}, original.Wrap(err)
		}
		msgs = append(msgs, MsgWithdrawDelegatorReward{
			ValidatorAddr: valAddr,
			DelegatorAddr: delegator,
		})

		d.Info().Str("delegator", delegator.String()).
			Str("validator", onlyFromValidator).
			Msg("execute withdrawDelegatorReward transaction")
		break
	default:
		msgs = append(msgs, MsgWithdrawDelegatorRewardsAll{
			DelegatorAddr: delegator,
		})

		d.Info().Str("delegator", delegator.String()).
			Str("validator", onlyFromValidator).
			Msg("execute withdrawDelegatorRewardsAll transaction")
		break
	}
	return d.BuildAndSend(msgs, baseTx)
}
