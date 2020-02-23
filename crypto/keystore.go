package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils"
	"github.com/irisnet/irishub-sdk-go/utils/uuid"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
)

func NewKeyStoreKeyManager(file string, auth string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromKeyStore(file, auth)
	return &k, err
}

func (m *keyManager) ExportAsKeyStore(password string) (types.Keystore, error) {
	return generateKeyStore(m.GetPrivKey(), password)
}

func (m *keyManager) recoveryFromKeyStore(keystoreFile string, auth string) error {
	if auth == "" {
		return fmt.Errorf("Password is missing ")
	}
	keyJson, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return err
	}
	var encryptedKey types.Keystore
	err = json.Unmarshal(keyJson, &encryptedKey)
	if err != nil {
		return err
	}
	keyBytes, err := decryptKey(encryptedKey, auth)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	m.privKey = secp256k1.PrivKeySecp256k1(keyBytesArray)
	return nil
}

func generateKeyStore(privateKey crypto.PrivKey, password string) (types.Keystore, error) {
	addr := types.AccAddress(privateKey.PubKey().Address())
	salt, err := utils.GenerateRandomBytes(32)
	if err != nil {
		return types.Keystore{}, err
	}
	iv, err := utils.GenerateRandomBytes(16)
	if err != nil {
		return types.Keystore{}, err
	}

	derivedKey := pbkdf2.Key([]byte(password), salt, 262144, 32, sha256.New)
	encryptKey := derivedKey[:32]
	secpPrivateKey, ok := privateKey.(secp256k1.PrivKeySecp256k1)
	if !ok {
		return types.Keystore{}, fmt.Errorf(" Only PrivKeySecp256k1 key is supported ")
	}
	cipherText, err := aesCTRXOR(encryptKey, secpPrivateKey[:], iv)
	if err != nil {
		return types.Keystore{}, err
	}

	hasher := sha3.NewLegacyKeccak512()
	hasher.Write(derivedKey[16:32])
	hasher.Write(cipherText)
	mac := hasher.Sum(nil)

	id, err := uuid.NewV4()
	if err != nil {
		return types.Keystore{}, err
	}
	return types.Keystore{
		Address: addr.String(),
		Id:      id.String(),
		Version: 1,
		Crypto: types.Crypto{
			CipherText: hex.EncodeToString(cipherText),
			CipherParams: types.CipherParams{
				IV: hex.EncodeToString(iv),
			},
			Cipher: "aes-128-ctr",
			Kdf:    "pbkdf2",
			KdfParams: types.KdfParams{
				DkLen: 32,
				Salt:  hex.EncodeToString(salt),
				C:     262144,
				Prf:   "hmac-sha256",
			},
			Mac: hex.EncodeToString(mac),
		},
	}, nil
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

func decryptKey(keyProtected types.Keystore, auth string) ([]byte, error) {
	mac, err := hex.DecodeString(keyProtected.Crypto.Mac)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(keyProtected.Crypto.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(keyProtected.Crypto.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(keyProtected.Crypto, auth)
	if err != nil {
		return nil, err
	}

	bufferValue := make([]byte, len(cipherText)+16)
	copy(bufferValue[0:16], derivedKey[16:32])
	copy(bufferValue[16:], cipherText[:])
	calculatedMAC := sha256.Sum256([]byte((bufferValue)))
	if !bytes.Equal(calculatedMAC[:], mac) {
		return nil, fmt.Errorf("decrypt failed")
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}

func getKDFKey(crypto types.Crypto, auth string) ([]byte, error) {
	authArray := []byte(auth)
	kdfParams := crypto.KdfParams
	if kdfParams.Salt == "" || kdfParams.DkLen == 0 ||
		kdfParams.C == 0 || kdfParams.Prf == "" {
		return nil, errors.New("invalid KDF params, must contains c, dklen, prf and salt")
	}
	salt, err := hex.DecodeString(kdfParams.Salt)
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(kdfParams.DkLen)

	if crypto.Kdf == "pbkdf2" {
		c := ensureInt(kdfParams.C)
		if kdfParams.Prf != "hmac-sha256" {
			return nil, fmt.Errorf("Unsupported PBKDF2 PRF: %s", kdfParams.Prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}
	return nil, fmt.Errorf("Unsupported KDF: %s", crypto.Kdf)
}

func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}
