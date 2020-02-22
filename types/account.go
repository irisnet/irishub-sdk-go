package types

import (
	"errors"

	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/tendermint/tendermint/crypto"
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

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	GetCoins() Coins
	SetCoins(Coins) error

	GetMemoRegexp() string
	SetMemoRegexp(string)
}

var _ Account = (*BaseAccount)(nil)

type BaseAccount struct {
	Address       AccAddress    `json:"address"`
	Coins         Coins         `json:"coins"`
	PubKey        crypto.PubKey `json:"public_key"`
	AccountNumber uint64        `json:"account_number"`
	Sequence      uint64        `json:"sequence"`
	MemoRegexp    string        `json:"memo_regexp"`
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
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetCoins() Coins {
	return acc.Coins
}

// Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins Coins) error {
	acc.Coins = coins
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
	return acc.MemoRegexp
}

// Implements sdk.Account.
func (acc *BaseAccount) SetMemoRegexp(regexp string) {
	acc.MemoRegexp = regexp
}

func RegisterCodec(cdc Codec) {
	cdc.RegisterInterface((*Account)(nil))
	cdc.RegisterInterface((*Msg)(nil))
	cdc.RegisterConcrete(&BaseAccount{}, "irishub/bank/Account")
	cdc.RegisterConcrete(StdTx{}, "irishub/bank/StdTx")
	// These are all written here instead of
	cdc.RegisterInterface((*crypto.PubKey)(nil))
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoName)
	cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{},
		multisig.PubKeyMultisigThresholdAminoRoute)

	cdc.RegisterInterface((*crypto.PrivKey)(nil))
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{},
		ed25519.PrivKeyAminoName)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{},
		secp256k1.PrivKeyAminoName)

	//cryptoAmino.RegisterAmino(cdc)
}
