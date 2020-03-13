package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/irisnet/irishub-sdk-go/client"
	"github.com/irisnet/irishub-sdk-go/types"
)

const (
	NodeURI = "localhost:26657"
	ChainID = "test"
	Online  = true
	Network = types.Testnet
	Mode    = types.Commit
	Fee     = "600000000000000000iris-atto"
	Gas     = 20000
)

type TestClient struct {
	sender types.AccAddress
	client.SDKClient
}

func NewClient() TestClient {
	tc := TestClient{}
	keystore := getKeystore()
	fees, err := types.ParseCoins(Fee)
	if err != nil {
		panic(err)
	}

	c := client.NewSDKClient(types.SDKConfig{
		NodeURI:   NodeURI,
		Network:   Network,
		ChainID:   ChainID,
		Gas:       Gas,
		Fee:       fees,
		KeyDAO:    createTestKeyDAO(),
		Mode:      Mode,
		Online:    Online,
		StoreType: types.Key,
		Level:     "debug",
	})

	//init account
	address, err := c.Keys().Import("test1", tc.Password(), keystore)
	if err != nil {
		panic(err)
	}

	tc.SDKClient = c
	tc.sender = types.MustAccAddressFromBech32(address)
	return tc
}

func (tc TestClient) Sender() types.AccAddress {
	return tc.sender
}

func (tc TestClient) Password() string {
	return "11111111"
}

func createTestKeyDAO() types.KeyDAO {
	dao := TestKeyDAO{
		store: map[string]types.Store{},
	}
	return types.NewKeyDAO(&dao, nil)
}

type TestKeyDAO struct {
	store map[string]types.Store
}

func (dao *TestKeyDAO) Write(name string, store types.Store) error {
	dao.store[name] = store
	return nil
}

func (dao *TestKeyDAO) Read(name string) (types.Store, error) {
	return dao.store[name], nil
}

func (dao *TestKeyDAO) Delete(name string) error {
	delete(dao.store, name)
	return nil
}

func getKeystore() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path = filepath.Dir(path)
	path = strings.TrimRight(path, "modules")
	path = filepath.Join(path, "test/scripts/keystore1.json")
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bz)
}
