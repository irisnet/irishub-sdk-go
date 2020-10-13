package original

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AES struct{}

func (a AES) Encrypt(orig string, key string) (string, error) {
	origData := []byte(orig)
	k := a.generateKey(key)
	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	origData = a.pkcs7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

func (a AES) Decrypt(cryted string, key string) (string, error) {
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := a.generateKey(key)
	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	orig := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(orig, crytedByte)
	orig = a.pkcs7UnPadding(orig)
	return string(orig), nil
}

func (a AES) pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (a AES) pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (a AES) generateKey(key string) []byte {
	kz := []byte(key)
	keyLen := len(kz)
	if keyLen >= 32 {
		return kz[:32]
	}
	if keyLen >= 24 {
		return kz[:24]
	}
	if keyLen >= 16 {
		return kz[:16]
	}
	padding := 16 - keyLen
	padText := bytes.Repeat([]byte{byte(0)}, padding)
	return append(kz, padText...)
}
