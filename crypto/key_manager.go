package crypto

import (
	"encoding/hex"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types"
	"strings"

	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultBIP39Passphrase = ""
)

type KeyManager interface {
	Sign(data []byte) ([]byte, error)
	GetPrivKey() crypto.PrivKey
	ExportAsMnemonic() (string, error)
	ExportAsPrivateKey() (string, error)
	ExportAsKeyStore(password string) (types.Keystore, error)
}

type keyManager struct {
	privKey  crypto.PrivKey
	mnemonic string
}

func NewMnemonicKeyManager(mnemonic string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromMnemonic(mnemonic, FullPath)
	return &k, err
}

func NewPrivateKeyManager(priKey string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromPrivateKey(priKey)
	return &k, err
}

func NewKeyManager() (KeyManager, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return NewMnemonicKeyManager(mnemonic)
}

func (m *keyManager) ExportAsMnemonic() (string, error) {
	if m.mnemonic == "" {
		return "", fmt.Errorf("This key manager is not recover from mnemonic or anto generated ")
	}
	return m.mnemonic, nil
}

func (m *keyManager) ExportAsPrivateKey() (string, error) {
	secpPrivateKey, ok := m.privKey.(secp256k1.PrivKeySecp256k1)
	if !ok {
		return "", fmt.Errorf(" Only PrivKeySecp256k1 key is supported ")
	}
	return hex.EncodeToString(secpPrivateKey[:]), nil
}

func (m *keyManager) Sign(data []byte) ([]byte, error) {
	return m.privKey.Sign(data)
}

func (m *keyManager) GetPrivKey() crypto.PrivKey {
	return m.privKey
}

func (m *keyManager) recoveryFromMnemonic(mnemonic, keyPath string) error {
	words := strings.Split(mnemonic, " ")
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("mnemonic length should either be 12 or 24")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, defaultBIP39Passphrase)
	if err != nil {
		return err
	}
	// create master key and derive first key:
	masterPriv, ch := ComputeMastersFromSeed(seed)
	derivedPriv, err := DerivePrivateKeyForPath(masterPriv, ch, keyPath)
	if err != nil {
		return err
	}
	priKey := secp256k1.PrivKeySecp256k1(derivedPriv)
	if err != nil {
		return err
	}
	m.privKey = priKey
	m.mnemonic = mnemonic
	return nil
}

func (m *keyManager) recoveryFromPrivateKey(privateKey string) error {
	priBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}

	if len(priBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], priBytes[:32])
	priKey := secp256k1.PrivKeySecp256k1(keyBytesArray)
	m.privKey = priKey
	return nil
}
