package types

type StoreType int

const (
	Keystore StoreType = 0
	PrivKey  StoreType = 1
)

var (
	_ Store = PrivKeyInfo{}
	_ Store = KeystoreInfo{}
)

type Store interface{
	GetType() StoreType
}
type PrivKeyInfo struct {
	PrivKey string `json:"priv_key"`
	Address string `json:"address"`
}

func (p PrivKeyInfo) GetType() StoreType {
	return PrivKey
}

type KeystoreInfo struct {
	Keystore string `json:"keystore"`
}

func (k KeystoreInfo) GetType() StoreType {
	return Keystore
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
	Encrypt(data string, password string) (string, error)
	Decrypt(data string, password string) (string, error)
}

type defaultKeyDAOImpl struct {
	AccountAccess
	AES
}

// NewKeyDaoWithAES return a KeyDAO object. by default,the SDK's own encryption(AES) method is used
func NewKeyDaoWithAES(account AccountAccess) KeyDAO {
	return defaultKeyDAOImpl{
		AccountAccess: account,
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
