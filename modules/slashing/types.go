package slashing

import (
	"errors"
	"fmt"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	ModuleName = "slashing"
)

var (
	_ sdk.Msg = MsgUnjail{}

	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

type MsgUnjail struct {
	ValidatorAddr sdk.ValAddress `json:"address"` // address of the validator operator
}

//nolint
func (msg MsgUnjail) Route() string { return ModuleName }
func (msg MsgUnjail) Type() string  { return "unjail" }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgUnjail) ValidateBasic() error {
	if msg.ValidatorAddr == nil {
		return errors.New("validator is missed")
	}
	return nil
}

type paramsV017 struct {
	MaxEvidenceAge          int64         `json:"max_evidence_age"`
	SignedBlocksWindow      int64         `json:"signed_blocks_window"`
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window"`
	DoubleSignJailDuration  time.Duration `json:"double_sign_jail_duration"`
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration"`
	CensorshipJailDuration  time.Duration `json:"censorship_jail_duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime"`
	SlashFractionCensorship sdk.Dec       `json:"slash_fraction_censorship"`
}

func (p paramsV017) Convert() interface{} {
	return rpc.SlashingParams{
		MaxEvidenceAge:          fmt.Sprintf("%d", p.MaxEvidenceAge),
		SignedBlocksWindow:      p.SignedBlocksWindow,
		MinSignedPerWindow:      p.MinSignedPerWindow.String(),
		DoubleSignJailDuration:  p.DoubleSignJailDuration.String(),
		DowntimeJailDuration:    p.DowntimeJailDuration.String(),
		SlashFractionDoubleSign: p.SlashFractionDoubleSign.String(),
		SlashFractionDowntime:   p.SlashFractionDowntime.String(),
	}
}

type params struct {
	MaxEvidenceAge          time.Duration `json:"max_evidence_age"`
	SignedBlocksWindow      int64         `json:"signed_blocks_window"`
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window"`
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime"`
}

func (params params) Convert() interface{} {
	return rpc.SlashingParams{
		MaxEvidenceAge:          params.MaxEvidenceAge.String(),
		SignedBlocksWindow:      params.SignedBlocksWindow,
		MinSignedPerWindow:      params.MinSignedPerWindow.String(),
		DowntimeJailDuration:    params.DowntimeJailDuration.String(),
		SlashFractionDoubleSign: params.SlashFractionDoubleSign.String(),
		SlashFractionDowntime:   params.SlashFractionDowntime.String(),
	}
}

// Signing info for a validator
type validatorSigningInfoV017 struct {
	StartHeight         int64     `json:"start_height"`          // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64     `json:"index_offset"`          // index offset into signed block bit array
	JailedUntil         time.Time `json:"jailed_until"`          // timestamp validator cannot be unjailed until
	MissedBlocksCounter int64     `json:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
}

// validatorSigningInfo defines the signing info for a validator
type validatorSigningInfo struct {
	Address             sdk.ConsAddress `json:"address"`               // validator consensus address
	StartHeight         int64           `json:"start_height"`          // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64           `json:"index_offset"`          // index offset into signed block bit array
	JailedUntil         time.Time       `json:"jailed_until"`          // timestamp validator cannot be unjailed until
	Tombstoned          bool            `json:"tombstoned"`            // whether or not a validator has been tombstoned (killed out of validator set)
	MissedBlocksCounter int64           `json:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
}

func (vsi validatorSigningInfo) Convert() interface{} {
	return rpc.ValidatorSigningInfo{
		Address:             vsi.Address.String(),
		StartHeight:         vsi.StartHeight,
		IndexOffset:         vsi.IndexOffset,
		JailedUntil:         vsi.JailedUntil,
		Tombstoned:          vsi.Tombstoned,
		MissedBlocksCounter: vsi.MissedBlocksCounter,
	}
}

func registerCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "irishub/slashing/MsgUnjail")
	cdc.RegisterConcrete(&paramsV017{}, "irishub/slashing/Params")
}
