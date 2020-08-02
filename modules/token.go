package modules

import (
	"fmt"
	"strings"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/cache"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type tokenQuery struct {
	q sdk.Queries
	*log.Logger
	cache.Cache
}

func (l tokenQuery) QueryToken(symbol string) (sdk.Token, error) {
	symbol = strings.ToLower(symbol)
	token, err := l.Get(l.prefixKey(symbol))
	if err == nil {
		return token.(sdk.Token), nil
	}

	param := struct {
		Symbol string
	}{
		Symbol: symbol,
	}

	//symbol = strings.TrimSuffix(symbol, "-min")
	var t sdk.Token
	if err := l.q.QueryWithResponse("custom/token/token", param, &t); err != nil {
		return sdk.Token{}, err
	}

	l.SaveTokens(t)
	return t, nil
}

func (l tokenQuery) SaveTokens(tokens ...sdk.Token) {
	for _, t := range tokens {
		err1 := l.Set(l.prefixKey(t.Symbol), t)
		err2 := l.Set(l.prefixKey(t.GetMinUnit()), t)
		if err1 != nil || err2 != nil {
			l.Warn().
				Str("symbol", t.Symbol).
				Msg("cache token failed")
		}
	}
}

func (l tokenQuery) ToMinCoin(coins ...sdk.DecCoin) (dstCoins sdk.Coins, err sdk.Error) {
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

func (l tokenQuery) ToMainCoin(coins ...sdk.Coin) (dstCoins sdk.DecCoins, err sdk.Error) {
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

func (l tokenQuery) prefixKey(symbol string) string {
	return fmt.Sprintf("token:%s", symbol)
}
