package types

import (
	"github.com/tendermint/tendermint/crypto"

	commoncodec "github.com/irisnet/core-sdk-go/common/codec"
)

type (
	// Generator defines an interface a client can utilize to generate an
	// application-defined concrete transaction type. The type returned must
	// implement ClientTx.
	Generator interface {
		NewTx() ClientTx
		NewFee() ClientFee
		NewSignature() ClientSignature
	}

	ClientFee interface {
		Fee
		SetGas(uint64)
		SetAmount(Coins)
	}

	ClientSignature interface {
		Signature
		SetPubKey(crypto.PubKey) error
		SetSignature([]byte)
	}

	// ClientTx defines an interface which an application-defined concrete transaction
	// type must implement. Namely, it must be able to set messages, generate
	// signatures, and provide canonical bytes to sign over. The transaction must
	// also know how to encode itself.
	ClientTx interface {
		Tx
		commoncodec.ProtoMarshaler

		SetMsgs(...Msg) error
		GetSignatures() []Signature
		SetSignatures(...ClientSignature) error
		GetFee() Fee
		SetFee(ClientFee) error
		GetMemo() string
		SetMemo(string)

		// CanonicalSignBytes returns the canonical JSON bytes to sign over, given a
		// chain ID, along with an account and sequence number. The JSON encoding
		// ensures all field names adhere to their proto definition, default values
		// are omitted, and follows the JSON Canonical Form.
		CanonicalSignBytes(cid string, num, seq uint64) ([]byte, error)
	}
)
