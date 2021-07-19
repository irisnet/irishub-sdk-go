package staking

import (
	"time"

	sdk "github.com/irisnet/core-sdk-go/types"
)

// expose Staking module api for user
type Client interface {
	sdk.Module

	CreateValidator(request CreateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	EditValidator(request EditValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	Delegate(request DelegateRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	Undelegate(request UndelegateRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	BeginRedelegate(request BeginRedelegateRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryValidators(status string, page, size uint64) (QueryValidatorsResp, sdk.Error)
	QueryValidator(validatorAddr string) (QueryValidatorResp, sdk.Error)
	QueryValidatorDelegations(validatorAddr string, page, size uint64) (QueryValidatorDelegationsResp, sdk.Error)
	QueryValidatorUnbondingDelegations(validatorAddr string, page, size uint64) (QueryValidatorUnbondingDelegationsResp, sdk.Error)
	QueryDelegation(delegatorAddr string, validatorAddr string) (QueryDelegationResp, sdk.Error)
	QueryUnbondingDelegation(delegatorAddr string, validatorAddr string) (QueryUnbondingDelegationResp, sdk.Error)
	QueryDelegatorDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorDelegationsResp, sdk.Error)
	QueryDelegatorUnbondingDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorUnbondingDelegationsResp, sdk.Error)
	QueryRedelegations(request QueryRedelegationsReq) (QueryRedelegationsResp, sdk.Error)
	QueryDelegatorValidators(delegatorAddr string, page, size uint64) (QueryDelegatorValidatorsResp, sdk.Error)
	QueryDelegatorValidator(delegatorAddr string, validatorAddr string) (QueryValidatorResp, sdk.Error)
	QueryHistoricalInfo(height int64) (QueryHistoricalInfoResp, sdk.Error)
	QueryPool() (QueryPoolResp, sdk.Error)
	QueryParams() (QueryParamsResp, sdk.Error)
}

type CreateValidatorRequest struct {
	Moniker           string      `json:"moniker"`
	Rate              sdk.Dec     `json:"rate"`
	MaxRate           sdk.Dec     `json:"max_rate"`
	MaxChangeRate     sdk.Dec     `json:"max_change_rate"`
	MinSelfDelegation sdk.Int     `json:"min_self_delegation"`
	Pubkey            string      `json:"pubkey"`
	Value             sdk.DecCoin `json:"value"`
}

type EditValidatorRequest struct {
	Moniker           string  `json:"moniker"`
	Identity          string  `json:"identity"`
	Website           string  `json:"website"`
	SecurityContact   string  `json:"security_contact"`
	Details           string  `json:"details"`
	CommissionRate    sdk.Dec `json:"commission_rate"`
	MinSelfDelegation sdk.Int `json:"min_self_delegation"`
}

type DelegateRequest struct {
	ValidatorAddr string      `json:"validator_address"`
	Amount        sdk.DecCoin `json:"amount"`
}

type UndelegateRequest struct {
	ValidatorAddr string      `json:"validator_address"`
	Amount        sdk.DecCoin `json:"amount"`
}

type BeginRedelegateRequest struct {
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Amount              sdk.DecCoin
}

type (
	description struct {
		Moniker         string `json:"moniker"`
		Identity        string `json:"identity"`
		Website         string `json:"website"`
		SecurityContact string `json:"security_contact"`
		Details         string `json:"details"`
	}
	commission struct {
		commissionRates
		UpdateTime time.Time `json:"update_time"`
	}
	commissionRates struct {
		Rate          sdk.Dec `json:"rate"`
		MaxRate       sdk.Dec `json:"max_rate"`
		MaxChangeRate sdk.Dec `json:"max_change_rate"`
	}

	QueryValidatorsResp struct {
		Validators []QueryValidatorResp `json:"validators"`
		Total      uint64               `json:"total"`
	}

	QueryValidatorResp struct {
		OperatorAddress   string      `json:"operator_address"`
		ConsensusPubkey   string      `json:"consensus_pubkey"`
		Jailed            bool        `json:"jailed"`
		Status            string      `json:"status"`
		Tokens            sdk.Int     `json:"tokens"`
		DelegatorShares   sdk.Dec     `json:"delegator_shares"`
		Description       description `json:"description"`
		UnbondingHeight   int64       `json:"unbonding_height"`
		UnbondingTime     time.Time   `json:"unbonding_time"`
		Commission        commission  `json:"commission"`
		MinSelfDelegation sdk.Int     `json:"min_self_delegation"`
	}
)

type (
	delegation struct {
		DelegatorAddress string  `json:"delegator_address"`
		Shares           sdk.Dec `json:"shares"`
		ValidatorAddress string  `json:"validator_address"`
	}

	QueryDelegationResp struct {
		Delegation delegation `json:"delegation"`
		Balance    sdk.Coin   `json:"balance"`
	}

	QueryValidatorDelegationsResp struct {
		DelegationResponses []QueryDelegationResp `json:"delegation_responses"`
		Total               uint64                `json:"total"`
	}
)

type (
	unbondingDelegationEntry struct {
		CreationHeight int64     `json:"creation_height"`
		CompletionTime time.Time `json:"completion_time"`
		InitialBalance sdk.Int   `json:"initial_balance"`
		Balance        sdk.Int   `json:"balance"`
	}

	QueryUnbondingDelegationResp struct {
		DelegatorAddress string                     `json:"delegator_address"`
		ValidatorAddress string                     `json:"validator_address"`
		Entries          []unbondingDelegationEntry `json:"entries"`
	}

	QueryValidatorUnbondingDelegationsResp struct {
		UnbondingResponses []QueryUnbondingDelegationResp `json:"unbonding_responses"`
		Total              uint64                         `json:"total"`
	}
)

type QueryDelegatorDelegationsResp struct {
	DelegationResponses []QueryDelegationResp `json:"delegation_responses"`
	Total               uint64                `json:"total"`
}

type QueryDelegatorUnbondingDelegationsResp struct {
	UnbondingDelegations []QueryUnbondingDelegationResp `json:"unbonding_delegations"`
	Total                uint64                         `json:"total"`
}

type (
	QueryRedelegationsReq struct {
		DelegatorAddr    string `json:"delegator_addr"`
		SrcValidatorAddr string `json:"src_validator_addr"`
		DstValidatorAddr string `json:"dst_validator_addr"`
		Page             uint64 `json:"page"`
		Size             uint64 `json:"size"`
	}

	QueryRedelegationsResp struct {
		RedelegationResponses []RedelegationResp `json:"redelegation_responses"`
		Total                 uint64             `json:"total"`
	}

	redelegationEntry struct {
		CreationHeight int64     `json:"creation_height"`
		CompletionTime time.Time `json:"completion_time"`
		InitialBalance sdk.Int   `json:"initial_balance"`
		SharesDst      sdk.Dec   `json:"shares_dst"`
	}
	redelegationEntryResponse struct {
		RedelegationEntry redelegationEntry `json:"redelegation_entry"`
		Balance           sdk.Int           `json:"balance"`
	}
	redelegation struct {
		DelegatorAddress    string              `json:"delegator_address"`
		ValidatorSrcAddress string              `json:"validator_src_address"`
		ValidatorDstAddress string              `json:"validator_dst_address"`
		Entries             []redelegationEntry `json:"entries"`
	}

	RedelegationResp struct {
		Redelegation redelegation                `json:"redelegation"`
		Entries      []redelegationEntryResponse `json:"entries"`
	}
)

type QueryDelegatorValidatorsResp struct {
	Validator []QueryValidatorResp `json:"validator"`
	Total     uint64               `json:"total"`
}

type QueryHistoricalInfoResp struct {
	Header sdk.Header           `json:"header"`
	Valset []QueryValidatorResp `json:"valset"`
}

type QueryPoolResp struct {
	NotBondedTokens sdk.Int `json:"not_bonded_tokens"`
	BondedTokens    sdk.Int `json:"bonded_tokens"`
}

type QueryParamsResp struct {
	UnbondingTime     time.Duration `json:"unbonding_time"`
	MaxValidators     uint32        `json:"max_validators"`
	MaxEntries        uint32        `json:"max_entries"`
	HistoricalEntries uint32        `json:"historical_entries"`
	BondDenom         string        `json:"bond_denom"`
}
