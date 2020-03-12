package adapter

import (
	"encoding/json"
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
	store, err := adapter.keyDAO.Read(name, password)
	if err != nil {
		return signature, err
	}

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
	store, err := adapter.keyDAO.Read(name, password)
	if err != nil {
		return addr, err
	}

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
	store, err := adapter.keyDAO.Read(name, password)
	if err != nil || store != nil {
		return "", errors.New(fmt.Sprintf("%s has existed", name))
	}

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

func (adapter daoAdapter) Import(name, password string, keystore string) (address string, err error) {
	store, err := adapter.keyDAO.Read(name, password)
	if err != nil || store != nil {
		return "", errors.New(fmt.Sprintf("%s has existed", name))
	}

	km, err := crypto.NewKeyStoreKeyManager(keystore, password)
	if err != nil {
		return "", err
	}

	address = types.AccAddress(km.GetPrivKey().PubKey().Address()).String()
	err = adapter.keyDAO.Write(name, keystore)
	return
}

func (adapter daoAdapter) Export(name, password, encryptKeystorePwd string) (keystore string, err error) {
	store, err := adapter.keyDAO.Read(name, password)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s not existed", name))
	}
	var km crypto.KeyManager
	switch store := store.(type) {
	case types.KeyInfo:
		km, err = crypto.NewPrivateKeyManager(store.PrivKey)
		if err != nil {
			return "", err
		}
	case types.KeystoreInfo:
		km, err = crypto.NewKeyStoreKeyManager(store.KeystoreJSON, password)
		if err != nil {
			return "", err
		}
	}
	keyStore, err := km.ExportAsKeystore(encryptKeystorePwd)
	if err != nil {
		return "", err
	}

	keyStore.Address = types.AccAddress(km.GetPrivKey().PubKey().Address()).String()
	bz, err := json.Marshal(keyStore)
	if err != nil {
		return "", err
	}
	return string(bz), nil
}

func (adapter daoAdapter) Delete(name, password string) error {
	return adapter.keyDAO.Delete(name, password)
}

func (adapter daoAdapter) Query(name string) (address string, err error) {
	store, err := adapter.keyDAO.Read(name, "")
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s not existed", name))
	}
	switch store := store.(type) {
	case types.KeyInfo:
		address = store.Address
	case types.KeystoreInfo:
		var keystore crypto.Keystore
		err := json.Unmarshal([]byte(store.KeystoreJSON), &keystore)
		if err != nil {
			return "", err
		}
		address = keystore.Address
	}
	return address, nil
}

func adapt(km crypto.KeyManager, storeType types.StoreType, password string) (address string, store types.Store, err error) {
	address = types.AccAddress(km.GetPrivKey().PubKey().Address()).String()
	switch storeType {
	case types.Keystore:
		keystore, err := km.ExportAsKeystore(password)
		if err != nil {
			return "", "", err
		}

		keystore.Address = address
		bz, err := json.Marshal(keystore)
		if err != nil {
			return "", "", err
		}

		store = types.KeystoreInfo{
			KeystoreJSON: string(bz),
		}
		return address, store, nil
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
