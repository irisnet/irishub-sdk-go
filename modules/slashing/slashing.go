//
// In Proof-of-Stake blockchain, validators will get block provisions by staking their token.
// But if they failed to keep online, they will be punished by slashing a small portion of their staked tokens.
// The offline validators will be removed from the validator set and put into jail, which means their voting power is zero.
// During the jail period, these nodes are not even validator candidates. Once the jail period ends, they can send [[unjail]] transactions to free themselves and become validator candidates again.
//
// [More Details](https://www.irisnet.org/docs/features/slashing.html)
//
package slashing

import (
	"fmt"

	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/crypto"
)

type slashingClient struct {
	sdk.AbstractClient
	*log.Logger
}

func New(ac sdk.AbstractClient) sdk.Slashing {
	return slashingClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

//Unjail is responsible for unjail a validator previously jailed
func (s slashingClient) Unjail(baseTx sdk.BaseTx) (sdk.Result, error) {
	address, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	valAddr := sdk.ValAddress(address)
	msg := MsgUnjail{
		ValidatorAddr: valAddr,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//QueryParams return parameter for slashing at genesis
func (s slashingClient) QueryParams() (sdk.SlashingParams, error) {
	return s.queryParamsV017()
}

//QueryValidatorSigningInfo return the specified validator sign information
func (s slashingClient) QueryValidatorSigningInfo(validatorConPubKey string) (sdk.ValidatorSigningInfo, error) {
	pk, err := sdk.GetConsPubKeyBech32(validatorConPubKey)
	if err != nil {
		return sdk.ValidatorSigningInfo{}, err
	}
	return s.querySigningInfoV017(pk)
}

func (s slashingClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (s slashingClient) Name() string {
	return ModuleName
}

func (s slashingClient) queryParamsV017() (sdk.SlashingParams, error) {
	param := struct {
		Module string
	}{
		Module: s.Name(),
	}

	var params ParamsV017
	err := s.Query("custom/params/module", param, &params)
	if err != nil {
		return sdk.SlashingParams{}, err
	}
	return sdk.SlashingParams{
		MaxEvidenceAge:          fmt.Sprintf("%d", params.MaxEvidenceAge),
		SignedBlocksWindow:      params.SignedBlocksWindow,
		MinSignedPerWindow:      params.MinSignedPerWindow.String(),
		DoubleSignJailDuration:  params.DoubleSignJailDuration.String(),
		DowntimeJailDuration:    params.DowntimeJailDuration.String(),
		SlashFractionDoubleSign: params.SlashFractionDoubleSign.String(),
		SlashFractionDowntime:   params.SlashFractionDowntime.String(),
	}, nil
}

func (s slashingClient) queryParamsV100() (sdk.SlashingParams, error) {
	var params Params
	err := s.Query("custom/%s/parameters", s.Name(), &params)
	if err != nil {
		return sdk.SlashingParams{}, err
	}
	return params.ToSDKResponse(), nil
}

func (s slashingClient) querySigningInfoV017(pk crypto.PubKey) (sdk.ValidatorSigningInfo, error) {
	key := append([]byte{0x01}, pk.Bytes()...)
	res, err := s.QueryStore(key, s.Name())
	if err != nil {
		return sdk.ValidatorSigningInfo{}, err
	}

	var signingInfo ValidatorSigningInfoV017
	err = cdc.UnmarshalBinaryLengthPrefixed(res, &signingInfo)
	if err != nil {
		return sdk.ValidatorSigningInfo{}, err
	}

	consAddr := sdk.ConsAddress(pk.Address())
	return sdk.ValidatorSigningInfo{
		Address:             consAddr.String(),
		StartHeight:         signingInfo.StartHeight,
		IndexOffset:         signingInfo.IndexOffset,
		JailedUntil:         signingInfo.JailedUntil,
		Tombstoned:          false,
		MissedBlocksCounter: signingInfo.MissedBlocksCounter,
	}, nil
}

func (s slashingClient) querySigningInfoV100(pk crypto.PubKey) (sdk.ValidatorSigningInfo, error) {
	consAddr := sdk.ConsAddress(pk.Address())
	param := struct {
		ConsAddress sdk.ConsAddress
	}{
		ConsAddress: consAddr,
	}

	var signingInfo ValidatorSigningInfo
	err := s.Query("custom/slashing/signingInfo", param, &signingInfo)
	if err != nil {
		return sdk.ValidatorSigningInfo{}, err
	}
	return signingInfo.ToSDKResponse(), err
}
