package types

type Keystore struct {
	Address string `json:"address"`
	Id      string `json:"id"`
	Version int    `json:"version"`
	Crypto  Crypto `json:"crypto"`
}

type Crypto struct {
	CipherText   string       `json:"cipher_text"`
	CipherParams CipherParams `json:"cipher_params"`
	Cipher       string       `json:"cipher"`
	Kdf          string       `json:"kdf"`
	KdfParams    KdfParams    `json:"kdf_params"`
	Mac          string       `json:"mac"`
}

type CipherParams struct {
	IV string `json:"iv"`
}

type KdfParams struct {
	DkLen int    `json:"dk_len"`
	Salt  string `json:"salt"`
	C     int16  `json:"c"`
	Prf   string `json:"prf"`
}
