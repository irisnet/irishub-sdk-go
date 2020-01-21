package client_test

import (
	"testing"

	"github.com/irisnet/irishub-sdk-go/client"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	client.Client
}

func (c *ClientTestSuite) SetupTest() {
	cfg := types.SDKConfig{
		NodeURI: "localhost:26657",
		Network: types.Mainnet,
		ChainID: "irishub-test",
		Gas:     20000,
		Fee:     "600000000000000000iris-atto",
		KeyDAO:  createTestKeyDAO(),
		Mode:    types.Commit,
		Online:  true,
	}
	c.Client = client.NewClient(cfg)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func createTestKeyDAO() TestKeyDAO {
	dao := TestKeyDAO{
		store: map[string]types.KeyStore{},
	}
	keystore := TestKeystore{
		Private: "8D03FEDB094224959DD12016D24782429216246BC03084211C0305F9767C3C38",
		Address: "iaa1x3f572u057lv88mva2q3z40ls8pup9hsa74f9x",
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
