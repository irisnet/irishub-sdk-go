package types

type Token struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Scale         uint32 `json:"scale"`
	MinUnit       string `json:"min_unit"`
	InitialSupply uint64 `json:"initial_supply"`
	MaxSupply     uint64 `json:"max_supply"`
	Mintable      bool   `json:"mintable"`
	Owner         string `json:"owner"`
}

// GetCoinType returns CnType
func (t Token) GetCoinType() CoinType {
	return CoinType{
		Name:     t.Name,
		MinUnit:  NewUnit(t.MinUnit, uint8(t.Scale)),
		MainUnit: NewUnit(t.Symbol, 0),
		Desc:     t.Name,
	}
}

type Tokens []Token
