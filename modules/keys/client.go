package keys

import "github.com/irisnet/irishub-sdk-go/types"

var (
	_ types.Keys = keysClient{}
)

type keysClient struct {
	types.KeyManager
}

//TODO
func (k keysClient) Add(name, password string) (address string, mnemonic string, err error) {
	panic("implement me")
}

//TODO
func (k keysClient) Recover(name, password, mnemonic string, derive bool, index int, salt string) (address string, err error) {
	panic("implement me")
}

//TODO
func (k keysClient) Import(name, password, keystore string) (address string, err error) {
	panic("implement me")
}

//TODO
func (k keysClient) Export(name, password, encryptKeystorePwd string) (keystore string, err error) {
	panic("implement me")
}

//TODO
func (k keysClient) Delete(name, password string) error {
	panic("implement me")
}
