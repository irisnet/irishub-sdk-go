package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

type Keys interface {
	original.Module
	Add(name, password string) (address string, mnemonic string, err original.Error)
	Recover(name, password, mnemonic string) (address string, err original.Error)
	Import(name, password, keystore string) (address string, err original.Error)
	Export(name, password, encryptKeystorePwd string) (keystore string, err original.Error)
	Delete(name string) original.Error
	Show(name string) (string, original.Error)
}
