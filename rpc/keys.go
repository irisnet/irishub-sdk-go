package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Keys interface {
	sdk.Module
	Add(name, password string) (address string, mnemonic string, err sdk.Error)
	Recover(name, password, mnemonic string) (address string, err sdk.Error)
	Import(name, password, keystore string) (address string, err sdk.Error)
	Export(name, password, encryptKeystorePwd string) (keystore string, err sdk.Error)
	Delete(name, password string) sdk.Error
	Show(name string) (string, sdk.Error)
}
