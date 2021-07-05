package staking

import (
	"context"

	"github.com/irisnet/core-sdk-go/common"
	"github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/query"
)

type stakingClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(baseClient sdk.BaseClient, marshaler codec.Marshaler) Client {
	return &stakingClient{
		BaseClient: baseClient,
		Marshaler:  marshaler,
	}
}

func (sc stakingClient) Name() string {
	return ModuleName
}

func (sc stakingClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (sc stakingClient) CreateValidator(request CreateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(delegatorAddr.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	values, err := sc.ToMinCoin(request.Value)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	pk, e := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, request.Pubkey)
	if e != nil {
		return sdk.ResultTx{}, sdk.Wrap(e)
	}
	pkAny, e := types.PackAny(pk)
	if e != nil {
		return sdk.ResultTx{}, sdk.Wrap(e)
	}

	msg := &MsgCreateValidator{
		Description: Description{
			Moniker: request.Moniker,
		},
		Commission: CommissionRates{
			Rate:          request.Rate,
			MaxRate:       request.MaxRate,
			MaxChangeRate: request.MaxChangeRate,
		},
		MinSelfDelegation: request.MinSelfDelegation,
		DelegatorAddress:  delegatorAddr.String(),
		ValidatorAddress:  valAddr.String(),
		Pubkey:            pkAny,
		Value:             values[0],
	}
	return sc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (sc stakingClient) EditValidator(request EditValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(delegatorAddr.String())
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgEditValidator{
		Description: Description{
			Moniker:         request.Moniker,
			Identity:        request.Identity,
			Website:         request.Website,
			SecurityContact: request.SecurityContact,
			Details:         request.Details,
		},
		ValidatorAddress:  valAddr.String(),
		CommissionRate:    &request.CommissionRate,
		MinSelfDelegation: &request.MinSelfDelegation,
	}
	return sc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (sc stakingClient) Delegate(request DelegateRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	coins, err := sc.ToMinCoin(request.Amount)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgDelegate{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: request.ValidatorAddr,
		Amount:           coins[0],
	}
	return sc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (sc stakingClient) Undelegate(request UndelegateRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	coins, err := sc.ToMinCoin(request.Amount)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgUndelegate{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: request.ValidatorAddr,
		Amount:           coins[0],
	}
	return sc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (sc stakingClient) BeginRedelegate(request BeginRedelegateRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegatorAddr, err := sc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	coins, err := sc.ToMinCoin(request.Amount)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgBeginRedelegate{
		DelegatorAddress:    delegatorAddr.String(),
		ValidatorSrcAddress: request.ValidatorSrcAddress,
		ValidatorDstAddress: request.ValidatorDstAddress,
		Amount:              coins[0],
	}
	return sc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// QueryValidators when status is "" will return all status' validator
// about status, you can see BondStatus_value
func (sc stakingClient) QueryValidators(status string, page, size uint64) (QueryValidatorsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorsResp{}, sdk.Wrap(err)
	}

	offset, limit := common.ParsePage(page, size)
	res, err := NewQueryClient(conn).Validators(
		context.Background(),
		&QueryValidatorsRequest{
			Status: status,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryValidatorsResp{}, sdk.Wrap(err)
	}
	return res.Convert(sc.Marshaler).(QueryValidatorsResp), nil
}

func (sc stakingClient) QueryValidator(validatorAddr string) (QueryValidatorResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Validator(
		context.Background(),
		&QueryValidatorRequest{
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryValidatorResp{}, sdk.Wrap(err)
	}
	return res.Validator.Convert(sc.Marshaler).(QueryValidatorResp), nil
}

func (sc stakingClient) QueryValidatorDelegations(validatorAddr string, page, size uint64) (QueryValidatorDelegationsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorDelegationsResp{}, sdk.Wrap(err)
	}

	offset, limit := common.ParsePage(page, size)
	res, err := NewQueryClient(conn).ValidatorDelegations(
		context.Background(),
		&QueryValidatorDelegationsRequest{
			ValidatorAddr: validatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryValidatorDelegationsResp{}, sdk.Wrap(err)
	}
	return res.Convert().(QueryValidatorDelegationsResp), nil
}

func (sc stakingClient) QueryValidatorUnbondingDelegations(validatorAddr string, page, size uint64) (QueryValidatorUnbondingDelegationsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorUnbondingDelegationsResp{}, sdk.Wrap(err)
	}

	offset, limit := common.ParsePage(page, size)
	res, err := NewQueryClient(conn).ValidatorUnbondingDelegations(
		context.Background(),
		&QueryValidatorUnbondingDelegationsRequest{
			ValidatorAddr: validatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryValidatorUnbondingDelegationsResp{}, sdk.Wrap(err)
	}
	return res.Convert().(QueryValidatorUnbondingDelegationsResp), nil
}

func (sc stakingClient) QueryDelegation(delegatorAddr string, validatorAddr string) (QueryDelegationResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegationResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Delegation(
		context.Background(),
		&QueryDelegationRequest{
			DelegatorAddr: delegatorAddr,
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryDelegationResp{}, sdk.Wrap(err)
	}
	return res.DelegationResponse.Convert().(QueryDelegationResp), nil
}

func (sc stakingClient) QueryUnbondingDelegation(delegatorAddr string, validatorAddr string) (QueryUnbondingDelegationResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryUnbondingDelegationResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).UnbondingDelegation(
		context.Background(),
		&QueryUnbondingDelegationRequest{
			DelegatorAddr: delegatorAddr,
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryUnbondingDelegationResp{}, sdk.Wrap(err)
	}
	return res.Unbond.Convert().(QueryUnbondingDelegationResp), nil
}

func (sc stakingClient) QueryDelegatorDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorDelegationsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegatorDelegationsResp{}, sdk.Wrap(err)
	}

	offset, limit := common.ParsePage(page, size)
	res, err := NewQueryClient(conn).DelegatorDelegations(
		context.Background(),
		&QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryDelegatorDelegationsResp{}, sdk.Wrap(err)
	}
	return res.Convert().(QueryDelegatorDelegationsResp), nil
}

func (sc stakingClient) QueryDelegatorUnbondingDelegations(delegatorAddr string, page, size uint64) (QueryDelegatorUnbondingDelegationsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegatorUnbondingDelegationsResp{}, sdk.Wrap(err)
	}

	offset, limit := common.ParsePage(page, size)
	res, err := NewQueryClient(conn).DelegatorUnbondingDelegations(
		context.Background(),
		&QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: delegatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryDelegatorUnbondingDelegationsResp{}, sdk.Wrap(err)
	}
	return res.Convert().(QueryDelegatorUnbondingDelegationsResp), nil
}

func (sc stakingClient) QueryRedelegations(request QueryRedelegationsReq) (QueryRedelegationsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryRedelegationsResp{}, sdk.Wrap(err)
	}

	offset, limit := common.ParsePage(request.Page, request.Size)
	res, err := NewQueryClient(conn).Redelegations(
		context.Background(),
		&QueryRedelegationsRequest{
			DelegatorAddr:    request.DelegatorAddr,
			SrcValidatorAddr: request.SrcValidatorAddr,
			DstValidatorAddr: request.DstValidatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryRedelegationsResp{}, sdk.Wrap(err)
	}
	return res.Convert().(QueryRedelegationsResp), nil
}

func (sc stakingClient) QueryDelegatorValidators(delegatorAddr string, page, size uint64) (QueryDelegatorValidatorsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryDelegatorValidatorsResp{}, sdk.Wrap(err)
	}

	offset, limit := common.ParsePage(page, size)
	res, err := NewQueryClient(conn).DelegatorValidators(
		context.Background(),
		&QueryDelegatorValidatorsRequest{
			DelegatorAddr: delegatorAddr,
			Pagination: &query.PageRequest{
				Offset:     offset,
				Limit:      limit,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return QueryDelegatorValidatorsResp{}, sdk.Wrap(err)
	}
	return res.Convert(sc.Marshaler).(QueryDelegatorValidatorsResp), nil
}

func (sc stakingClient) QueryDelegatorValidator(delegatorAddr string, validatorAddr string) (QueryValidatorResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryValidatorResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).DelegatorValidator(
		context.Background(),
		&QueryDelegatorValidatorRequest{
			DelegatorAddr: delegatorAddr,
			ValidatorAddr: validatorAddr,
		},
	)
	if err != nil {
		return QueryValidatorResp{}, sdk.Wrap(err)
	}
	return res.Validator.Convert(sc.Marshaler).(QueryValidatorResp), nil
}

// QueryHistoricalInfo tendermint only save latest 100 block, previous block is aborted
func (sc stakingClient) QueryHistoricalInfo(height int64) (QueryHistoricalInfoResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryHistoricalInfoResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).HistoricalInfo(
		context.Background(),
		&QueryHistoricalInfoRequest{
			Height: height,
		},
	)
	if err != nil {
		return QueryHistoricalInfoResp{}, sdk.Wrap(err)
	}
	return res.Convert(sc.Marshaler).(QueryHistoricalInfoResp), nil
}

func (sc stakingClient) QueryPool() (QueryPoolResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryPoolResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Pool(
		context.Background(),
		&QueryPoolRequest{},
	)
	if err != nil {
		return QueryPoolResp{}, sdk.Wrap(err)
	}

	return QueryPoolResp{
		NotBondedTokens: res.Pool.NotBondedTokens,
		BondedTokens:    res.Pool.BondedTokens,
	}, nil
}

func (sc stakingClient) QueryParams() (QueryParamsResp, sdk.Error) {
	conn, err := sc.GenConn()

	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{},
	)
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}
	return res.Convert().(QueryParamsResp), nil
}
