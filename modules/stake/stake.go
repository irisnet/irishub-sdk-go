package stake

import (
	"github.com/irisnet/irishub-sdk-go/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func NewStakeClient(tm types.TxCtxManager) Stake {
	return stakeClient{
		TxCtxManager: tm,
	}
}

func (s stakeClient) QueryDelegation(delegatorAddr, validatorAddr string) (delegation types.Delegation, err error) {
	delAddr, err := types.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return delegation, err
	}
	varAddr, err := types.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return delegation, err
	}
	param := QueryBondsParams{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	err = s.Query("custom/stake/delegation", param, &delegation)
	return
}

func (s stakeClient) QueryDelegations(delegatorAddr string) (delegations types.Delegations, err error) {
	delAddr, err := types.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return delegations, err
	}
	param := QueryDelegatorParams{
		DelegatorAddr: delAddr,
	}

	err = s.Query("custom/stake/delegatorDelegations", param, &delegations)
	return
}

func (s stakeClient) QueryUnbondingDelegation(delegatorAddr, validatorAddr string) (types.UnbondingDelegation, error) {
	panic("implement me")
}

func (s stakeClient) QueryUnbondingDelegations(delegatorAddr, validatorAddr string) (types.UnbondingDelegations, error) {
	panic("implement me")
}

func (s stakeClient) QueryRedelegation(delegatorAddr, srcValidatorAddr, dstValidatorAddr string) (types.Redelegation, error) {
	panic("implement me")
}

func (s stakeClient) QueryRedelegations(delegatorAddr string) (types.Redelegation, error) {
	panic("implement me")
}

func (s stakeClient) QueryDelegationsTo(validatorAddr string) (types.Delegations, error) {
	panic("implement me")
}

func (s stakeClient) QueryUnbondingDelegationsFrom(validatorAddr string) (types.UnbondingDelegations, error) {
	panic("implement me")
}

func (s stakeClient) QueryRedelegationsFrom(validatorAddr string) (types.Redelegations, error) {
	panic("implement me")
}

func (s stakeClient) QueryValidator(valAddr string) (validator types.Validator, err error) {
	varAddr, err := types.ValAddressFromBech32(valAddr)
	if err != nil {
		return validator, err
	}
	param := QueryValidatorParams{
		ValidatorAddr: varAddr,
	}

	err = s.Query("custom/stake/validator", param, &validator)
	return
}

func (s stakeClient) QueryValidators(page, size int) (types.Validators, error) {
	panic("implement me")
}

func (s stakeClient) QueryAllValidators() (validators types.Validators, err error) {
	bz, err := s.QueryStore(validatorsKey, stakeStore)
	if err != nil {
		return validators, err
	}
	var resKVs []cmn.KVPair
	cdc := s.GetCodec()

	if err = cdc.UnmarshalBinaryLengthPrefixed(bz, &resKVs); err != nil {
		return validators, err
	}
	for _, kv := range resKVs {
		addr := kv.Key[1:]
		validator := MustUnmarshalValidator(cdc, addr, kv.Value)
		validators = append(validators, validator)
	}
	return
}

func (s stakeClient) QueryPool() (types.StakePool, error) {
	panic("implement me")
}

func (s stakeClient) QueryParams() (types.StakeParams, error) {
	panic("implement me")
}

func (s stakeClient) Delegate(validatorAddr string, amount types.Coin, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}

func (s stakeClient) Unbond(validatorAddr string, amount string, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}

func (s stakeClient) Redelegate(validatorSrcAddr, validatorDstAddr, amount string, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}
