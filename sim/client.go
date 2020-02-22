package sim

import (
	"github.com/irisnet/irishub-sdk-go/client"
	"github.com/irisnet/irishub-sdk-go/types"
)

const (
	Addr    = "faa1d3mf696gvtwq2dfx03ghe64akf6t5vyz6pe3le"
	ValAddr = "iva1x3f572u057lv88mva2q3z40ls8pup9hsg0lxcp"
	PrivKey = "927be78a5f5b63bb95ff34ed9c6e4b39b6af6d2f9f59731452de659cac9b19db"
	NodeURI = "localhost:26657"
	ChainID = "irishub-test"
	Online  = true
	Network = types.Testnet
	Mode    = types.Commit
	Fee     = "600000000000000000iris-atto"
	Gas     = 20000
)

type TestClient struct {
	types.Bank
}

func NewClient() client.Client {
	return client.New(types.SDKConfig{
		NodeURI: NodeURI,
		Network: Network,
		ChainID: ChainID,
		Gas:     Gas,
		Fee:     Fee,
		KeyDAO:  createTestKeyDAO(),
		Mode:    Mode,
		Online:  Online,
	})
}

func (TestClient) GetTestSender() string {
	return Addr
}

func (TestClient) GetTestValidator() string {
	return ValAddr
}

func createTestKeyDAO() TestKeyDAO {
	dao := TestKeyDAO{
		store: map[string]types.KeyStore{},
	}
	keystore := TestKeystore{
		Private: PrivKey,
		Address: Addr,
	}
	_ = dao.Write("test1", keystore)
	return dao
}

type TestKeyDAO struct {
	store map[string]types.KeyStore
}

func (dao TestKeyDAO) Write(name string, keystore types.KeyStore) error {
	dao.store[name] = keystore
	return nil
}

func (dao TestKeyDAO) Read(name string) types.KeyStore {
	return dao.store[name]
}

func (dao TestKeyDAO) Delete(name string) error {
	return nil
}

type TestKeystore struct {
	Private string
	Address string
}

func (t TestKeystore) GetPrivate() string {
	return t.Private
}
func (t TestKeystore) GetAddress() string {
	return t.Address
}
