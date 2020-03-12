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
	Write(name string, store Store) error
	Read(name, password string) (Store, error)
	Delete(name, password string) error
}

type KeyManager interface {
	Sign(name, password string, data []byte) (Signature, error)
	QueryAddress(name, password string) (addr AccAddress, err error)
	Insert(name, password string) (string, string, error)
	Recover(name, password, mnemonic string) (string, error)
	Import(name, password string, keystore string) (address string, err error)
	Export(name, password, encryptKeystorePwd string) (keystore string, err error)
	Delete(name, password string) error
	Query(name string) (address string, err error)
}
