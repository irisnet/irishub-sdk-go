package slashing

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type slashingClient struct {
	original.BaseClient
	*log.Logger
}

func Create(ac original.BaseClient) rpc.Slashing {
	return slashingClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

//QueryParams return parameter for slashing at genesis
func (s slashingClient) QueryParams() (rpc.SlashingParams, original.Error) {
	return s.queryParamsV017()
}

//QueryValidatorSigningInfo return the specified validator sign information
func (s slashingClient) QueryValidatorSigningInfo(validatorConPubKey string) (rpc.ValidatorSigningInfo, original.Error) {
	pk, err := original.GetConsPubKeyBech32(validatorConPubKey)
	if err != nil {
		return rpc.ValidatorSigningInfo{}, original.Wrap(err)
	}
	return s.querySigningInfoV100(pk)
}

func (s slashingClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (s slashingClient) Name() string {
	return ModuleName
}

func (s slashingClient) queryParamsV017() (rpc.SlashingParams, original.Error) {
	param := struct {
		Module string
	}{
		Module: s.Name(),
	}

	var params paramsV017
	if err := s.QueryWithResponse("custom/params/module", param, &params); err != nil {
		return rpc.SlashingParams{}, original.Wrap(err)
	}
	return params.Convert().(rpc.SlashingParams), nil
}

func (s slashingClient) queryParamsV100() (rpc.SlashingParams, error) {
	var params params
	err := s.QueryWithResponse("custom/%s/parameters", s.Name(), &params)
	if err != nil {
		return rpc.SlashingParams{}, err
	}
	return params.Convert().(rpc.SlashingParams), nil
}

func (s slashingClient) querySigningInfoV017(pk crypto.PubKey) (rpc.ValidatorSigningInfo, original.Error) {
	key := append([]byte{0x01}, pk.Address().Bytes()...)
	res, err := s.QueryStore(key, s.Name())
	if err != nil {
		return rpc.ValidatorSigningInfo{}, original.Wrap(err)
	}

	var signingInfo validatorSigningInfoV017
	err = cdc.UnmarshalBinaryBare(res, &signingInfo)
	if err != nil {
		return rpc.ValidatorSigningInfo{}, original.Wrap(err)
	}

	consAddr := original.ConsAddress(pk.Address())
	return rpc.ValidatorSigningInfo{
		Address:             consAddr.String(),
		StartHeight:         signingInfo.StartHeight,
		IndexOffset:         signingInfo.IndexOffset,
		JailedUntil:         signingInfo.JailedUntil,
		Tombstoned:          false,
		MissedBlocksCounter: signingInfo.MissedBlocksCounter,
	}, nil
}

func (s slashingClient) querySigningInfoV100(pk crypto.PubKey) (rpc.ValidatorSigningInfo, original.Error) {
	consAddr := original.ConsAddress(pk.Address())
	param := struct {
		ConsAddress original.ConsAddress `json:"cons_address"`
	}{
		ConsAddress: consAddr,
	}

	var signingInfo validatorSigningInfo
	err := s.QueryWithResponse("custom/slashing/signingInfo", param, &signingInfo)
	if err != nil {
		return rpc.ValidatorSigningInfo{}, original.Wrap(err)
	}
	return signingInfo.Convert().(rpc.ValidatorSigningInfo), original.Wrap(err)
}
