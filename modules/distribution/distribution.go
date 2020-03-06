package distribution

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type distributionClient struct {
	sdk.AbstractClient
	*log.Logger
}

func (d distributionClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (d distributionClient) Name() string {
	return ModuleName
}

func New(ac sdk.AbstractClient) rpc.Distribution {
	return distributionClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

func (d distributionClient) QueryRewards(delegator string) (rpc.Rewards, error) {
	address, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return rpc.Rewards{}, err
	}

	param := struct {
		Address sdk.AccAddress
	}{
		Address: address,
	}

	var rewards rewards
	if err = d.QueryWithResponse("custom/distr/rewards", param, &rewards); err != nil {
		return rpc.Rewards{}, err
	}
	return rewards.Convert().(rpc.Rewards), nil
}

func (d distributionClient) SetWithdrawAddr(withdrawAddr string, baseTx sdk.BaseTx) (sdk.Result, error) {
	delegator, err := d.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	withdraw, err := sdk.AccAddressFromBech32(withdrawAddr)
	if err != nil {
		return nil, err
	}

	msg := MsgSetWithdrawAddress{
		DelegatorAddr: delegator,
		WithdrawAddr:  withdraw,
	}
	d.Info().Str("delegator", delegator.String()).
		Str("withdrawAddr", withdrawAddr).
		Msg("execute setWithdrawAddr transaction")
	return d.Broadcast(baseTx, []sdk.Msg{msg})
}

func (d distributionClient) WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx sdk.BaseTx) (sdk.Result, error) {
	delegator, err := d.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	var msgs []sdk.Msg
	switch {
	case isValidator:
		msgs = append(msgs, MsgWithdrawValidatorRewardsAll{
			ValidatorAddr: sdk.ValAddress(delegator.Bytes()),
		})

		d.Info().Str("delegator", delegator.String()).
			Msg("execute withdrawValidatorRewardsAll transaction")
		break
	case onlyFromValidator != "":
		valAddr, err := sdk.ValAddressFromBech32(onlyFromValidator)
		if err != nil {
			return nil, err
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
	return d.Broadcast(baseTx, msgs)
}
