package client

import (
	"fmt"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	kmg "github.com/irisnet/core-sdk-go/common/crypto"
	cryptoamino "github.com/irisnet/core-sdk-go/common/crypto/codec"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/secp256k1"
	"github.com/irisnet/core-sdk-go/common/crypto/keys/sm2"
	commoncryptotypes "github.com/irisnet/core-sdk-go/common/crypto/types"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
)

type KeyManager struct {
	KeyDAO store.KeyDAO
	Algo   string
}

func (k KeyManager) Add(name, password string) (string, string, types.Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, types.Wrap(err)
}

func (k KeyManager) Sign(name, password string, data []byte) ([]byte, tmcrypto.PubKey, error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	km, err := kmg.NewPrivateKeyManager([]byte(info.PrivKeyArmor), string(info.Algo))
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	signByte, err := km.Sign(data)
	if err != nil {
		return nil, nil, err
	}

	return signByte, FromTmPubKey(info.Algo, km.ExportPubKey()), nil
}

func (k KeyManager) Insert(name, password string) (string, string, error) {
	if k.KeyDAO.Has(name) {
		return "", "", fmt.Errorf("name %s has existed", name)
	}

	km, err := kmg.NewAlgoKeyManager(k.Algo)
	if err != nil {
		return "", "", err
	}

	mnemonic, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       cryptoamino.MarshalPubkey(pubKey),
		PrivKeyArmor: string(cryptoamino.MarshalPrivKey(priv)),
		Algo:         k.Algo,
	}

	if err = k.KeyDAO.Write(name, password, info); err != nil {
		return "", "", err
	}
	return address, mnemonic, nil
}

func (k KeyManager) Recover(name, password, mnemonic, hdPath string) (string, error) {
	if k.KeyDAO.Has(name) {
		return "", fmt.Errorf("name %s has existed", name)
	}
	var (
		km  kmg.KeyManager
		err error
	)
	if hdPath == "" {
		km, err = kmg.NewMnemonicKeyManager(mnemonic, k.Algo)
	} else {
		km, err = kmg.NewMnemonicKeyManagerWithHDPath(mnemonic, k.Algo, hdPath)
	}

	if err != nil {
		return "", err
	}

	_, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       cryptoamino.MarshalPubkey(pubKey),
		PrivKeyArmor: string(cryptoamino.MarshalPrivKey(priv)),
		Algo:         k.Algo,
	}

	if err = k.KeyDAO.Write(name, password, info); err != nil {
		return "", err
	}

	return address, nil
}

func (k KeyManager) Import(name, password, armor string) (string, error) {
	if k.KeyDAO.Has(name) {
		return "", fmt.Errorf("%s has existed", name)
	}

	km := kmg.NewKeyManager()

	priv, _, err := km.ImportPrivKey(armor, password)
	if err != nil {
		return "", err
	}

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       cryptoamino.MarshalPubkey(pubKey),
		PrivKeyArmor: string(cryptoamino.MarshalPrivKey(priv)),
		Algo:         k.Algo,
	}

	err = k.KeyDAO.Write(name, password, info)
	if err != nil {
		return "", err
	}
	return address, nil
}

func (k KeyManager) Export(name, password string) (armor string, err error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return armor, fmt.Errorf("name %s not exist", name)
	}

	km, err := kmg.NewPrivateKeyManager([]byte(info.PrivKeyArmor), info.Algo)
	if err != nil {
		return "", err
	}

	return km.ExportPrivKey(password)
}

func (k KeyManager) Delete(name, password string) error {
	return k.KeyDAO.Delete(name, password)
}

func (k KeyManager) Find(name, password string) (tmcrypto.PubKey, types.AccAddress, error) {
	info, err := k.KeyDAO.Read(name, password)
	if err != nil {
		return nil, nil, types.WrapWithMessage(err, "name %s not exist", name)
	}

	pubKey, err := cryptoamino.PubKeyFromBytes(info.PubKey)
	if err != nil {
		return nil, nil, types.WrapWithMessage(err, "name %s not exist", name)
	}

	return FromTmPubKey(info.Algo, pubKey), types.AccAddress(pubKey.Address().Bytes()), nil
}

func FromTmPubKey(Algo string, pubKey tmcrypto.PubKey) commoncryptotypes.PubKey {
	var pubkey commoncryptotypes.PubKey
	pubkeyBytes := pubKey.Bytes()
	switch Algo {
	case "sm2":
		pubkey = &sm2.PubKey{Key: pubkeyBytes}
	case "secp256k1":
		pubkey = &secp256k1.PubKey{Key: pubkeyBytes}
	}
	return pubkey
}

type Client interface {
	Add(name, password string) (address string, mnemonic string, err types.Error)
	Recover(name, password, mnemonic string) (address string, err types.Error)
	RecoverWithHDPath(name, password, mnemonic, hdPath string) (address string, err types.Error)
	Import(name, password, privKeyArmor string) (address string, err types.Error)
	Export(name, password string) (privKeyArmor string, err types.Error)
	Delete(name, password string) types.Error
	Show(name, password string) (string, types.Error)
}

type keysClient struct {
	types.KeyManager
}

func NewClient(keyManager types.KeyManager) Client {
	return keysClient{keyManager}
}

func (k keysClient) Add(name, password string) (string, string, types.Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, types.Wrap(err)
}

func (k keysClient) Recover(name, password, mnemonic string) (string, types.Error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic, "")
	return address, types.Wrap(err)
}

func (k keysClient) RecoverWithHDPath(name, password, mnemonic, hdPath string) (string, types.Error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic, hdPath)
	return address, types.Wrap(err)
}

func (k keysClient) Import(name, password, privKeyArmor string) (string, types.Error) {
	address, err := k.KeyManager.Import(name, password, privKeyArmor)
	return address, types.Wrap(err)
}

func (k keysClient) Export(name, password string) (string, types.Error) {
	keystore, err := k.KeyManager.Export(name, password)
	return keystore, types.Wrap(err)
}

func (k keysClient) Delete(name, password string) types.Error {
	err := k.KeyManager.Delete(name, password)
	return types.Wrap(err)
}

func (k keysClient) Show(name, password string) (string, types.Error) {
	_, address, err := k.KeyManager.Find(name, password)
	if err != nil {
		return "", types.Wrap(err)
	}
	return address.String(), nil
}
