package types

import (
	"bytes"
	"errors"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
)

// AccountI is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements AccountI.
type AccountI interface {
	GetAddress() AccAddress
	SetAddress(AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error
}

var _ AccountI = (*BaseAccount)(nil)

type BaseAccount struct {
	Address       AccAddress `json:"address"`
	PubKey        []byte     `json:"public_key"`
	AccountNumber uint64     `json:"account_number"`
	Sequence      uint64     `json:"sequence"`
}

// GetAddress - Implements sdk.AccountI.
func (acc BaseAccount) GetAddress() AccAddress {
	return acc.Address
}

// SetAddress - Implements sdk.AccountI.
func (acc *BaseAccount) SetAddress(addr AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}

	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.AccountI.
func (acc BaseAccount) GetPubKey() (pk crypto.PubKey) {
	if len(acc.PubKey) == 0 {
		return nil
	}

	amino.MustUnmarshalBinaryBare(acc.PubKey, &pk)
	return pk
}

// SetPubKey - Implements sdk.AccountI.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	if pubKey == nil {
		acc.PubKey = nil
	} else {
		acc.PubKey = pubKey.Bytes()
	}

	return nil
}

// GetAccountNumber - Implements AccountI
func (acc BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements AccountI
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.AccountI.
func (acc BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.AccountI.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if len(acc.PubKey) != 0 && acc.Address != nil &&
		!bytes.Equal(acc.GetPubKey().Address().Bytes(), acc.Address.Bytes()) {
		return errors.New("account address and pubkey address do not match")
	}

	return nil
}

func (acc BaseAccount) Convert() interface{} {
	return acc
}
