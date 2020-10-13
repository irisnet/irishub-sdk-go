package keys

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
)

const ModuleName = "keys"

type keysClient struct {
	original.KeyManager
}

func Create(keyManager original.KeyManager) rpc.Keys {
	return keysClient{
		KeyManager: keyManager,
	}
}

func (k keysClient) Add(name, password string) (string, string, original.Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, original.Wrap(err)
}

func (k keysClient) Recover(name, password, mnemonic string) (string, original.Error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic)
	return address, original.Wrap(err)
}

func (k keysClient) Import(name, password, keystore string) (string, original.Error) {
	address, err := k.KeyManager.Import(name, password, keystore)
	return address, original.Wrap(err)
}

func (k keysClient) Export(name, srcPwd, dstPwd string) (string, original.Error) {
	keystore, err := k.KeyManager.Export(name, srcPwd, dstPwd)
	return keystore, original.Wrap(err)
}

func (k keysClient) Delete(name string) original.Error {
	err := k.KeyManager.Delete(name)
	return original.Wrap(err)
}

func (k keysClient) Show(name string) (string, original.Error) {
	address, err := k.KeyManager.Query(name)
	if err != nil {
		return "", original.Wrap(err)
	}
	return address.String(), nil
}

func (k keysClient) RegisterCodec(_ original.Codec) {
	//do nothing
}

func (k keysClient) Name() string {
	return ModuleName
}
