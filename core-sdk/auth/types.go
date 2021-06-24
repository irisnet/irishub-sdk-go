package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	commoncodec "github.com/irisnet/irishub-sdk-go/common/codec"
	codectypes "github.com/irisnet/irishub-sdk-go/common/codec/types"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/crypto"
)

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.

type Account interface {
	GetAddress() sdk.AccAddress
	SetAddress(sdk.AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error
}

//BaseAccount Have they all been implemented
var _ Account = (*BaseAccount)(nil)

// GetAddress - Implements sdk.AccountI.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(acc.Address)
	return addr
}

// SetAddress - Implements sdk.AccountI.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr.String()
	return nil
}

// GetPubKey - Implements sdk.AccountI.
func (acc BaseAccount) GetPubKey() (pk crypto.PubKey) {
	if acc.PubKey == nil {
		return nil
	}
	content, ok := acc.PubKey.GetCachedValue().(crypto.PubKey)
	if !ok {
		return nil
	}
	return content
}

// SetPubKey - Implements sdk.AccountI.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	if pubKey == nil {
		acc.PubKey = nil
	} else {
		protoMsg, ok := pubKey.(proto.Message)
		if !ok {
			return sdk.Wrap(fmt.Errorf("err invalid key, can't proto encode %T", protoMsg))
		}

		any, err := codectypes.NewAnyWithValue(protoMsg)
		if err != nil {
			return err
		}
		acc.PubKey = any
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

//json.Marshal BaseAccount
func (acc BaseAccount) String() string {
	out, _ := json.Marshal(acc)
	return string(out)
}

// Convert return a sdk.BaseAccount
func (acc *BaseAccount) Convert() interface{} {
	// error don't use it
	return nil
}

// Convert return a sdk.BaseAccount
// in order to unpack pubKey so not use Convert()
func (acc *BaseAccount) ConvertAccount(cdc commoncodec.Marshaler) interface{} {
	account := sdk.BaseAccount{
		Address:       acc.Address,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	var pkStr string
	if acc.PubKey == nil {
		return account
	}

	var pk crypto.PubKey
	if err := cdc.UnpackAny(acc.PubKey, &pk); err != nil {
		return sdk.BaseAccount{}
	}

	pkStr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pk)
	if err != nil {
		panic(err)
	}

	account.PubKey = pkStr
	return account
}
