package modules

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub-sdk-go/tools/cache"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type localToken struct {
	q sdk.Queries
	*log.Logger
	cache.Cache
}

func (l localToken) QueryToken(symbol string) (sdk.Token, error) {
	symbol = strings.ToLower(symbol)
	if symbol == sdk.IRIS.Symbol || symbol == sdk.IRIS.MinUnit {
		return sdk.IRIS, nil
	}

	token, err := l.Get(l.keyWithPrefix(symbol))
	if err == nil {
		return token.(sdk.Token), nil
	}

	param := struct {
		Symbol string
	}{
		Symbol: symbol,
	}

	symbol = strings.TrimSuffix(symbol, "-min")
	var t sdk.Token
	if err := l.q.QueryWithResponse("custom/asset/token", param, &t); err != nil {
		return sdk.Token{}, err
	}

	l.SaveTokens(t)
	return t, nil
}

func (l localToken) SaveTokens(tokens ...sdk.Token) {
	for _, t := range tokens {
		err1 := l.Set(l.keyWithPrefix(t.Symbol), t)
		err2 := l.Set(l.keyWithPrefix(t.GetMinUnit()), t)
		if err1 != nil || err2 != nil {
			l.Warn().
				Str("symbol", t.Symbol).
				Msg("cache token failed")
		}
	}
}

func (l localToken) ToMinCoin(coins ...sdk.DecCoin) (dstCoins sdk.Coins, err sdk.Error) {
	for _, coin := range coins {
		token, err := l.QueryToken(coin.Denom)
		if err != nil {
			return nil, sdk.Wrap(err)
		}

		minCoin, err := token.GetCoinType().ConvertToMinCoin(coin)
		if err != nil {
			return nil, sdk.Wrap(err)
		}
		dstCoins = append(dstCoins, minCoin)
	}
	return dstCoins.Sort(), nil
}

func (l localToken) ToMainCoin(coins ...sdk.Coin) (dstCoins sdk.DecCoins, err sdk.Error) {
	for _, coin := range coins {
		token, err := l.QueryToken(coin.Denom)
		if err != nil {
			return dstCoins, sdk.Wrap(err)
		}

		mainCoin, err := token.GetCoinType().ConvertToMainCoin(coin)
		if err != nil {
			return dstCoins, sdk.Wrap(err)
		}
		dstCoins = append(dstCoins, mainCoin)
	}
	return dstCoins.Sort(), nil
}

func (l localToken) keyWithPrefix(symbol string) string {
	return fmt.Sprintf("token:%s", symbol)
}
