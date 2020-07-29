package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"

	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub-sdk-go/utils/bech32"
)

const (
	// AddrLen defines a valid address length
	AddrLen = 20
)

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address string) (AccAddress, Error) {
	bech32PrefixAccAddr := GetAddrPrefixCfg().GetBech32AccountAddrPrefix()
	bz, err := bech32.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return nil, Wrap(err)
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

// Returns boolean for whether two AccAddresses are equal
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
func ValAddressFromBech32(address string) (ValAddress, Error) {
	bech32PrefixValAddr := GetAddrPrefixCfg().GetBech32ValidatorAddrPrefix()
	bz, err := bech32.GetFromBech32(address, bech32PrefixValAddr)
	if err != nil {
		return nil, Wrap(err)
	}

	return ValAddress(bz), nil
}

// Returns boolean for whether two ValAddresses are equal
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

// ConsAddress defines a wrapper around bytes meant to present a consensus node.
// When marshaled to a string or JSON, it uses Bech32.
type ConsAddress []byte

// String implements the Stringer interface.
func (ca ConsAddress) String() string {
	bech32PrefixConsAddr := GetAddrPrefixCfg().GetBech32ConsensusAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixConsAddr, ca.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Bytes returns the raw address bytes.
func (ca ConsAddress) Bytes() []byte {
	return ca
}

// ConsAddressFromHex creates a ConsAddress from a hex string.
func ConsAddressFromHex(address string) (addr ConsAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return ConsAddress(bz), nil
}

// ConsAddressFromBech32 creates a ConsAddress from a Bech32 string.
func ConsAddressFromBech32(address string) (addr ConsAddress, err error) {
	bech32PrefixConsAddr := GetAddrPrefixCfg().GetBech32ConsensusAddrPrefix()
	bz, err := GetFromBech32(address, bech32PrefixConsAddr)
	if err != nil {
		return nil, err
	}

	return ConsAddress(bz), nil
}

// Bech32ifyConsPub returns a Bech32 encoded string containing the
// Bech32PrefixConsPub prefixfor a given consensus node's PubKey.
func Bech32ifyConsPub(pub crypto.PubKey) (string, error) {
	bech32PrefixConsPub := GetAddrPrefixCfg().GetBech32ConsensusPubPrefix()
	return bech32.ConvertAndEncode(bech32PrefixConsPub, pub.Bytes())
}

// GetConsPubKeyBech32 creates a PubKey for a consensus node with a given public
// key string using the Bech32 Bech32PrefixConsPub prefix.
func GetConsPubKeyBech32(pubkey string) (pk crypto.PubKey, err error) {
	bech32PrefixConsPub := GetAddrPrefixCfg().GetBech32ConsensusPubPrefix()
	bz, err := GetFromBech32(pubkey, bech32PrefixConsPub)
	if err != nil {
		return nil, err
	}

	pk, err = cryptoAmino.PubKeyFromBytes(bz)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

// GetFromBech32 decodes a bytestring from a Bech32 encoded string.
func GetFromBech32(bech32str, prefix string) ([]byte, error) {
	if len(bech32str) == 0 {
		return nil, errors.New("decoding Bech32 address failed: must provide an address")
	}

	hrp, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, err
	}

	if hrp != prefix {
		return nil, fmt.Errorf("invalid Bech32 prefix; expected %s, got %s", prefix, hrp)
	}

	return bz, nil
}
