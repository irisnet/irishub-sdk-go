package stake

import (
	"errors"
	"time"

	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	stakeStore    = "stake"
	validatorsKey = []byte{0x21} // prefix for each key to a validator
)

type Stake interface {
	QueryDelegation(delegatorAddr, validatorAddr string) (types.Delegation, error)
	QueryDelegations(delegatorAddr string) (types.Delegations, error)
	QueryUnbondingDelegation(delegatorAddr, validatorAddr string) (types.UnbondingDelegation, error)
	QueryUnbondingDelegations(delegatorAddr, validatorAddr string) (types.UnbondingDelegations, error)
	QueryRedelegation(delegatorAddr, srcValidatorAddr, dstValidatorAddr string) (types.Redelegation, error)
	QueryRedelegations(delegatorAddr string) (types.Redelegations, error)
	QueryDelegationsTo(validatorAddr string) (types.Delegations, error)
	QueryUnbondingDelegationsFrom(validatorAddr string) (types.UnbondingDelegations, error)
	QueryRedelegationsFrom(validatorAddr string) (types.Redelegations, error)
	QueryValidator(address string) (types.Validator, error)
	QueryValidators(page uint64, size uint16) (types.Validators, error)
	QueryAllValidators() (types.Validators, error)
	QueryPool() (types.StakePool, error)
	QueryParams() (types.StakeParams, error)
	Delegate(validatorAddr string, amount types.Coin, baseTx types.BaseTx) (types.Result, error)
	Undelegate(validatorAddr string, amount types.Coin, baseTx types.BaseTx) (types.Result, error)
	Redelegate(validatorSrcAddr, validatorDstAddr string, amount types.Coin, baseTx types.BaseTx) (types.Result, error)
}

type stakeClient struct {
	types.TxCtxManager
}

// unmarshal a redelegation from a store key and value
func mustUnmarshalValidator(cdc types.Codec, operatorAddr, value []byte) types.Validator {
	validator, err := unmarshalValidator(cdc, operatorAddr, value)
	if err != nil {
		panic(err)
	}
	return validator
}

// unmarshal a redelegation from a store key and value
func unmarshalValidator(cdc types.Codec, operatorAddr, value []byte) (validator types.Validator, err error) {
	if len(operatorAddr) != types.AddrLen {
		err = errors.New("bad address")
		return
	}
	var storeValue validatorValue
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &storeValue)
	if err != nil {
		return
	}

	conPubkey, err := types.Bech32ifyConsPub(storeValue.ConsPubKey)
	return types.Validator{
		OperatorAddress: types.ValAddress(operatorAddr).String(),
		ConsensusPubkey: conPubkey,
		Jailed:          storeValue.Jailed,
		Tokens:          storeValue.Tokens,
		Status:          storeValue.Status,
		DelegatorShares: storeValue.DelegatorShares,
		Description:     storeValue.Description,
		BondHeight:      storeValue.BondHeight,
		UnbondingHeight: storeValue.UnbondingHeight,
		UnbondingTime:   storeValue.UnbondingMinTime,
		Commission:      storeValue.Commission,
	}, nil
}

// what's kept in the store value
type validatorValue struct {
	ConsPubKey       crypto.PubKey
	Jailed           bool
	Status           types.BondStatus
	Tokens           types.Dec
	DelegatorShares  types.Dec
	Description      types.Description
	BondHeight       int64
	UnbondingHeight  int64
	UnbondingMinTime time.Time
	Commission       types.Commission
}

// defines the params for the following queries:
// - 'custom/stake/delegation'
// - 'custom/stake/unbondingDelegation'
// - 'custom/stake/delegatorValidator'
type QueryBondsParams struct {
	DelegatorAddr types.AccAddress
	ValidatorAddr types.ValAddress
}

// defines the params for the following queries:
// - 'custom/stake/delegatorDelegations'
// - 'custom/stake/delegatorUnbondingDelegations'
// - 'custom/stake/delegatorRedelegations'
// - 'custom/stake/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr types.AccAddress
}

// defines the params for the following queries:
// - 'custom/stake/validator'
// - 'custom/stake/validatorDelegations'
// - 'custom/stake/validatorUnbondingDelegations'
// - 'custom/stake/validatorRedelegations'
type QueryValidatorParams struct {
	ValidatorAddr types.ValAddress
}

type QueryRedelegationParams struct {
	DelegatorAddr types.AccAddress
	ValSrcAddr    types.ValAddress
	ValDstAddr    types.ValAddress
}
