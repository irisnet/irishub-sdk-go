package types

import (
	"errors"
)

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.
type Account interface {
	GetAddress() AccAddress
	SetAddress(AccAddress) error // errors if already set.

	GetPubKey() []byte // can return nil.
	SetPubKey([]byte) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	//GetCoins() Coins
	//SetCoins(Coins) error

	GetMemoRegexp() string
	SetMemoRegexp(string)
}

var _ Account = (*BaseAccount)(nil)

/*type BaseAccount struct {
	Address       AccAddress    `json:"address"`
	Coins         Coins         `json:"coins"`
	PubKey        crypto.PubKey `json:"public_key"`
	AccountNumber uint64        `json:"account_number"`
	Sequence      uint64        `json:"sequence"`
	MemoRegexp    string        `json:"memo_regexp"`
}*/

type BaseAccount struct {
	Address       AccAddress `json:"address"`
	PubKey        []byte     `json:"public_key"`
	AccountNumber uint64     `json:"account_number"`
	Sequence      uint64     `json:"sequence"`
}

// Implements sdk.Account.
func (acc BaseAccount) GetAddress() AccAddress {
	return acc.Address
}

// Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// Implements sdk.Account.
func (acc BaseAccount) GetPubKey() []byte {
	return acc.PubKey
}

// Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey []byte) error {
	acc.PubKey = pubKey
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetCoins() string {
	return ""
}

// Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins Coins) error {
	acc.PubKey = []byte("coins")
	return nil
}

// Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetMemoRegexp() string {
	return "acc.MemoRegexp"
}

// Implements sdk.Account.
func (acc *BaseAccount) SetMemoRegexp(regexp string) {
	//acc.MemoRegexp = regexp
}

func (acc *BaseAccount) Convert() interface{} {
	return acc
}
