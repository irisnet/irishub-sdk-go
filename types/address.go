package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bech32"

	"github.com/irisnet/irishub-sdk-go/utils"
)

const (
	// AddrLen defines a valid address length
	AddrLen = 20
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
	addr, err := AccAddressFromBech32(address)
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

// ----------------------------------------------------------------------------
// validator operator
// ----------------------------------------------------------------------------

// ValAddress defines a wrapper around bytes meant to present a validator's
// operator. When marshaled to a string or JSON, it uses Bech32.
type ValAddress []byte

// ValAddressFromHex creates a ValAddress from a hex string.
func ValAddressFromHex(address string) (addr ValAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return ValAddress(bz), nil
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromBech32(address string) (addr ValAddress, err error) {
	bech32PrefixValAddr := GetAddrPrefixCfg().GetBech32ValidatorAddrPrefix()
	bz, err := utils.GetFromBech32(address, bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}

	return ValAddress(bz), nil
}

// Returns boolean for whether two ValAddresses are Equal
func (va ValAddress) Equals(va2 ValAddress) bool {
	if va.Empty() && va2.Empty() {
		return true
	}

	return bytes.Compare(va.Bytes(), va2.Bytes()) == 0
}

// Returns boolean for whether an AccAddress is empty
func (va ValAddress) Empty() bool {
	if va == nil {
		return true
	}

	va2 := ValAddress{}
	return bytes.Compare(va.Bytes(), va2.Bytes()) == 0
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (va ValAddress) Marshal() ([]byte, error) {
	return va, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (va *ValAddress) Unmarshal(data []byte) error {
	*va = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (va ValAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(va.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (va *ValAddress) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	va2, err := ValAddressFromBech32(s)
	if err != nil {
		return err
	}

	*va = va2
	return nil
}

// Bytes returns the raw address bytes.
func (va ValAddress) Bytes() []byte {
	return va
}

// String implements the Stringer interface.
func (va ValAddress) String() string {
	bech32PrefixValAddr := GetAddrPrefixCfg().GetBech32ValidatorAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixValAddr, va.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (va ValAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(fmt.Sprintf("%s", va.String())))
	case 'p':
		_, _ = s.Write([]byte(fmt.Sprintf("%p", va)))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%X", []byte(va))))
	}
}

// Bech32ifyConsPub returns a Bech32 encoded string containing the
// Bech32PrefixConsPub prefixfor a given consensus node's PubKey.
func Bech32ifyConsPub(pub crypto.PubKey) (string, error) {
	bech32PrefixConsPub := GetAddrPrefixCfg().GetBech32ConsensusPubPrefix()
	return bech32.ConvertAndEncode(bech32PrefixConsPub, pub.Bytes())
}
