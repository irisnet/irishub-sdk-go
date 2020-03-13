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
	PrivKey string `json:"priv_key"`
	Address string `json:"address"`
}
type KeystoreInfo struct {
	Keystore string `json:"keystore"`
}

type KeyDAO interface {
	AccountAccess
	Crypto
}

type AccountAccess interface {
	Write(name string, store Store) error
	Read(name string) (Store, error)
	Delete(name string) error
}

type Crypto interface {
	Encrypt(privKey string, password string) (string, error)
	Decrypt(privKey string, password string) (string, error)
}

type keyDAOImpl struct {
	AccountAccess
	Crypto
}

func NewKeyDAO(account AccountAccess, crypto Crypto) KeyDAO {
	if crypto == nil {
		//default
		crypto = AES{}
	}
	return keyDAOImpl{
		AccountAccess: account,
		Crypto:        crypto,
	}
}

type KeyManager interface {
	Sign(name, password string, data []byte) (Signature, error)
	Insert(name, password string) (string, string, error)
	Recover(name, password, mnemonic string) (string, error)
	Import(name, password string, keystore string) (address string, err error)
	Export(name, password, encryptKeystorePwd string) (keystore string, err error)
	Delete(name string) error
	Query(name string) (address AccAddress, err error)
}
