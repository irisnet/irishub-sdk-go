package types

type PageRequest struct {
	// key is a value returned in PageResponse.next_key to begin
	// querying the next page most efficiently. Only one of offset or key
	// should be set.
	Key []byte `json:"key,omitempty"`
	// offset is a numeric offset that can be used when key is unavailable.
	// It is less efficient than using key. Only one of offset or key should
	// be set.
	Offset uint64 ` json:"offset,omitempty"`
	// limit is the total number of results to be returned in the result page.
	// If left empty it will default to a value to be set by each app.
	Limit uint64 ` json:"limit,omitempty"`
	// count_total is set to true  to indicate that the result set should include
	// a count of the total number of items available for pagination in UIs.
	// count_total is only respected when offset is used. It is ignored when key
	// is set.
	CountTotal bool ` json:"count_total,omitempty"`
}

type PoolInfo struct {
	Id string `json:"id"`
	// escrow account for deposit tokens
	EscrowAddress string `json:"escrow_address"`
	// main token balance
	Standard Coin `json:"standard"`
	// counterparty token balance
	Token Coin `json:"token"`
	// liquidity token balance
	Lpt Coin `json:"lpt"`
	// liquidity pool fee
	Fee string `json:"fee"`
}
