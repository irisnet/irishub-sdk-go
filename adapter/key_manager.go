// Package adapter is to adapt to the user DAO layer, the user can not override this implementation
//
//
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

//NewDAOAdapter return a apapter for user DAO
func NewDAOAdapter(dao types.KeyDAO, storeType types.StoreType) types.KeyManager {
	return daoAdapter{
		keyDAO:    dao,
		storeType: storeType,
	}
}

func (adapter daoAdapter) Sign(name, password string, data []byte) (signature types.Signature, err error) {
	store, err := adapter.keyDAO.Read(name)
	if err != nil {
		return signature, err
	}

	var mm crypto.KeyManager
	switch store := store.(type) {
	case types.PrivKeyInfo:
		privKey, err := adapter.keyDAO.Decrypt(store.PrivKey, password)
		if err != nil {
			return signature, err
		}
		mm, err = crypto.NewPrivateKeyManager(privKey)
		if err != nil {
			return signature, err
		}
	case types.KeystoreInfo:
		mm, err = crypto.NewKeyStoreKeyManager(store.Keystore, password)
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

func (adapter daoAdapter) Insert(name, password string) (string, string, error) {
	km, err := crypto.NewKeyManager()
	if err != nil {
		return "", "", err
	}
	address, store, err := adapter.apply(km, password)
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
	store, err := adapter.keyDAO.Read(name)
	if err != nil || store != nil {
		return "", errors.New(fmt.Sprintf("%s has existed", name))
	}

	km, err := crypto.NewMnemonicKeyManager(mnemonic)
	if err != nil {
		return "", err
	}

	address, store, err := adapter.apply(km, password)
	if err != nil {
		return address, err
	}

	err = adapter.keyDAO.Write(name, store)
	return address, err
}

func (adapter daoAdapter) Import(name, password string, keystore string) (string, error) {
	//store, err := adapter.keyDAO.Read(name)
	//if err != nil || store != nil {
	//	return "", fmt.Errorf("%s has existed", name)
	//}

	km, err := crypto.NewKeyStoreKeyManager(keystore, password)
	if err != nil {
		return "", err
	}
	address, s, err := adapter.apply(km, password)
	if err != nil {
		return "", err
	}
	return address, adapter.keyDAO.Write(name, s)
}

func (adapter daoAdapter) Export(name, password, encryptKeystorePwd string) (keystore string, err error) {
	store, err := adapter.keyDAO.Read(name)
	if err != nil {
		return "", fmt.Errorf("%s not existed", name)
	}
	var km crypto.KeyManager
	switch store := store.(type) {
	case types.PrivKeyInfo:
		privKey, err := adapter.keyDAO.Decrypt(store.PrivKey, password)
		if err != nil {
			return "", err
		}
		km, err = crypto.NewPrivateKeyManager(privKey)
		if err != nil {
			return "", err
		}
	case types.KeystoreInfo:
		km, err = crypto.NewKeyStoreKeyManager(store.Keystore, password)
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

func (adapter daoAdapter) Delete(name string) error {
	return adapter.keyDAO.Delete(name)
}

func (adapter daoAdapter) Query(name string) (types.AccAddress, error) {
	store, err := adapter.keyDAO.Read(name)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s not existed", name))
	}
	switch store := store.(type) {
	case types.PrivKeyInfo:
		return types.AccAddressFromBech32(store.Address)
	case types.KeystoreInfo:
		var keystore crypto.Keystore
		err := json.Unmarshal([]byte(store.Keystore), &keystore)
		if err != nil {
			return nil, err
		}
		return types.AccAddressFromBech32(keystore.Address)
	}
	return nil, errors.New("invalid Store")
}

func (adapter daoAdapter) apply(km crypto.KeyManager, password string) (address string, store types.Store, err error) {
	address = types.AccAddress(km.GetPrivKey().PubKey().Address()).String()
	switch adapter.storeType {
	case types.Keystore:
		keystore, err := km.ExportAsKeystore(password)
		if err != nil {
			return "", store, err
		}

		keystore.Address = address
		bz, err := json.Marshal(keystore)
		if err != nil {
			return "", store, err
		}

		store = types.KeystoreInfo{
			Keystore: string(bz),
		}
		return address, store, nil
	case types.PrivKey:
		privKey, err := km.ExportAsPrivateKey()
		if err != nil {
			return address, store, err
		}
		pk, err := adapter.keyDAO.Encrypt(privKey, password)
		if err != nil {
			return "", nil, err
		}
		store = types.PrivKeyInfo{
			PrivKey: pk,
			Address: address,
		}
		return address, store, nil
	}
	return address, store, fmt.Errorf("invalid storeType:%d", adapter.storeType)
}
