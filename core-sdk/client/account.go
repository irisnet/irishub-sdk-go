package client

import (
	"context"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/auth"
	"github.com/irisnet/irishub-sdk-go/bank"
	cache "github.com/irisnet/irishub-sdk-go/common/cache"
	commoncodec "github.com/irisnet/irishub-sdk-go/common/codec"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

// Must be used with locker, otherwise there are thread safety issues
type accountQuery struct {
	sdk.Queries
	sdk.GRPCClient
	log.Logger
	cache.Cache
	cdc        commoncodec.Marshaler
	km         sdk.KeyManager
	expiration time.Duration
}

func (a accountQuery) QueryAndRefreshAccount(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := a.Get(a.prefixKey(address))
	if err != nil {
		return a.refresh(address)
	}

	acc := account.(accountInfo)
	baseAcc := sdk.BaseAccount{
		Address:       address,
		AccountNumber: acc.N,
		Sequence:      acc.S + 1,
	}
	a.saveAccount(baseAcc)

	a.Debug("query account from cache", "address", address)
	return baseAcc, nil
}

func (a accountQuery) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	conn, err := a.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	request := &auth.QueryAccountRequest{
		Address: address,
	}

	response, err := auth.NewQueryClient(conn).Account(context.Background(), request)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	var baseAccount auth.Account
	if err := a.cdc.UnpackAny(response.Account, &baseAccount); err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	account := baseAccount.(*auth.BaseAccount).ConvertAccount(a.cdc).(sdk.BaseAccount)

	breq := &bank.QueryAllBalancesRequest{
		Address:    address,
		Pagination: nil,
	}
	balances, err := bank.NewQueryClient(conn).AllBalances(context.Background(), breq)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	account.Coins = balances.Balances
	return account, nil
}

func (a accountQuery) QueryAddress(name, password string) (sdk.AccAddress, sdk.Error) {
	addr, err := a.Get(a.prefixKey(name))
	if err == nil {
		address, err := sdk.AccAddressFromBech32(addr.(string))
		if err != nil {
			a.Debug("invalid address", "name", name)
			_ = a.Remove(a.prefixKey(name))
		} else {
			return address, nil
		}
	}

	_, address, err := a.km.Find(name, password)
	if err != nil {
		a.Debug("can't find account", "name", name)
		return address, sdk.Wrap(err)
	}

	if err := a.SetWithExpire(a.prefixKey(name), address.String(), a.expiration); err != nil {
		a.Debug("cache user failed", "name", name)
	}
	a.Debug("query user from cache", "name", name, "address", address.String())
	return address, nil
}

func (a accountQuery) removeCache(address string) bool {
	return a.Remove(a.prefixKey(address))
}

func (a accountQuery) refresh(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := a.QueryAccount(address)
	if err != nil {
		a.Error("update cache failed", "address", address, "errMsg", err.Error())
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	a.saveAccount(account)
	return account, nil
}

func (a accountQuery) saveAccount(account sdk.BaseAccount) {
	address := account.Address
	info := accountInfo{
		N: account.AccountNumber,
		S: account.Sequence,
	}
	if err := a.SetWithExpire(a.prefixKey(address), info, a.expiration); err != nil {
		a.Debug("cache user failed", "address", account.Address)
		return
	}
	a.Debug("cache account", "address", address, "expiration", a.expiration.String())
}

func (a accountQuery) prefixKey(address string) string {
	return fmt.Sprintf("account:%s", address)
}

type accountInfo struct {
	N uint64 `json:"n"`
	S uint64 `json:"s"`
}
