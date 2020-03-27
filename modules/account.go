package modules

import (
	"fmt"
	"time"

	"github.com/irisnet/irishub-sdk-go/tools/cache"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

// Must be used with locker, otherwise there are thread safety issues
type localAccount struct {
	sdk.Queries
	*log.Logger
	cache.Cache

	keyManager sdk.KeyManager
	expiration time.Duration
}

func (l localAccount) Refresh(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := l.QueryAccount(address)
	if err != nil {
		l.Err(err).
			Str("address", address).
			Msg("update cache failed")
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	l.saveAccount(account)
	return account, nil
}

func (l localAccount) QueryAndRefreshAccount(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := l.Cache.Get(l.keyWithPrefix(address))
	if err != nil {
		return l.Refresh(address)
	}

	acc := account.(accountInfo)
	baseAcc := sdk.BaseAccount{
		Address:       sdk.MustAccAddressFromBech32(address),
		AccountNumber: acc.N,
		Sequence:      acc.S + 1,
	}
	l.saveAccount(baseAcc)

	l.Debug().
		Str("address", address).
		Msg("query account from cache")
	return baseAcc, nil
}

func (l localAccount) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	param := struct {
		Address sdk.AccAddress
	}{
		Address: addr,
	}

	var account sdk.BaseAccount
	if err := l.QueryWithResponse("custom/acc/account", param, &account); err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}
	l.Debug().
		Str("address", address).
		Msg("query account from chain")
	return account, nil
}

func (l localAccount) QueryAddress(name string) (sdk.AccAddress, sdk.Error) {
	addr, err := l.Cache.Get(l.keyWithPrefix(name))
	if err == nil {
		address, err := sdk.AccAddressFromBech32(addr.(string))
		if err != nil {
			l.Warn().
				Str("name", name).
				Msg("invalid address")
			_ = l.Remove(l.keyWithPrefix(name))
		} else {
			return address, nil
		}
	}

	address, err := l.keyManager.Query(name)
	if err != nil {
		l.Warn().
			Str("name", name).
			Msg("can't find account")
		return address, sdk.Wrap(err)
	}

	if err := l.SetWithExpire(l.keyWithPrefix(name), address.String(), l.expiration); err != nil {
		l.Warn().
			Str("name", name).
			Msg("cache user failed")
	}
	l.Debug().
		Str("name", name).
		Str("address", address.String()).
		Msg("query user from cache")
	return address, nil
}

func (l localAccount) saveAccount(account sdk.BaseAccount) {
	address := account.Address.String()
	info := accountInfo{
		N: account.AccountNumber,
		S: account.Sequence,
	}
	if err := l.SetWithExpire(l.keyWithPrefix(address), info, l.expiration); err != nil {
		l.Warn().
			Str("address", address).
			Msg("cache account failed")
		return
	}
	l.Debug().
		Str("address", address).
		Msgf("cache account %s", l.expiration.String())
}

func (l localAccount) keyWithPrefix(address string) string {
	return fmt.Sprintf("account:%s", address)
}

type accountInfo struct {
	N uint64 `json:"n"`
	S uint64 `json:"s"`
}
