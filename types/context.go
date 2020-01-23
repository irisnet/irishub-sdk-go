package types

import (
	"fmt"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/pkg/errors"

	"github.com/irisnet/irishub-sdk-go/crypto"
	"github.com/irisnet/irishub-sdk-go/net"
)

type AbstractClient interface {
	Broadcast(baseTx BaseTx, msg []Msg) (Result, error)
	Query(path string, data interface{}, result interface{}) error
	QueryStore(key cmn.HexBytes, storeName string) ([]byte, error)
	QueryAccount(address string) (BaseAccount, error)
	GetSender(name string) AccAddress
	GetRPC() net.RPCClient
	GetCodec() Codec
}

// TxContext implements a transaction context created in SDK modules.
type TxContext struct {
	Codec         Codec
	AccountNumber uint64
	Sequence      uint64
	Gas           uint64
	ChainID       string
	Memo          string
	Fee           string
	KeyDAO        KeyDAO
	Online        bool
	Network       Network
	Mode          BroadcastMode
	RPC           net.RPCClient
	Simulate      bool
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
func (txCtx *TxContext) WithFee(fee string) *TxContext {
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
func (txCtx *TxContext) WithKeyDAO(keyDao KeyDAO) *TxContext {
	txCtx.KeyDAO = keyDao
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

// WithRPC returns a pointer of the context with a rpc.
func (txCtx *TxContext) WithRPC(rpc net.RPCClient) *TxContext {
	txCtx.RPC = rpc
	return txCtx
}

// WithRPC returns a pointer of the context with a simulate.
func (txCtx *TxContext) WithSimulate(simulate bool) *TxContext {
	txCtx.Simulate = simulate
	return txCtx
}

func (txCtx TxContext) BuildAndSign(name string, msgs []Msg) ([]byte, error) {
	msg, err := txCtx.Build(msgs)
	if err != nil {
		return nil, err
	}
	return txCtx.Sign(name, msg)
}

// Build builds a single message to be signed from a TxContext given a set of
// messages. It returns an error if a Fee is supplied but cannot be parsed.
func (txCtx TxContext) Build(msgs []Msg) (StdSignMsg, error) {
	chainID := txCtx.ChainID
	if chainID == "" {
		return StdSignMsg{}, errors.Errorf("chain ID required but not specified")
	}

	fee := Coins{}
	if txCtx.Fee != "" {
		parsedFee, err := ParseCoins(txCtx.Fee)
		if err != nil {
			return StdSignMsg{}, fmt.Errorf("encountered error in parsing transaction Fee: %s", err.Error())
		}

		fee = parsedFee
	}

	return StdSignMsg{
		ChainID:       txCtx.ChainID,
		AccountNumber: txCtx.AccountNumber,
		Sequence:      txCtx.Sequence,
		Memo:          txCtx.Memo,
		Msgs:          msgs,
		Fee:           NewStdFee(txCtx.Gas, fee...),
	}, nil
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (txCtx TxContext) Sign(name string, msg StdSignMsg) ([]byte, error) {
	sig, err := txCtx.makeSignature(name, msg)
	if err != nil {
		return nil, err
	}
	tx := NewStdTx(msg.Msgs, msg.Fee, []StdSignature{sig}, msg.Memo)
	return txCtx.Codec.MarshalBinaryLengthPrefixed(tx)
}

func (txCtx TxContext) makeSignature(name string, msg StdSignMsg) (sig StdSignature, err error) {
	sig = StdSignature{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
	}
	if !txCtx.Simulate {
		keystore := txCtx.KeyDAO.Read(name)
		keyManager, err := crypto.NewPrivateKeyManager(keystore.GetPrivate())
		if err != nil {
			return sig, err
		}
		sigBytes, err := keyManager.Sign(msg.Bytes(txCtx.Codec))
		if err != nil {
			return sig, err
		}
		sig.PubKey = keyManager.GetPrivKey().PubKey()
		sig.Signature = sigBytes
	}
	return sig, nil
}
