package adapter

import (
	"errors"

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

func (adapter daoAdapter) Insert(name, password string) (address, mnemonic string, err error) {
	km, err := crypto.NewKeyManager()
	if err != nil {
		return address, mnemonic, err
	}
	var store types.Store

	address = types.AccAddress(km.GetPrivKey().PubKey().Address()).String()
	switch adapter.storeType {
	case types.Keystore:
		store, err = km.ExportAsKeystore(password)
		mnemonic, err = km.ExportAsMnemonic()
		return
	case types.Key:
		privKey, err := km.ExportAsPrivateKey()
		if err != nil {
			return address, mnemonic, err
		}
		store = types.KeyInfo{
			PrivKey:  privKey,
			Address:  address,
			Password: password,
		}
	}
	err = adapter.keyDAO.Write(name, store)
	return
}

func (adapter daoAdapter) Recover(name, password, mnemonic string) (address string, err error) {
	km, err := crypto.NewMnemonicKeyManager(mnemonic)
	if err != nil {
		return address, err
	}
	var store types.Store

	address = types.AccAddress(km.GetPrivKey().PubKey().Address()).String()
	switch adapter.storeType {
	case types.Keystore:
		store, err = km.ExportAsKeystore(password)
		mnemonic, err = km.ExportAsMnemonic()
		return
	case types.Key:
		privKey, err := km.ExportAsPrivateKey()
		if err != nil {
			return address, err
		}
		store = types.KeyInfo{
			PrivKey:  privKey,
			Address:  address,
			Password: password,
		}
	}
	err = adapter.keyDAO.Write(name, store)
	return
}
