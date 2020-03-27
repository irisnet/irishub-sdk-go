package modules

import (
	"time"

	"github.com/irisnet/irishub-sdk-go/tools/cache"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

// Must be used with locker, otherwise there are thread safety issues
type localAccount struct {
	sdk.Query
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
	account, err := l.Cache.Get(address)
	if err != nil {
		return l.Refresh(address)
	}

	acc := account.(sdk.BaseAccount)
	acc.Sequence += 1
	l.saveAccount(acc)
	return acc, nil
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
	l.Info().
		Str("address", address).
		Msg("query account from chain")
	return account, nil
}

func (l localAccount) QueryAddress(name string) (sdk.AccAddress, sdk.Error) {
	address, err := l.keyManager.Query(name)
	if err != nil {
		l.Err(err).
			Str("name", name).
			Msg("can't find account")
		return address, sdk.Wrap(err)
	}
	return address, nil
}

func (l localAccount) saveAccount(account sdk.BaseAccount) {
	address := account.Address.String()
	if err := l.SetWithExpire(address, account, l.expiration); err != nil {
		l.Err(err).
			Str("address", address).
			Msg("save or update cache failed")
		return
	}
	l.Info().
		Str("address", address).
		Msgf("cache account %s", l.expiration.String())
}
