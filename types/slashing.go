package types

import "time"

type Slashing interface {
	Module
	Unjail(baseTx BaseTx) (Result, error)
	QueryParams() (SlashingParams, error)
	QueryValidatorSigningInfo(validatorConPubKey string) (ValidatorSigningInfo, error)
}

type SlashingParams struct {
	MaxEvidenceAge          string `json:"max_evidence_age"`
	SignedBlocksWindow      int64  `json:"signed_blocks_window"`
	MinSignedPerWindow      string `json:"min_signed_per_window"`
	DoubleSignJailDuration  string `json:"double_sign_jail_duration"`
	DowntimeJailDuration    string `json:"downtime_jail_duration"`
	CensorshipJailDuration  string `json:"censorship_jail_duration"` // delete by v1.0.0
	SlashFractionDoubleSign string `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `json:"slash_fraction_downtime"`
	SlashFractionCensorship string `json:"slash_fraction_censorship"` // delete by v1.0.0
}

// ValidatorSigningInfo defines the signing info for a validator
type ValidatorSigningInfo struct {
	Address             string    `json:"address"`               // validator consensus address
	StartHeight         int64     `json:"start_height"`          // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64     `json:"index_offset"`          // index offset into signed block bit array
	JailedUntil         time.Time `json:"jailed_until"`          // timestamp validator cannot be unjailed until
	Tombstoned          bool      `json:"tombstoned"`            // whether or not a validator has been tombstoned (killed out of validator set)
	MissedBlocksCounter int64     `json:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
}
