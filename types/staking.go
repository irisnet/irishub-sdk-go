package types

type StakingTx interface {
	Delegate(valAddr string, amount Coin, baseTx BaseTx) (Result, error)
	Undelegate(valAddr string, amount Coin, baseTx BaseTx) (Result, error)
	Redelegate(srcValAddr, dstValAddr string, amount Coin, baseTx BaseTx) (Result, error)
}

type StakingQuery interface {
	QueryDelegation(delAddr, valAddr string) (delegation Delegation, err error)
	QueryDelegations(delAddr string) (delegations Delegations, err error)

	QueryUnbondingDelegation(delAddr, valAddr string) (ubd UnbondingDelegation, err error)
	QueryUnbondingDelegations(delAddr string) (ubds UnbondingDelegations, err error)

	QueryRedelegation(delAddr, srcValAddr, dstValAddr string) (rd Redelegation, err error)
	QueryRedelegations(delAddr string) (rds Redelegations, err error)

	QueryDelegationsTo(valAddr string) (delegations Delegations, err error)
	QueryUnbondingDelegationsFrom(valAddr string) (ubds UnbondingDelegations, err error)
	QueryRedelegationsFrom(valAddr string) (rds Redelegations, err error)

	QueryValidator(address string) (validator Validator, err error)
	QueryValidators(page uint64, size uint16) (validators Validators, err error)

	QueryPool() (pool StakePool, err error)
	QueryParams() (params StakeParams, err error)
}

type Staking interface {
	StakingTx
	StakingQuery
}

type Delegation struct {
	DelegatorAddr string `json:"delegator_addr"`
	ValidatorAddr string `json:"validator_addr"`
	Shares        string `json:"shares"`
	Height        int64  `json:"height"`
}
type Delegations []Delegation

type UnbondingDelegations []UnbondingDelegation
type UnbondingDelegation struct {
	TxHash         string `json:"tx_hash"`
	DelegatorAddr  string `json:"delegator_addr"`
	ValidatorAddr  string `json:"validator_addr"`
	CreationHeight int64  `json:"creation_height"`
	MinTime        string `json:"min_time"`
	InitialBalance Coin   `json:"initial_balance"`
	Balance        Coin   `json:"balance"`
}

type Redelegations []Redelegation
type Redelegation struct {
	DelegatorAddr    string `json:"delegator_addr"`
	ValidatorSrcAddr string `json:"validator_src_addr"`
	ValidatorDstAddr string `json:"validator_dst_addr"`
	CreationHeight   int64  `json:"creation_height"`
	MinTime          string `json:"min_time"`
	InitialBalance   Coin   `json:"initial_balance"`
	Balance          Coin   `json:"balance"`
	SharesSrc        string `json:"shares_src"`
	SharesDst        string `json:"shares_dst"`
}

type TmValidator struct {
	Address          string `json:"address"`
	PubKey           string `json:"pub_key"`
	VotingPower      int64  `json:"voting_power"`
	ProposerPriority int64  `json:"proposer_priority"`
}

type Validators []Validator
type Validator struct {
	OperatorAddress string      `json:"operator_address"`
	ConsensusPubkey string      `json:"consensus_pubkey"`
	Jailed          bool        `json:"jailed"`
	Status          string      `json:"status"`
	Tokens          string      `json:"tokens"`
	DelegatorShares string      `json:"delegator_shares"`
	Description     Description `json:"description"`
	BondHeight      int64       `json:"bond_height"`
	UnbondingHeight int64       `json:"unbonding_height"`
	UnbondingTime   string      `json:"unbonding_time"`
	Commission      Commission  `json:"commission"`
}

// DelegatorShareExRate gets the exchange rate of tokens over delegator shares.
// UNITS: tokens/delegator-shares
func (v Validator) DelegatorShareExRate() Dec {
	delegatorShares, err := NewDecFromStr(v.DelegatorShares)
	if err != nil {
		return ZeroDec()
	}

	tokens, err := NewDecFromStr(v.Tokens)
	if err != nil {
		return ZeroDec()
	}
	if delegatorShares.IsZero() {
		return OneDec()
	}
	return tokens.Quo(delegatorShares)
}

type Commission struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
	UpdateTime    string `json:"update_time"`
}

type Description struct {
	Moniker  string `json:"moniker"`
	Identity string `json:"identity"`
	Website  string `json:"website"`
	Details  string `json:"details"`
}

type StakePool struct {
	LooseTokens  string `json:"loose_tokens"`
	BondedTokens string `json:"bonded_tokens"`
}

type StakeParams struct {
	UnbondingTime string `json:"unbonding_time"`
	MaxValidators int    `json:"max_validators"`
}
