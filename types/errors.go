package types

import (
	"fmt"

	"github.com/pkg/errors"
)

// RootCodespace is the codespace for all errors defined in this package
const RootCodespace = "sdk"
const GoSDKCodespace = "go-sdk"

// undefinedCodespace when we explicitly declare no codespace
const undefinedCodespace = "undefined"

var (
	ErrUnkonwn = register(undefinedCodespace, 111222, "unknown error")
	Nil        = Error{}
)

func init() {
	_ = register(undefinedCodespace, 1, "internal")
	_ = register(RootCodespace, 2, "tx parse error")
	_ = register(RootCodespace, 3, "invalid sequence")
	_ = register(RootCodespace, 4, "unauthorized")
	_ = register(RootCodespace, 5, "insufficient funds")
	_ = register(RootCodespace, 6, "unknown request")
	_ = register(RootCodespace, 7, "invalid address")
	_ = register(RootCodespace, 8, "invalid pubkey")
	_ = register(RootCodespace, 9, "unknown address")
	_ = register(RootCodespace, 10, "invalid coins")
	_ = register(RootCodespace, 11, "out of gas")
	_ = register(RootCodespace, 12, "memo too large")
	_ = register(RootCodespace, 13, "insufficient fee")
	_ = register(RootCodespace, 14, "maximum number of signatures exceeded")
	_ = register(RootCodespace, 15, "no signatures supplied")
	_ = register(RootCodespace, 16, "failed to marshal JSON bytes")
	_ = register(RootCodespace, 17, "failed to unmarshal JSON bytes")
	_ = register(RootCodespace, 18, "invalid request")
	_ = register(RootCodespace, 19, "tx already in mempool")
	_ = register(RootCodespace, 20, "mempool is full")
	_ = register(RootCodespace, 21, "tx too large")
}

// register returns an error instance that should be used as the base for
// creating error instances during runtime.
//
// Popular root errors are declared in this package, but extensions may want to
// declare custom codes. This function ensures that no error code is used
// twice. Attempt to reuse an error code results in panic.
//
// Use this function only during a program startup phase.
func register(codespace string, code uint32, description string) Error {
	err := New(codespace, code, description)
	setUsed(err)

	return err
}

// usedCodes is keeping track of used codes to ensure their uniqueness. No two
// error instances should share the same (codespace, code) tuple.
var usedCodes = map[string]Error{}

func errorID(codespace string, code uint32) string {
	return fmt.Sprintf("%s:%d", codespace, code)
}

func setUsed(err Error) {
	usedCodes[errorID(err.codespace, err.code)] = err
}

func GetError(codespace string, code uint32, log ...string) Error {
	err, ok := usedCodes[errorID(codespace, code)]
	if ok {
		return err
	}
	if len(log) == 0 {
		return ErrUnkonwn
	}
	return Error{
		codespace: codespace,
		code:      code,
		desc:      log[0],
	}
}

// Error represents a root error.
//
// Weave framework is using root error to categorize issues. Each instance
// created during the runtime should wrap one of the declared root errors. This
// allows error tests and returning all errors to the client in a safe manner.
//
// All popular root errors are declared in this package. If an extension has to
// declare a custom root error, always use register function to ensure
// error code uniqueness.
type Error struct {
	codespace string
	code      uint32
	desc      string
}

func New(codespace string, code uint32, desc string) Error {
	return Error{
		codespace: codespace,
		code:      code,
		desc:      desc,
	}
}

func (e Error) Error() string {
	return e.desc
}

func (e Error) Code() uint32 {
	return e.code
}

func (e Error) Codespace() string {
	return e.codespace
}

func (e Error) IsNil() bool {
	return e.codespace == Nil.codespace &&
		e.code == Nil.code &&
		e.desc == Nil.desc
}

// Wrap extends given error with an additional information.
//
// If the wrapped error does not provide ABCICode method (ie. stdlib errors),
// it will be labeled as internal error.
//
// If err is nil, this returns nil, avoiding the need for an if statement when
// wrapping a error returned at the end of a function
func Wrap(err error) Error {
	if err == nil {
		return Nil
	}

	return Error{
		codespace: GoSDKCodespace,
		code:      ErrUnkonwn.code,
		desc:      err.Error(),
	}
}

// Wrapf extends given error with an additional information.
//
// This function works like Wrap function with additional functionality of
// formatting the input as specified.
func Wrapf(format string, args ...interface{}) Error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(errors.New(desc))
}
