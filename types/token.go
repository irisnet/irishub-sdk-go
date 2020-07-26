package types

import (
	"fmt"
	"strings"
)

var (
	IRIS = Token{
		Symbol:        iris,
		Name:          "IRIS Network",
		Scale:         18,
		MinUnit:       irisAtto,
		InitialSupply: 2000000000,
		MaxSupply:     1000000000000,
		Mintable:      true,
		Owner:         "",
	}
)

type TokenResponse struct {
	TypeName string `json:"type"`
	Value    Token  `json:"value"`
}

type Token struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Scale         uint8  `json:"scale"`
	MinUnit       string `json:"min_unit"`
	InitialSupply uint64 `json:"initial_supply"`
	MaxSupply     uint64 `json:"max_supply"`
	Mintable      bool   `json:"mintable"`
	Owner         string `json:"owner"`
}

// GetMinUnit returns MinUnit
func (t Token) GetMinUnit() string {
	symbol := strings.ToLower(strings.TrimSpace(t.Symbol))

	if symbol == IRIS.Symbol {
		return IRIS.MinUnit
	}

	return fmt.Sprintf("%s%s", symbol, minDenomSuffix)
}

// GetCoinType returns CoinType
func (t Token) GetCoinType() CoinType {
	return CoinType{
		Name:     t.Name,
		MinUnit:  NewUnit(t.GetMinUnit(), t.Scale),
		MainUnit: NewUnit(t.Symbol, 0),
		Desc:     t.Name,
	}
}

func (t Token) Convert() interface{} {
	return t
}

type Tokens []TokenResponse

func (t Tokens) Convert() interface{} {
	return t
}
