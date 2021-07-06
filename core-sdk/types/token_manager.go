package types

type DefaultTokenManager struct{}

func (TokenManager DefaultTokenManager) QueryToken(denom string) (Token, error) {
	return Token{}, nil
}

func (TokenManager DefaultTokenManager) SaveTokens(tokens ...Token) {
	return
}

func (TokenManager DefaultTokenManager) ToMinCoin(coins ...DecCoin) (Coins, Error) {
	for i := range coins {
		if coins[i].Denom == "iris" {
			coins[i].Denom = "uiris"
			coins[i].Amount = coins[i].Amount.MulInt(NewIntWithDecimal(1, 6))
		}
	}
	ucoins, _ := DecCoins(coins).TruncateDecimal()
	return ucoins, nil
}

func (TokenManager DefaultTokenManager) ToMainCoin(coins ...Coin) (DecCoins, Error) {
	decCoins := make(DecCoins, len(coins), 0)
	for _, coin := range coins {
		if coin.Denom == "uiris" {
			amtount := NewDecFromInt(coin.Amount).Mul(NewDecWithPrec(1, 6))
			decCoins = append(decCoins, NewDecCoinFromDec("iris", amtount))
		}
	}
	return decCoins, nil
}
