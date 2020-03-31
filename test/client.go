package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	sdk "github.com/irisnet/irishub-sdk-go"
	"github.com/irisnet/irishub-sdk-go/types"
)

const (
	NodeURI = "localhost:26657"
	ChainID = "test"
	Network = types.Testnet
	Mode    = types.Commit
	Fee     = "0.6iris"
	Gas     = 20000
)

type MockClient struct {
	sdk.Client
	user MockAccount
}

type MockAccount struct {
	Name, Password string
	Address        types.AccAddress
}

func NewMockClient() MockClient {
	tc := MockClient{
		user: MockAccount{
			Name:     "test1",
			Password: "11111111",
		},
	}
	fees, err := types.ParseDecCoins(Fee)
	if err != nil {
		panic(err)
	}

	c := sdk.NewClient(types.ClientConfig{
		NodeURI:   NodeURI,
		Network:   Network,
		ChainID:   ChainID,
		Gas:       Gas,
		Fee:       fees,
		KeyDAO:    types.NewDefaultKeyDAO(&Memory{}),
		Mode:      Mode,
		StoreType: types.Key,
		Timeout:   5 * time.Second,
		Level:     "info",
	})

	//init account
	keystore := getKeystore()
	address, err := c.Keys().Import("test1", tc.user.Password, keystore)
	if err != nil {
		panic(err)
	}

	tc.Client = c
	tc.user.Address = types.MustAccAddressFromBech32(address)
	return tc
}
func (tc MockClient) Account() MockAccount {
	return tc.user
}

type Memory map[string]types.Store

func (m Memory) Write(name string, store types.Store) error {
	m[name] = store
	return nil
}

func (m Memory) Read(name string) (types.Store, error) {
	return m[name], nil
}

func (m Memory) Delete(name string) error {
	delete(m, name)
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
