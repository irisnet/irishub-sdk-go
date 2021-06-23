package sm2

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"log"
)

//Generate key pair
func GenerateKey() *sm2.PrivateKey {
	//rand.Reader
	privateKey, err := sm2.GenerateKey(rand.Reader) // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

//Gets the public key from the private key
func GetPublickey(privateKey *sm2.PrivateKey) *sm2.PublicKey {
	return &privateKey.PublicKey
}

//SM2 public key encryption
func PublicKeyEncrypt(publicKey *sm2.PublicKey, msg []byte) []byte {
	ciphertxt, err := publicKey.EncryptAsn1(msg, rand.Reader) //sm2公钥加密
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("加密结果:%x\n", ciphertxt)
	return ciphertxt
}

//SM2 Private key decrypt
func PrivateKeyDecrypt(privateKey *sm2.PrivateKey, ciphertxt []byte) (plaintext []byte) {
	plaintext, err := privateKey.DecryptAsn1(ciphertxt) //sm2私钥解密
	if err != nil {
		log.Fatal(err)
	}
	//if !bytes.Equal(msg, plaintext) {
	//	//log.Fatal("原文不匹配")
	//}
	return plaintext
}

//Text matching: compares decrypted information with original text
//msg: original text
//plaintext: Decrypted plaintext
func TextMatch(msg, plaintext []byte) bool {
	return bytes.Equal(msg, plaintext)
}

//PrivateKey Sign
func PrivateKeySign(privateKey *sm2.PrivateKey, msg []byte) ([]byte, error) {
	sign, err := privateKey.Sign(rand.Reader, msg, nil) //sm2私钥签名
	if err != nil {
		log.Fatal(err)
	}
	return sign, err
}

// PublicKey Verify
func PublicKeyVerify(publicKey *sm2.PublicKey, msg, sign []byte) bool {
	isok := publicKey.Verify(msg, sign) //sm2公钥验签
	//fmt.Printf("Verified: %v\n", isok)
	return isok
}

func Sm2() {
	privateKey, err := sm2.GenerateKey(rand.Reader) // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("privateKey", privateKey)
	publicKey := &privateKey.PublicKey //从私钥中获取公钥
	fmt.Println("publicKey ", publicKey)

	msg := []byte("Tongji Fintech Research Institute")

	ciphertxt, err := publicKey.EncryptAsn1(msg, rand.Reader) //sm2公钥加密
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("加密结果:%x\n", ciphertxt)

	plaintext, err := privateKey.DecryptAsn1(ciphertxt) //sm2私钥解密
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(msg, plaintext) {
		log.Fatal("原文不匹配")
	}

	sign, err := privateKey.Sign(rand.Reader, msg, nil) //sm2私钥签名
	if err != nil {
		log.Fatal(err)
	}
	isok := publicKey.Verify(msg, sign) //sm2公钥验签
	fmt.Printf("Verified: %v\n", isok)
}
