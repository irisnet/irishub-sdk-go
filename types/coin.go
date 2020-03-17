package types

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	// Denominations can be 3 ~ 21 characters long.
	reABS           = `([a-z][0-9a-z]{2}[:])?`
	reCoinName      = reABS + `(([a-z][a-z0-9]{2,7}|x)\.)?([a-z][a-z0-9]{2,7})`
	reDenom         = reCoinName + `(-[a-z]{3,5})?`
	reAmount        = `[0-9]+(\.[0-9]+)?`
	reSpace         = `[[:space:]]*`
	reDenomCompiled = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDenom))
	reCoinCompiled  = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmount, reSpace, reDenom))

	reDecAmt  = `^(0|([1-9]*))(\.\d+)?$`
	reSpc     = `[[:space:]]*`
	reDecCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmount, reSpc, reDenom))

	irisAtto       = "iris-atto"
	minDenomSuffix = "-min"
	iris           = "iris"
)

type Coin struct {
	Denom string `json:"denom"`
	// To allow the use of unsigned integers (see: #1273) a larger refactor will
	// need to be made. So we use signed integers for now with safety measures in
	// place preventing negative values being used.
	Amount Int `json:"amount"`
}

// NewCoin returns a new coin with a denomination and amount. It will panic if
// the amount is negative.
func NewCoin(denom string, amount Int) Coin {
	if amount.i == nil {
		amount = ZeroInt()
	}

	if amount.IsNegative() {
		panic("negative coin amount")
	}

	return Coin{
		Denom:  denom,
		Amount: amount,
	}
}

// String provides a human-readable representation of a coin
func (coin Coin) String() string {
	return fmt.Sprintf("%v%v", coin.Amount, coin.Denom)
}

// IsPositive returns true if coin amount is positive.
//
func (coin Coin) IsPositive() bool {
	return coin.Amount.i != nil && coin.Amount.IsPositive()
}

// IsNegative returns true if the coin amount is negative and false otherwise.
//
func (coin Coin) IsNegative() bool {
	return coin.Amount.i != nil && coin.Amount.IsNegative()
}

func (coin Coin) IsValidIrisAtto() bool {
	return coin.Denom == "iris-atto" && coin.IsPositive()
}

// IsZero returns if this coin has zero amount
func (coin Coin) IsZero() bool {
	return coin.Amount.i == nil || coin.Amount.IsZero()
}

// IsValid returns true if the coin amount is non-negative
// and the coin is denominated in its minimum unit
func (coin Coin) IsValid() bool {
	if coin.IsNegative() {
		return false
	}
	return IsCoinMinDenomValid(coin.Denom)
}

// IsEqual returns true if the two sets of Coins have the same value
func (coin Coin) IsEqual(other Coin) bool {
	return coin.Denom == other.Denom && coin.Amount.Equal(other.Amount)
}

func IsCoinMinDenomValid(denom string) bool {
	if denom != irisAtto && (!strings.HasSuffix(denom, minDenomSuffix) || strings.HasPrefix(denom, iris+"-")) {
		return false
	}
	return reDenomCompiled.MatchString(denom)
}

// Adds amounts of two coins with same denom. If the coins differ in denom then
// it panics.
func (coin Coin) Add(coinB Coin) Coin {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, coinB.Denom))
	}

	return Coin{coin.Denom, coin.Amount.Add(coinB.Amount)}
}

func ParseCoin(coinStr string) (coin Coin, err error) {
	coinStr = strings.ToLower(strings.TrimSpace(coinStr))

	matches := reCoinCompiled.FindStringSubmatch(coinStr)
	if matches == nil {
		return Coin{}, fmt.Errorf("invalid coin expression: %s", coinStr)
	}

	denomStr, amountStr := matches[3], matches[1]

	amount, ok := NewIntFromString(amountStr)
	if !ok {
		return Coin{}, fmt.Errorf("failed to parse coin amount: %s", amountStr)
	}

	return NewCoin(denomStr, amount), nil
}

// Coins is a set of Coin, one per currency
type Coins []Coin

// NewCoins constructs a new coin set.
func NewCoins(coins ...Coin) Coins {
	// remove zeroes
	newCoins := removeZeroCoins(Coins(coins))
	if len(newCoins) == 0 {
		return Coins{}
	}

	newCoins.Sort()

	if !newCoins.IsValid() {
		panic(fmt.Sprintf("invalid coin set: %s", newCoins))
	}

	return newCoins
}

// Empty returns true if there are no coins and false otherwise.
func (coins Coins) Empty() bool {
	return len(coins) == 0
}

// IsEqual returns true if the two sets of Coins have the same value
func (coins Coins) IsEqual(coinsB Coins) bool {
	if len(coins) != len(coinsB) {
		return false
	}

	coins = coins.Sort()
	coinsB = coinsB.Sort()

	for i := 0; i < len(coins); i++ {
		if !coins[i].IsEqual(coinsB[i]) {
			return false
		}
	}

	return true
}

// IsValid asserts the coins are valid and sorted.
func (coins Coins) IsValid() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		return coins[0].IsValid()
	default:
		// check first coin
		if !coins[0].IsValid() {
			return false
		}

		lowDenom := coins[0].Denom
		for _, coin := range coins[1:] {
			if !coin.IsValid() {
				return false
			}
			if coin.Denom <= lowDenom {
				return false
			}

			// we compare each coin against the last denom
			lowDenom = coin.Denom
		}

		return true
	}
}

// Add adds two sets of coins.
//
// e.g.
// {2A} + {A, 2B} = {3A, 2B}
// {2A} + {0B} = {2A}
//
// NOTE: Add operates under the invariant that coins are sorted by
// denominations.
//
// CONTRACT: Add will never return Coins where one Coin has a negative
// amount. In other words, IsValid will always return true.
func (coins Coins) Add(coinsB ...Coin) Coins {
	sum, hasNeg := coins.SafeAdd(coinsB)
	if hasNeg {
		panic("negative coin amount")
	}

	return sum
}

// SafeAdd performs the same arithmetic as Add but returns a boolean if any
// negative coin amount was returned.
func (coins Coins) SafeAdd(coinsB Coins) (Coins, bool) {
	sum := coins.safeAdd(coinsB)
	return sum, sum.IsAnyNegative()
}

// safeAdd will perform addition of two coins sets. If both coin sets are
// empty, then an empty set is returned. If only a single set is empty, the
// other set is returned. Otherwise, the coins are compared in order of their
// denomination and addition only occurs when the denominations match, otherwise
// the coin is simply added to the sum assuming it's not zero.
func (coins Coins) safeAdd(coinsB Coins) Coins {
	sum := ([]Coin)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(coins), len(coinsB)

	for {
		if indexA == lenA {
			if indexB == lenB {
				// return nil coins if both sets are empty
				return sum
			}

			// return set B (excluding zero coins) if set A is empty
			return append(sum, removeZeroCoins(coinsB[indexB:])...)
		} else if indexB == lenB {
			// return set A (excluding zero coins) if set B is empty
			return append(sum, removeZeroCoins(coins[indexA:])...)
		}

		coinA, coinB := coins[indexA], coinsB[indexB]

		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1: // coin A denom < coin B denom
			if !coinA.IsZero() {
				sum = append(sum, coinA)
			}

			indexA++

		case 0: // coin A denom == coin B denom
			res := coinA.Add(coinB)
			if !res.IsZero() {
				sum = append(sum, res)
			}

			indexA++
			indexB++

		case 1: // coin A denom > coin B denom
			if !coinB.IsZero() {
				sum = append(sum, coinB)
			}

			indexB++
		}
	}
}

// IsAnyNegative returns true if at least one coin has negative amount.
//
func (coins Coins) IsAnyNegative() bool {
	for _, coin := range coins {
		if coin.IsNegative() {
			return true
		}
	}

	return false
}

// removeZeroCoins removes all zero coins from the given coin set in-place.
func removeZeroCoins(coins Coins) Coins {
	i, l := 0, len(coins)
	for i < l {
		if coins[i].IsZero() {
			// remove coin
			coins = append(coins[:i], coins[i+1:]...)
			l--
		} else {
			i++
		}
	}

	return coins[:i]
}

func (coins Coins) String() string {
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v,", coin.String())
	}
	return out[:len(out)-1]
}

//-----------------------------------------------------------------------------
// Sort interface

//nolint
func (coins Coins) Len() int           { return len(coins) }
func (coins Coins) Less(i, j int) bool { return coins[i].Denom < coins[j].Denom }
func (coins Coins) Swap(i, j int)      { coins[i], coins[j] = coins[j], coins[i] }

var _ sort.Interface = Coins{}

// Sort is a helper function to sort the set of coins inplace
func (coins Coins) Sort() Coins {
	sort.Sort(coins)
	return coins
}

// validate returns an error if the Coin has a negative amount or if
// the denom is invalid.
func validate(denom string, amount Int) error {
	if err := ValidateDenom(denom); err != nil {
		return err
	}

	if amount.IsNegative() {
		return fmt.Errorf("negative coin amount: %v", amount)
	}

	return nil
}

// ValidateDenom validates a denomination string returning an error if it is
// invalid.
func ValidateDenom(denom string) error {
	if !reDenomCompiled.MatchString(denom) {
		return fmt.Errorf("invalid denom: %s", denom)
	}
	return nil
}

func mustValidateDenom(denom string) {
	if err := ValidateDenom(denom); err != nil {
		panic(err)
	}
}

// ParseCoins will parse out a list of coins separated by commas.
// If nothing is provided, it returns nil Coins.
// Returned coins are sorted.
func ParseCoins(coinsStr string) (coins Coins, err error) {
	if len(coinsStr) == 0 {
		return Coins{}, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin, err := ParseCoin(coinStr)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	// Sort coins for determinism.
	coins.Sort()

	return coins, nil
}

type findDupDescriptor interface {
	GetDenomByIndex(int) string
	Len() int
}

// findDup works on the assumption that coins is sorted
func findDup(coins findDupDescriptor) int {
	if coins.Len() <= 1 {
		return -1
	}

	prevDenom := coins.GetDenomByIndex(0)
	for i := 1; i < coins.Len(); i++ {
		if coins.GetDenomByIndex(i) == prevDenom {
			return i
		}
		prevDenom = coins.GetDenomByIndex(i)
	}

	return -1
}
