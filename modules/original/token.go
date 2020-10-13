package original

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"strings"

	"github.com/irisnet/irishub-sdk-go/utils/cache"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type tokenQuery struct {
	q original.Queries
	*log.Logger
	cache.Cache
}

func (l tokenQuery) QueryToken(symbol string) (original.Token, error) {
	symbol = strings.ToLower(symbol)
	token, err := l.Get(l.prefixKey(symbol))
	if err == nil {
		return token.(original.Token), nil
	}

	param := struct {
		Symbol string
	}{
		Symbol: symbol,
	}

	//symbol = strings.TrimSuffix(symbol, "-min")
	var t original.Token
	if err := l.q.QueryWithResponse("custom/token/token", param, &t); err != nil {
		return original.Token{}, err
	}

	l.SaveTokens(t)
	return t, nil
}

func (l tokenQuery) SaveTokens(tokens ...original.Token) {
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

func (l tokenQuery) ToMinCoin(coins ...original.DecCoin) (dstCoins original.Coins, err original.Error) {
	for _, coin := range coins {
		token, err := l.QueryToken(coin.Denom)
		if err != nil {
			return nil, original.Wrap(err)
		}

		minCoin, err := token.GetCoinType().ConvertToMinCoin(coin)
		if err != nil {
			return nil, original.Wrap(err)
		}
		dstCoins = append(dstCoins, minCoin)
	}
	return dstCoins.Sort(), nil
}

func (l tokenQuery) ToMainCoin(coins ...original.Coin) (dstCoins original.DecCoins, err original.Error) {
	for _, coin := range coins {
		token, err := l.QueryToken(coin.Denom)
		if err != nil {
			return dstCoins, original.Wrap(err)
		}

		mainCoin, err := token.GetCoinType().ConvertToMainCoin(coin)
		if err != nil {
			return dstCoins, original.Wrap(err)
		}
		dstCoins = append(dstCoins, mainCoin)
	}
	return dstCoins.Sort(), nil
}

func (l tokenQuery) prefixKey(symbol string) string {
	return fmt.Sprintf("token:%s", symbol)
}
