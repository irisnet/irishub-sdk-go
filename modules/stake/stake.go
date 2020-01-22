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

func (s stakeClient) QueryUnbondingDelegation(delegatorAddr, validatorAddr string) (ubd types.UnbondingDelegation, err error) {
	delAddr, err := types.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return ubd, err
	}
	varAddr, err := types.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return ubd, err
	}
	param := QueryBondsParams{
		DelegatorAddr: delAddr,
		ValidatorAddr: varAddr,
	}

	err = s.Query("custom/stake/unbondingDelegation", param, &ubd)
	return
}

func (s stakeClient) QueryUnbondingDelegations(delegatorAddr, validatorAddr string) (ubds types.UnbondingDelegations, err error) {
	delAddr, err := types.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return ubds, err
	}
	param := QueryDelegatorParams{
		DelegatorAddr: delAddr,
	}

	err = s.Query("custom/stake/delegatorUnbondingDelegations", param, &ubds)
	return
}

func (s stakeClient) QueryRedelegation(delegatorAddr, srcValidatorAddr, dstValidatorAddr string) (rd types.Redelegation, err error) {
	delAddr, err := types.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rd, err
	}
	srcVarAddr, err := types.ValAddressFromBech32(srcValidatorAddr)
	if err != nil {
		return rd, err
	}

	dstVarAddr, err := types.ValAddressFromBech32(dstValidatorAddr)
	if err != nil {
		return rd, err
	}

	param := QueryRedelegationParams{
		DelegatorAddr: delAddr,
		ValSrcAddr:    srcVarAddr,
		ValDstAddr:    dstVarAddr,
	}

	err = s.Query("custom/stake/redelegation", param, &rd)
	return
}

func (s stakeClient) QueryRedelegations(delegatorAddr string) (rds types.Redelegations, err error) {
	delAddr, err := types.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return rds, err
	}
	param := QueryDelegatorParams{
		DelegatorAddr: delAddr,
	}

	err = s.Query("custom/stake/delegatorRedelegations", param, &rds)
	return
}

func (s stakeClient) QueryDelegationsTo(validatorAddr string) (delegations types.Delegations, err error) {
	varAddr, err := types.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return delegations, err
	}
	param := QueryValidatorParams{
		ValidatorAddr: varAddr,
	}

	err = s.Query("custom/stake/validatorDelegations", param, &delegations)
	return
}

func (s stakeClient) QueryUnbondingDelegationsFrom(validatorAddr string) (ubds types.UnbondingDelegations, err error) {
	varAddr, err := types.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return ubds, err
	}
	param := QueryValidatorParams{
		ValidatorAddr: varAddr,
	}

	err = s.Query("custom/stake/validatorUnbondingDelegations", param, &ubds)
	return
}

func (s stakeClient) QueryRedelegationsFrom(validatorAddr string) (rds types.Redelegations, err error) {
	varAddr, err := types.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return rds, err
	}
	param := QueryValidatorParams{
		ValidatorAddr: varAddr,
	}

	err = s.Query("custom/stake/validatorRedelegations", param, &rds)
	return
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

func (s stakeClient) QueryValidators(page uint64, size uint16) (validators types.Validators, err error) {
	param := types.PaginationParams{
		Page: page,
		Size: size,
	}
	err = s.Query("custom/stake/validators", param, &validators)
	return
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
		validator := mustUnmarshalValidator(cdc, addr, kv.Value)
		validators = append(validators, validator)
	}
	return
}

func (s stakeClient) QueryPool() (pool types.StakePool, err error) {
	err = s.Query("custom/stake/pool", nil, &pool)
	return
}

func (s stakeClient) QueryParams() (params types.StakeParams, err error) {
	err = s.Query("custom/stake/parameters", nil, &params)
	return
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
