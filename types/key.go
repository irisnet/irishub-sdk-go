package types

type StoreType int

const (
	Keystore StoreType = 0
	Key      StoreType = 1
)

var (
	_ Store = KeyInfo{}
	_ Store = KeystoreInfo{}
)

type Store interface{}
type KeyInfo struct {
	PrivKey  string `json:"priv_key"`
	Address  string `json:"address"`
	Password string `json:"password"`
}
type KeystoreInfo struct {
	KeystoreJSON string `json:"keystore_json"`
}

type KeyDAO interface {
	Write(name string, keystore Store) error
	Read(name string) (Store, error)
	Delete(name string) error
}

type KeyManager interface {
	Sign(name, password string, data []byte) (Signature, error)
	QueryAddress(name, password string) (addr AccAddress, err error)
	Insert(name, password string) (string, string, error)
	Recover(name, password, mnemonic string) (string, error)
}

type Keys interface {
	Add(name, password string) (address string, mnemonic string, err error)
	Recover(name, password, mnemonic string) (address string, err error)
	Import(name, password, keystore string) (address string, err error)
	Export(name, password, encryptKeystorePwd string) (keystore string, err error)
	Delete(name, password string) error
}
