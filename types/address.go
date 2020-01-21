package types

import (
	"bytes"
	"encoding/json"
	"github.com/irisnet/irishub-sdk-go/utils"
	"github.com/tendermint/tendermint/libs/bech32"
)

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address string) (addr AccAddress, err error) {
	bech32PrefixAccAddr := GetAddrPrefixCfg().GetBech32AccountAddrPrefix()
	bz, err := utils.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}

	return AccAddress(bz), nil
}

func MustAccAddressFromBech32(address string) AccAddress {
	addr,err := AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return addr
}

// String implements the Stringer interface.
func (aa AccAddress) String() string {
	bech32PrefixAccAddr := GetAddrPrefixCfg().GetBech32AccountAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixAccAddr, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Returns boolean for whether two AccAddresses are Equal
func (aa AccAddress) Equals(aa2 AccAddress) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Compare(aa.Bytes(), aa2.Bytes()) == 0
}

// Returns boolean for whether an AccAddress is empty
func (aa AccAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := AccAddress{}
	return bytes.Compare(aa.Bytes(), aa2.Bytes()) == 0
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa AccAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *AccAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *AccAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AccAddressFromBech32(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa AccAddress) Bytes() []byte {
	return aa
}
