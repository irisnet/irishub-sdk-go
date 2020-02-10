package distr

import "github.com/irisnet/irishub-sdk-go/types"

func New(ac types.AbstractClient) Distr {
	return distrClient{
		AbstractClient: ac,
	}
}

func (d distrClient) QueryWithdrawAddress(delAddr string) (res types.AccAddress, err error) {
	delegator, err := types.AccAddressFromBech32(delAddr)
	if err != nil {
		return res, err
	}
	param := QueryDelegatorParams{
		DelegatorAddress: delegator,
	}
	err = d.Query("custom/distr/withdraw_addr", param, &res)
	return res, err
}

func (d distrClient) QueryDelegationDistInfo(delAddr, valAddr string) (res types.DelegationDistInfo, err error) {
	panic("implement me")
}

func (d distrClient) QueryAllDelegationDistInfo(delAddr string) (types.DelegationDistInfo, error) {
	panic("implement me")
}

func (d distrClient) QueryValidatorDistInfo(valAddr string) (types.DelegationDistInfo, error) {
	panic("implement me")
}

func (d distrClient) QueryRewards(address string) (res types.Rewards, err error) {
	delegator, err := types.AccAddressFromBech32(address)
	if err != nil {
		return res, err
	}
	param := QueryDelegatorParams{
		DelegatorAddress: delegator,
	}
	err = d.Query("custom/distr/rewards", param, &res)
	return res, err
}

func (d distrClient) SetWithdrawAddress(delAddr, valAddr string, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}

func (d distrClient) WithdrawDelegationRewardsAll(delAddr string, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}

func (d distrClient) WithdrawDelegationReward(delAddr, valAddr string, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}

func (d distrClient) WithdrawValidatorRewardsAll(valAddr string, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}
