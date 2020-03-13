package types

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AES struct{}

func (a AES) Encrypt(orig string, key string) (string, error) {
	// 转成字节数组
	origData := []byte(orig)
	k := a.generateKey(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = a.pkcs7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

func (a AES) Decrypt(cryted string, key string) (string, error) {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := a.generateKey(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
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
