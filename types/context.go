package types

import (
	"github.com/pkg/errors"
)

// TxContext implements a transaction context created in SDK modules.
type TxContext struct {
	Codec         Codec
	AccountNumber uint64
	Sequence      uint64
	Gas           uint64
	ChainID       string
	Memo          string
	Fee           Coins
	Online        bool
	Network       Network
	Mode          BroadcastMode
	Simulate      bool
	Password      string
	KeyManager    KeyManager
}

// WithCodec returns a pointer of the context with an updated codec.
func (txCtx *TxContext) WithCodec(cdc Codec) *TxContext {
	txCtx.Codec = cdc
	return txCtx
}

// WithChainID returns a pointer of the context with an updated ChainID.
func (txCtx *TxContext) WithChainID(chainID string) *TxContext {
	txCtx.ChainID = chainID
	return txCtx
}

// WithGas returns a pointer of the context with an updated Gas.
func (txCtx *TxContext) WithGas(gas uint64) *TxContext {
	txCtx.Gas = gas
	return txCtx
}

// WithFee returns a pointer of the context with an updated Fee.
func (txCtx *TxContext) WithFee(fee Coins) *TxContext {
	txCtx.Fee = fee
	return txCtx
}

// WithSequence returns a pointer of the context with an updated sequence number.
func (txCtx *TxContext) WithSequence(sequence uint64) *TxContext {
	txCtx.Sequence = sequence
	return txCtx
}

// WithMemo returns a pointer of the context with an updated memo.
func (txCtx *TxContext) WithMemo(memo string) *TxContext {
	txCtx.Memo = memo
	return txCtx
}

// WithAccountNumber returns a pointer of the context with an account number.
func (txCtx *TxContext) WithAccountNumber(accnum uint64) *TxContext {
	txCtx.AccountNumber = accnum
	return txCtx
}

// WithAccountNumber returns a pointer of the context with a keyDao.
func (txCtx *TxContext) WithKeyManager(keyManager KeyManager) *TxContext {
	txCtx.KeyManager = keyManager
	return txCtx
}

// WithOnline returns a pointer of the context with an online
func (txCtx *TxContext) WithOnline(online bool) *TxContext {
	txCtx.Online = online
	return txCtx
}

// WithNetwork returns a pointer of the context with a Network.
func (txCtx *TxContext) WithNetwork(network Network) *TxContext {
	txCtx.Network = network
	return txCtx
}

// WithMode returns a pointer of the context with a Mode.
func (txCtx *TxContext) WithMode(mode BroadcastMode) *TxContext {
	txCtx.Mode = mode
	return txCtx
}

// WithRPC returns a pointer of the context with a simulate.
func (txCtx *TxContext) WithSimulate(simulate bool) *TxContext {
	txCtx.Simulate = simulate
	return txCtx
}

// WithRPC returns a pointer of the context with a password.
func (txCtx *TxContext) WithPassword(password string) *TxContext {
	txCtx.Password = password
	return txCtx
}

func (txCtx *TxContext) BuildAndSign(name string, msgs []Msg) (StdTx, error) {
	msg, err := txCtx.Build(msgs)
	if err != nil {
		return StdTx{}, err
	}
	return txCtx.Sign(name, msg)
}

// Build builds a single message to be signed from a TxContext given a set of
// messages. It returns an error if a Fee is supplied but cannot be parsed.
func (txCtx *TxContext) Build(msgs []Msg) (StdSignMsg, error) {
	chainID := txCtx.ChainID
	if chainID == "" {
		return StdSignMsg{}, errors.Errorf("chain ID required but not specified")
	}
	return StdSignMsg{
		ChainID:       txCtx.ChainID,
		AccountNumber: txCtx.AccountNumber,
		Sequence:      txCtx.Sequence,
		Memo:          txCtx.Memo,
		Msgs:          msgs,
		Fee:           NewStdFee(txCtx.Gas, txCtx.Fee...),
	}, nil
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (txCtx *TxContext) Sign(name string, msg StdSignMsg) (StdTx, error) {
	sig, err := txCtx.makeSignature(name, msg)
	if err != nil {
		return StdTx{}, err
	}
	return NewStdTx(msg.Msgs, msg.Fee, []StdSignature{sig}, msg.Memo), nil
}

func (txCtx *TxContext) makeSignature(name string, msg StdSignMsg) (sig StdSignature, err error) {
	sig = StdSignature{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
	}
	if !txCtx.Simulate {
		signature, err := txCtx.KeyManager.Sign(name, "", msg.Bytes(txCtx.Codec))
		if err != nil {
			return sig, err
		}
		sig.PubKey = signature.PubKey
		sig.Signature = signature.Signature
	}
	return sig, nil
}
