package distr

import "github.com/irisnet/irishub-sdk-go/types"

type Client interface {
	QueryWithdrawAddress(delAddr string) (types.AccAddress, error)
	QueryDelegationDistInfo(delAddr, valAddr string) (types.DelegationDistInfo, error)
	QueryAllDelegationDistInfo(delAddr string) (types.DelegationDistInfo, error)
	QueryValidatorDistInfo(valAddr string) (types.DelegationDistInfo, error)
	QueryRewards(address string) (types.Result, error)
	SetWithdrawAddress(delAddr, valAddr string, baseTx types.BaseTx) (types.Result, error)
	WithdrawDelegationRewardsAll(delAddr string, baseTx types.BaseTx) (types.Result, error)
	WithdrawDelegationReward(delAddr, valAddr string, baseTx types.BaseTx) (types.Result, error)
	WithdrawValidatorRewardsAll(valAddr string, baseTx types.BaseTx) (types.Result, error)
}

type distrClient struct {
	types.AbstractClient
}

func (d distrClient) QueryWithdrawAddress(delAddr string) (types.AccAddress, error) {
	panic("implement me")
}

func (d distrClient) QueryDelegationDistInfo(delAddr, valAddr string) (types.DelegationDistInfo, error) {
	panic("implement me")
}

func (d distrClient) QueryAllDelegationDistInfo(delAddr string) (types.DelegationDistInfo, error) {
	panic("implement me")
}

func (d distrClient) QueryValidatorDistInfo(valAddr string) (types.DelegationDistInfo, error) {
	panic("implement me")
}

func (d distrClient) QueryRewards(address string) (types.Result, error) {
	panic("implement me")
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
