package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

// NOTE: never use new(Decimal) or else we will panic unmarshalling into the
// nil embedded big.Int
type Decimal struct {
	i *big.Int
}

// number of decimal places
const (
	precision = 18

	// bytes required to represent the above precision
	// Ceiling[Log2[999 999 999 999 999 999]]
	decimalPrecisionBits = 64
)

var (
	precisionReuseV1       = new(big.Int).Exp(big.NewInt(10), big.NewInt(precision), nil)
	fivePrecisionV1        = new(big.Int).Quo(precisionReuseV1, big.NewInt(2))
	precisionMultipliersV1 []*big.Int
)

// Set precision multipliers
func init() {
	precisionMultipliersV1 = make([]*big.Int, precision+1)
	for i := 0; i <= precision; i++ {
		precisionMultipliersV1[i] = calcprecisionMultiplier(int64(i))
	}
}

func precisionIntV1() *big.Int {
	return new(big.Int).Set(precisionReuseV1)
}

func ZeroDecimal() Decimal     { return Decimal{new(big.Int).Set(zeroInt)} }
func OneDecimal() Decimal      { return Decimal{precisionIntV1()} }
func SmallestDecimal() Decimal { return Decimal{new(big.Int).Set(oneInt)} }

// calculate the precision multiplier
func calcprecisionMultiplier(prec int64) *big.Int {
	if prec > precision {
		panic(fmt.Sprintf("too much precision, maximum %v, provided %v", precision, prec))
	}
	zerosToAdd := precision - prec
	multiplier := new(big.Int).Exp(tenInt, big.NewInt(zerosToAdd), nil)
	return multiplier
}

// get the precision multiplier, do not mutate result
func precisionMultiplierV1(prec int64) *big.Int {
	if prec > precision {
		panic(fmt.Sprintf("too much precision, maximum %v, provided %v", precision, prec))
	}
	return precisionMultipliersV1[prec]
}

//______________________________________________________________________________________________

// create a new Decimal from integer assuming whole number
func NewDecimal(i int64) Decimal {
	return NewDecimalWithPrec(i, 0)
}

// create a new Decimal from integer with decimal place at prec
// CONTRACT: prec <= precision
func NewDecimalWithPrec(i, prec int64) Decimal {
	return Decimal{
		new(big.Int).Mul(big.NewInt(i), precisionMultiplierV1(prec)),
	}
}

// create a new Decimal from big integer assuming whole numbers
// CONTRACT: prec <= precision
func NewDecimalFromBigInt(i *big.Int) Decimal {
	return NewDecimalFromBigIntWithPrec(i, 0)
}

// create a new Decimal from big integer assuming whole numbers
// CONTRACT: prec <= precision
func NewDecimalFromBigIntWithPrec(i *big.Int, prec int64) Decimal {
	return Decimal{
		new(big.Int).Mul(i, precisionMultiplierV1(prec)),
	}
}

// create a new Decimal from big integer assuming whole numbers
// CONTRACT: prec <= precision
func NewDecimalFromInt(i Int) Decimal {
	return NewDecimalFromIntWithPrec(i, 0)
}

// create a new Decimal from big integer with decimal place at prec
// CONTRACT: prec <= precision
func NewDecimalFromIntWithPrec(i Int, prec int64) Decimal {
	return Decimal{
		new(big.Int).Mul(i.BigInt(), precisionMultiplierV1(prec)),
	}
}

// create a decimal from an input decimal string.
// valid must come in the form:
//   (-) whole integers (.) decimal integers
// examples of acceptable input include:
//   -123.456
//   456.7890
//   345
//   -456789
//
// NOTE - An error will return if more decimal places
// are provided in the string than the constant precision.
//
// CONTRACT - This function does not mutate the input str.
func NewDecimalFromStr(str string) (Decimal, error) {
	if len(str) == 0 {
		return Decimal{}, ErrEmptyDecimalStr
	}

	// first extract any negative symbol
	neg := false
	if str[0] == '-' {
		neg = true
		str = str[1:]
	}

	if len(str) == 0 {
		return Decimal{}, ErrEmptyDecimalStr
	}

	strs := strings.Split(str, ".")
	lenDecs := 0
	combinedStr := strs[0]

	if len(strs) == 2 { // has a decimal place
		lenDecs = len(strs[1])
		if lenDecs == 0 || len(combinedStr) == 0 {
			return Decimal{}, ErrInvalidDecimalLength
		}
		combinedStr += strs[1]

	} else if len(strs) > 2 {
		return Decimal{}, ErrInvalidDecimalStr
	}

	if lenDecs > precision {
		return Decimal{}, fmt.Errorf("invalid precision; max: %d, got: %d", precision, lenDecs)
	}

	// add some extra zero's to correct to the precision factor
	zerosToAdd := precision - lenDecs
	zeros := fmt.Sprintf(`%0`+strconv.Itoa(zerosToAdd)+`s`, "")
	combinedStr += zeros

	combined, ok := new(big.Int).SetString(combinedStr, 10) // base 10
	if !ok {
		return Decimal{}, fmt.Errorf("failed to set decimal string: %s", combinedStr)
	}
	if neg {
		combined = new(big.Int).Neg(combined)
	}

	return Decimal{combined}, nil
}

// Decimal from string, panic on error
func MustNewDecimalFromStr(s string) Decimal {
	dec, err := NewDecimalFromStr(s)
	if err != nil {
		panic(err)
	}
	return dec
}

//______________________________________________________________________________________________
//nolint
func (d Decimal) IsNil() bool           { return d.i == nil }                     // is decimal nil
func (d Decimal) IsZero() bool          { return (d.i).Sign() == 0 }              // is equal to zero
func (d Decimal) IsNegative() bool      { return (d.i).Sign() == -1 }             // is negative
func (d Decimal) IsPositive() bool      { return (d.i).Sign() == 1 }              // is positive
func (d Decimal) Equal(d2 Decimal) bool { return (d.i).Cmp(d2.i) == 0 }           // equal decimals
func (d Decimal) GT(d2 Decimal) bool    { return (d.i).Cmp(d2.i) > 0 }            // greater than
func (d Decimal) GTE(d2 Decimal) bool   { return (d.i).Cmp(d2.i) >= 0 }           // greater than or equal
func (d Decimal) LT(d2 Decimal) bool    { return (d.i).Cmp(d2.i) < 0 }            // less than
func (d Decimal) LTE(d2 Decimal) bool   { return (d.i).Cmp(d2.i) <= 0 }           // less than or equal
func (d Decimal) Neg() Decimal          { return Decimal{new(big.Int).Neg(d.i)} } // reverse the decimal sign
func (d Decimal) Abs() Decimal          { return Decimal{new(big.Int).Abs(d.i)} } // absolute value

// BigInt returns a copy of the underlying big.Int.
func (d Decimal) BigInt() *big.Int {
	copy := new(big.Int)
	return copy.Set(d.i)
}

// addition
func (d Decimal) Add(d2 Decimal) Decimal {
	res := new(big.Int).Add(d.i, d2.i)

	if res.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{res}
}

// subtraction
func (d Decimal) Sub(d2 Decimal) Decimal {
	res := new(big.Int).Sub(d.i, d2.i)

	if res.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{res}
}

// multiplication
func (d Decimal) Mul(d2 Decimal) Decimal {
	mul := new(big.Int).Mul(d.i, d2.i)
	chopped := chopprecisionAndRound(mul)

	if chopped.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{chopped}
}

// multiplication truncate
func (d Decimal) MulTruncate(d2 Decimal) Decimal {
	mul := new(big.Int).Mul(d.i, d2.i)
	chopped := chopprecisionAndTruncate(mul)

	if chopped.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{chopped}
}

// multiplication
func (d Decimal) MulInt(i Int) Decimal {
	mul := new(big.Int).Mul(d.i, i.i)

	if mul.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{mul}
}

// MulInt64 - multiplication with int64
func (d Decimal) MulInt64(i int64) Decimal {
	mul := new(big.Int).Mul(d.i, big.NewInt(i))

	if mul.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{mul}
}

// quotient
func (d Decimal) Quo(d2 Decimal) Decimal {

	// multiply precision twice
	mul := new(big.Int).Mul(d.i, precisionReuseV1)
	mul.Mul(mul, precisionReuseV1)

	quo := new(big.Int).Quo(mul, d2.i)
	chopped := chopprecisionAndRound(quo)

	if chopped.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{chopped}
}

// quotient truncate
func (d Decimal) QuoTruncate(d2 Decimal) Decimal {

	// multiply precision twice
	mul := new(big.Int).Mul(d.i, precisionReuseV1)
	mul.Mul(mul, precisionReuseV1)

	quo := new(big.Int).Quo(mul, d2.i)
	chopped := chopprecisionAndTruncate(quo)

	if chopped.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{chopped}
}

// quotient, round up
func (d Decimal) QuoRoundUp(d2 Decimal) Decimal {
	// multiply precision twice
	mul := new(big.Int).Mul(d.i, precisionReuseV1)
	mul.Mul(mul, precisionReuseV1)

	quo := new(big.Int).Quo(mul, d2.i)
	chopped := chopprecisionAndRoundUp(quo)

	if chopped.BitLen() > 255+decimalPrecisionBits {
		panic("Int overflow")
	}
	return Decimal{chopped}
}

// quotient
func (d Decimal) QuoInt(i Int) Decimal {
	mul := new(big.Int).Quo(d.i, i.i)
	return Decimal{mul}
}

// QuoInt64 - quotient with int64
func (d Decimal) QuoInt64(i int64) Decimal {
	mul := new(big.Int).Quo(d.i, big.NewInt(i))
	return Decimal{mul}
}

// is integer, e.g. decimals are zero
func (d Decimal) IsInteger() bool {
	return new(big.Int).Rem(d.i, precisionReuseV1).Sign() == 0
}

// format decimal state
func (d Decimal) Format(s fmt.State, verb rune) {
	_, err := s.Write([]byte(d.String()))
	if err != nil {
		panic(err)
	}
}

func (d Decimal) String() string {
	if d.i == nil {
		return d.i.String()
	}

	isNeg := d.IsNegative()
	if d.IsNegative() {
		d = d.Neg()
	}

	bzInt, err := d.i.MarshalText()
	if err != nil {
		return ""
	}
	inputSize := len(bzInt)

	var bzStr []byte

	// case 1, purely decimal
	if inputSize <= precision {
		bzStr = make([]byte, precision+2)

		// 0. prefix
		bzStr[0] = byte('0')
		bzStr[1] = byte('.')

		// set relevant digits to 0
		for i := 0; i < precision-inputSize; i++ {
			bzStr[i+2] = byte('0')
		}

		// set final digits
		copy(bzStr[2+(precision-inputSize):], bzInt)

	} else {

		// inputSize + 1 to account for the decimal point that is being added
		bzStr = make([]byte, inputSize+1)
		decPointPlace := inputSize - precision

		copy(bzStr, bzInt[:decPointPlace])                   // pre-decimal digits
		bzStr[decPointPlace] = byte('.')                     // decimal point
		copy(bzStr[decPointPlace+1:], bzInt[decPointPlace:]) // post-decimal digits
	}

	// Remove trailing zeros
	res := strings.TrimRightFunc(string(bzStr), func(r rune) bool {
		c := string(r)
		return c == "0" || c == "."
	})

	if isNeg {
		return "-" + res
	}

	return res
}

//     ____
//  __|    |__   "chop 'em
//       ` \     round!"
// ___||  ~  _     -bankers
// |         |      __
// |       | |   __|__|__
// |_____:  /   | $$$    |
//              |________|

// Remove a precision amount of rightmost digits and perform bankers rounding
// on the remainder (gaussian rounding) on the digits which have been removed.
//
// Mutates the input. Use the non-mutative version if that is undesired
func chopprecisionAndRound(d *big.Int) *big.Int {

	// remove the negative and add it back when returning
	if d.Sign() == -1 {
		// make d positive, compute chopped value, and then un-mutate d
		d = d.Neg(d)
		d = chopprecisionAndRound(d)
		d = d.Neg(d)
		return d
	}

	// get the truncated quotient and remainder
	quo, rem := d, big.NewInt(0)
	quo, rem = quo.QuoRem(d, precisionReuseV1, rem)

	if rem.Sign() == 0 { // remainder is zero
		return quo
	}

	switch rem.Cmp(fivePrecisionV1) {
	case -1:
		return quo
	case 1:
		return quo.Add(quo, oneInt)
	default: // bankers rounding must take place
		// always round to an even number
		if quo.Bit(0) == 0 {
			return quo
		}
		return quo.Add(quo, oneInt)
	}
}

func chopprecisionAndRoundUp(d *big.Int) *big.Int {

	// remove the negative and add it back when returning
	if d.Sign() == -1 {
		// make d positive, compute chopped value, and then un-mutate d
		d = d.Neg(d)
		// truncate since d is negative...
		d = chopprecisionAndTruncate(d)
		d = d.Neg(d)
		return d
	}

	// get the truncated quotient and remainder
	quo, rem := d, big.NewInt(0)
	quo, rem = quo.QuoRem(d, precisionReuseV1, rem)

	if rem.Sign() == 0 { // remainder is zero
		return quo
	}

	return quo.Add(quo, oneInt)
}

func chopprecisionAndRoundNonMutative(d *big.Int) *big.Int {
	tmp := new(big.Int).Set(d)
	return chopprecisionAndRound(tmp)
}

// RoundInt64 rounds the decimal using bankers rounding
func (d Decimal) RoundInt64() int64 {
	chopped := chopprecisionAndRoundNonMutative(d.i)
	if !chopped.IsInt64() {
		panic("Int64() out of bound")
	}
	return chopped.Int64()
}

// RoundInt round the decimal using bankers rounding
func (d Decimal) RoundInt() Int {
	return NewIntFromBigInt(chopprecisionAndRoundNonMutative(d.i))
}

//___________________________________________________________________________________

// similar to chopprecisionAndRound, but always rounds down
func chopprecisionAndTruncate(d *big.Int) *big.Int {
	return d.Quo(d, precisionReuseV1)
}

func chopprecisionAndTruncateNonMutative(d *big.Int) *big.Int {
	tmp := new(big.Int).Set(d)
	return chopprecisionAndTruncate(tmp)
}

// TruncateInt64 truncates the decimals from the number and returns an int64
func (d Decimal) TruncateInt64() int64 {
	chopped := chopprecisionAndTruncateNonMutative(d.i)
	if !chopped.IsInt64() {
		panic("Int64() out of bound")
	}
	return chopped.Int64()
}

// TruncateInt truncates the decimals from the number and returns an Int
func (d Decimal) TruncateInt() Int {
	return NewIntFromBigInt(chopprecisionAndTruncateNonMutative(d.i))
}

// TruncateDec truncates the decimals from the number and returns a Decimal
func (d Decimal) TruncateDec() Decimal {
	return NewDecimalFromBigInt(chopprecisionAndTruncateNonMutative(d.i))
}

// Ceil returns the smallest interger value (as a decimal) that is greater than
// or equal to the given decimal.
func (d Decimal) Ceil() Decimal {
	tmp := new(big.Int).Set(d.i)

	quo, rem := tmp, big.NewInt(0)
	quo, rem = quo.QuoRem(tmp, precisionReuseV1, rem)

	// no need to round with a zero remainder regardless of sign
	if rem.Cmp(zeroInt) == 0 {
		return NewDecimalFromBigInt(quo)
	}

	if rem.Sign() == -1 {
		return NewDecimalFromBigInt(quo)
	}

	return NewDecimalFromBigInt(quo.Add(quo, oneInt))
}

//___________________________________________________________________________________

func init() {
	empty := new(big.Int)
	bz, _ := empty.MarshalText()
	nilJSON, _ = json.Marshal(string(bz))
}

// MarshalJSON marshals the decimal
func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.i == nil {
		return nilJSON, nil
	}
	return json.Marshal(d.String())
}

// UnmarshalJSON defines custom decoding scheme
func (d *Decimal) UnmarshalJSON(bz []byte) error {
	if d.i == nil {
		d.i = new(big.Int)
	}

	var text string
	err := json.Unmarshal(bz, &text)
	if err != nil {
		return err
	}

	// TODO: Reuse dec allocation
	newDec, err := NewDecimalFromStr(text)
	if err != nil {
		return err
	}

	d.i = newDec.i
	return nil
}

// MarshalYAML returns the YAML representation.
func (d Decimal) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

// Marshal implements the gogo proto custom type interface.
func (d Decimal) Marshal() ([]byte, error) {
	if d.i == nil {
		d.i = new(big.Int)
	}
	return d.i.MarshalText()
}

// MarshalTo implements the gogo proto custom type interface.
func (d *Decimal) MarshalTo(data []byte) (n int, err error) {
	if d.i == nil {
		d.i = new(big.Int)
	}
	if len(d.i.Bytes()) == 0 {
		copy(data, []byte{0x30})
		return 1, nil
	}

	bz, err := d.Marshal()
	if err != nil {
		return 0, err
	}

	copy(data, bz)
	return len(bz), nil
}

// Unmarshal implements the gogo proto custom type interface.
func (d *Decimal) Unmarshal(data []byte) error {
	if len(data) == 0 {
		d = nil
		return nil
	}

	if d.i == nil {
		d.i = new(big.Int)
	}

	if err := d.i.UnmarshalText(data); err != nil {
		return err
	}

	if d.i.BitLen() > maxBitLen {
		return fmt.Errorf("decimal out of range; got: %d, max: %d", d.i.BitLen(), maxBitLen)
	}

	return nil
}

// Size implements the gogo proto custom type interface.
func (d *Decimal) Size() int {
	bz, _ := d.Marshal()
	return len(bz)
}

// Override Amino binary serialization by proxying to protobuf.
func (d Decimal) MarshalAmino() ([]byte, error)   { return d.Marshal() }
func (d *Decimal) UnmarshalAmino(bz []byte) error { return d.Unmarshal(bz) }

//___________________________________________________________________________________
// helpers

// test if two decimal arrays are equal
func DecimalEqual(d1s, d2s []Decimal) bool {
	if len(d1s) != len(d2s) {
		return false
	}

	for i, d1 := range d1s {
		if !d1.Equal(d2s[i]) {
			return false
		}
	}
	return true
}

// minimum decimal between two
func MinDecimal(d1, d2 Decimal) Decimal {
	if d1.LT(d2) {
		return d1
	}
	return d2
}

// maximum decimal between two
func MaxDecimal(d1, d2 Decimal) Decimal {
	if d1.LT(d2) {
		return d2
	}
	return d1
}

// intended to be used with require/assert:  require.True(DecEq(...))
func DecimalEq(t *testing.T, exp, got Decimal) (*testing.T, bool, string, string, string) {
	return t, exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
