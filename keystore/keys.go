package keystore

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/crypto"
	sdksecp256k1 "github.com/irisnet/irishub-sdk-go/crypto/keys/secp256k1"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// RecoveryAndExportPrivKeyArmor return the new private key armor(after IrisHubV1.0) from a old keystoreFile(before IrisHubV0.16)
func RecoveryAndExportPrivKeyArmor(keystore []byte, password string) (armor string, err error) {
	priv, err := recoveryFromKeyStore(keystore, password)
	if err != nil {
		return "", err
	}
	return exportPrivKeyArmor(priv, password)
}

func recoveryFromKeyStore(keystore []byte, password string) (tmcrypto.PrivKey, error) {
	if password == "" {
		return nil, fmt.Errorf("Password is missing ")
	}

	var encryptedKey EncryptedKeyJSON
	if err := json.Unmarshal(keystore, &encryptedKey); err != nil {
		return nil, err
	}

	keyBytes, err := decryptKey(&encryptedKey, password)
	if err != nil {
		return nil, err
	}

	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("Len of Keyby tes is not equal to 32 ")
	}

	return secp256k1.PrivKey(keyBytes), nil
}

func exportPrivKeyArmor(privKey tmcrypto.PrivKey, password string) (armor string, err error) {
	priv := sdksecp256k1.PrivKey{
		Key: privKey.Bytes(),
	}
	return crypto.EncryptArmorPrivKey(&priv, password, "secp256k1"), nil
}
