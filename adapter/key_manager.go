package adapter

import (
	"errors"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/crypto"
	"github.com/irisnet/irishub-sdk-go/types"
)

type daoAdapter struct {
	keyDAO    types.KeyDAO
	storeType types.StoreType
}

func NewDAOAdapter(dao types.KeyDAO, storeType types.StoreType) types.KeyManager {
	return daoAdapter{
		keyDAO:    dao,
		storeType: storeType,
	}
}

func (adapter daoAdapter) Sign(name, password string, data []byte) (signature types.Signature, err error) {
	store := adapter.keyDAO.Read(name)

	var mm crypto.KeyManager
	switch store := store.(type) {
	case types.KeyInfo:
		mm, err = crypto.NewPrivateKeyManager(store.PrivKey)
		if err != nil {
			return signature, err
		}
	case types.KeystoreInfo:
		mm, err = crypto.NewKeyStoreKeyManager(store.KeystoreJSON, password)
		if err != nil {
			return signature, err
		}
	}
	signByte, err := mm.Sign(data)

	return types.Signature{
		PubKey:    mm.GetPrivKey().PubKey(),
		Signature: signByte,
	}, nil
}

func (adapter daoAdapter) QueryAddress(name, password string) (addr types.AccAddress, err error) {
	store := adapter.keyDAO.Read(name)

	var mm crypto.KeyManager
	switch store := store.(type) {
	case types.KeyInfo:
		mm, err = crypto.NewPrivateKeyManager(store.PrivKey)
		if err != nil {
			return addr, err
		}
		return types.AccAddressFromBech32(store.Address)
	case types.KeystoreInfo:
		mm, err = crypto.NewKeyStoreKeyManager(store.KeystoreJSON, password)
		if err != nil {
			return addr, err
		}
		accAddr := types.AccAddress(mm.GetPrivKey().PubKey().Address())
		return accAddr, nil
	}
	return addr, errors.New("invalid StoreType")
}

func (adapter daoAdapter) Insert(name, password string) (string, string, error) {
	km, err := crypto.NewKeyManager()
	if err != nil {
		return "", "", err
	}
	address, store, err := adapt(km, adapter.storeType, password)
	if err != nil {
		return "", "", err
	}

	mnemonic, err := km.ExportAsMnemonic()
	if err != nil {
		return "", "", err
	}

	err = adapter.keyDAO.Write(name, store)
	return address, mnemonic, err
}

func (adapter daoAdapter) Recover(name, password, mnemonic string) (string, error) {
	km, err := crypto.NewMnemonicKeyManager(mnemonic)
	if err != nil {
		return "", err
	}

	address, store, err := adapt(km, adapter.storeType, password)
	if err != nil {
		return address, err
	}

	err = adapter.keyDAO.Write(name, store)
	return address, err
}

func adapt(km crypto.KeyManager, storeType types.StoreType, password string) (address string, store types.Store, err error) {
	address = types.AccAddress(km.GetPrivKey().PubKey().Address()).String()
	switch storeType {
	case types.Keystore:
		store, err = km.ExportAsKeystore(password)
		return
	case types.Key:
		privKey, err := km.ExportAsPrivateKey()
		if err != nil {
			return address, store, err
		}
		store = types.KeyInfo{
			PrivKey:  privKey,
			Address:  address,
			Password: password,
		}
	}
	return address, store, fmt.Errorf("invalid storeType:%d", storeType)
}
