package types

import (
	"github.com/pkg/errors"
)

// TxContext implements a transaction context created in SDK modules.
type TxContext struct {
	codec         Codec
	accountNumber uint64
	sequence      uint64
	gas           uint64
	chainID       string
	memo          string
	fee           Coins
	network       Network
	mode          BroadcastMode
	simulate      bool
	password      string
	keyManager    KeyManager
}

// WithCodec returns a pointer of the context with an updated codec.
func (txCtx *TxContext) WithCodec(cdc Codec) *TxContext {
	txCtx.codec = cdc
	return txCtx
}

// WithCodec returns a pointer of the context with an updated codec.
func (txCtx *TxContext) Codec() Codec {
	return txCtx.codec
}

// WithChainID returns a pointer of the context with an updated ChainID.
func (txCtx *TxContext) WithChainID(chainID string) *TxContext {
	txCtx.chainID = chainID
	return txCtx
}

// WithChainID returns a pointer of the context with an updated ChainID.
func (txCtx *TxContext) ChainID() string {
	return txCtx.chainID
}

// WithGas returns a pointer of the context with an updated Gas.
func (txCtx *TxContext) WithGas(gas uint64) *TxContext {
	txCtx.gas = gas
	return txCtx
}

// WithGas returns a pointer of the context with an updated Gas.
func (txCtx *TxContext) Gas() uint64 {
	return txCtx.gas
}

// WithFee returns a pointer of the context with an updated Fee.
func (txCtx *TxContext) WithFee(fee Coins) *TxContext {
	txCtx.fee = fee
	return txCtx
}

// WithFee returns a pointer of the context with an updated Fee.
func (txCtx *TxContext) Fee() Coins {
	return txCtx.fee
}

// WithSequence returns a pointer of the context with an updated sequence number.
func (txCtx *TxContext) WithSequence(sequence uint64) *TxContext {
	txCtx.sequence = sequence
	return txCtx
}

// WithSequence returns a pointer of the context with an updated sequence number.
func (txCtx *TxContext) Sequence() uint64 {
	return txCtx.sequence
}

// WithMemo returns a pointer of the context with an updated memo.
func (txCtx *TxContext) WithMemo(memo string) *TxContext {
	txCtx.memo = memo
	return txCtx
}

// WithMemo returns a pointer of the context with an updated memo.
func (txCtx *TxContext) Memo() string {
	return txCtx.memo
}

// WithAccountNumber returns a pointer of the context with an account number.
func (txCtx *TxContext) WithAccountNumber(accnum uint64) *TxContext {
	txCtx.accountNumber = accnum
	return txCtx
}

// WithAccountNumber returns a pointer of the context with an account number.
func (txCtx *TxContext) AccountNumber() uint64 {
	return txCtx.accountNumber
}

// WithAccountNumber returns a pointer of the context with a keyDao.
func (txCtx *TxContext) WithKeyManager(keyManager KeyManager) *TxContext {
	txCtx.keyManager = keyManager
	return txCtx
}

// WithAccountNumber returns a pointer of the context with a keyDao.
func (txCtx *TxContext) KeyManager() KeyManager {
	return txCtx.keyManager
}

// WithNetwork returns a pointer of the context with a Network.
func (txCtx *TxContext) WithNetwork(network Network) *TxContext {
	txCtx.network = network
	return txCtx
}

// WithNetwork returns a pointer of the context with a Network.
func (txCtx *TxContext) Network() Network {
	return txCtx.network
}

// WithMode returns a pointer of the context with a Mode.
func (txCtx *TxContext) WithMode(mode BroadcastMode) *TxContext {
	txCtx.mode = mode
	return txCtx
}

// WithMode returns a pointer of the context with a Mode.
func (txCtx *TxContext) Mode() BroadcastMode {
	return txCtx.mode
}

// WithRPC returns a pointer of the context with a simulate.
func (txCtx *TxContext) WithSimulate(simulate bool) *TxContext {
	txCtx.simulate = simulate
	return txCtx
}

// WithRPC returns a pointer of the context with a simulate.
func (txCtx *TxContext) Simulate() bool {
	return txCtx.simulate
}

// WithRPC returns a pointer of the context with a password.
func (txCtx *TxContext) WithPassword(password string) *TxContext {
	txCtx.password = password
	return txCtx
}

// WithRPC returns a pointer of the context with a password.
func (txCtx *TxContext) Password() string {
	return txCtx.password
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
	if txCtx.chainID == "" {
		return StdSignMsg{}, errors.Errorf("chain ID required but not specified")
	}
	return StdSignMsg{
		ChainID:       txCtx.chainID,
		AccountNumber: txCtx.accountNumber,
		Sequence:      txCtx.sequence,
		Memo:          txCtx.memo,
		Msgs:          msgs,
		Fee:           NewStdFee(txCtx.gas, txCtx.fee...),
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
	if !txCtx.Simulate() {
		signature, err := txCtx.keyManager.Sign(name, txCtx.password, msg.Bytes(txCtx.codec))
		if err != nil {
			return sig, err
		}
		sig.PubKey = signature.PubKey
		sig.Signature = signature.Signature
	}
	return sig, nil
}
