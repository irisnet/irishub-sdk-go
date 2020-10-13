package original

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"time"

	"github.com/irisnet/irishub-sdk-go/utils/cache"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

// Must be used with locker, otherwise there are thread safety issues
type accountQuery struct {
	original.Queries
	*log.Logger
	cache.Cache

	keyManager original.KeyManager
	expiration time.Duration
}

func (a accountQuery) QueryAndRefreshAccount(address string) (original.BaseAccount, original.Error) {
	account, err := a.Get(a.prefixKey(address))
	if err != nil {
		return a.refresh(address)
	}

	acc := account.(accountInfo)
	baseAcc := original.BaseAccount{
		Address:       original.MustAccAddressFromBech32(address),
		AccountNumber: acc.N,
		Sequence:      acc.S + 1,
	}
	a.saveAccount(baseAcc)

	a.Debug().
		Str("address", address).
		Msg("query account from cache")
	return baseAcc, nil
}

func (a accountQuery) QueryAccount(address string) (original.BaseAccount, original.Error) {
	addr, err := original.AccAddressFromBech32(address)
	if err != nil {
		return original.BaseAccount{}, original.Wrap(err)
	}

	param := struct {
		Address original.AccAddress `json:"address"`
	}{
		Address: addr,
	}

	var account original.BaseAccount
	if err := a.QueryWithResponse("custom/auth/account", param, &account); err != nil {
		return original.BaseAccount{}, original.Wrap(err)
	}
	a.Debug().
		Str("address", address).
		Msg("query account from chain")
	return account, nil
}

func (a accountQuery) QueryAddress(name string) (original.AccAddress, original.Error) {
	addr, err := a.Get(a.prefixKey(name))
	if err == nil {
		address, err := original.AccAddressFromBech32(addr.(string))
		if err != nil {
			a.Warn().
				Str("name", name).
				Msg("invalid address")
			_ = a.Remove(a.prefixKey(name))
		} else {
			return address, nil
		}
	}

	address, err := a.keyManager.Query(name)
	if err != nil {
		a.Warn().
			Str("name", name).
			Msg("can't find account")
		return address, original.Wrap(err)
	}

	if err := a.SetWithExpire(a.prefixKey(name), address.String(), a.expiration); err != nil {
		a.Warn().
			Str("name", name).
			Msg("cache user failed")
	}
	a.Debug().
		Str("name", name).
		Str("address", address.String()).
		Msg("query user from cache")
	return address, nil
}

func (a accountQuery) removeCache(address string) bool {
	return a.Remove(a.prefixKey(address))
}

func (a accountQuery) refresh(address string) (original.BaseAccount, original.Error) {
	account, err := a.QueryAccount(address)
	if err != nil {
		a.Err(err).
			Str("address", address).
			Msg("update cache failed")
		return original.BaseAccount{}, original.Wrap(err)
	}

	a.saveAccount(account)
	return account, nil
}

func (a accountQuery) saveAccount(account original.BaseAccount) {
	address := account.Address.String()
	info := accountInfo{
		N: account.AccountNumber,
		S: account.Sequence,
	}
	if err := a.SetWithExpire(a.prefixKey(address), info, a.expiration); err != nil {
		a.Warn().
			Str("address", address).
			Msg("cache account failed")
		return
	}
	a.Debug().
		Str("address", address).
		Msgf("cache account %s", a.expiration.String())
}

func (a accountQuery) prefixKey(address string) string {
	return fmt.Sprintf("account:%s", address)
}

type accountInfo struct {
	N uint64 `json:"n"`
	S uint64 `json:"s"`
}
