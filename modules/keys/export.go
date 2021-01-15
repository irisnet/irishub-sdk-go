package keys

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

// Import: (newKeystore/ privKeyArmor) after IrisHubV1.0
// ImportFromOlderKeystore: this keystore before IrisHubV0.16
type Client interface {
	Add(name, password string) (address string, mnemonic string, err sdk.Error)
	Recover(name, password, mnemonic string) (address string, err sdk.Error)
	Import(name, password, privKeyArmor string) (address string, err sdk.Error)
	ImportFromOlderKeystore(name, password, keystore string) (address string, err sdk.Error)
	Export(name, password string) (privKeyArmor string, err sdk.Error)
	Delete(name, password string) sdk.Error
	Show(name, password string) (string, sdk.Error)
}
