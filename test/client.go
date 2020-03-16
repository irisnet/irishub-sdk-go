package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/irisnet/irishub-sdk-go/client"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	NodeURI = "localhost:26657"
	ChainID = "irishub-test"
	Online  = true
	Network = sdk.Testnet
	Mode    = sdk.Commit
	Fee     = "0.6iris"
	Gas     = 20000
)

type MockClient struct {
	client.SDKClient
	user MockAccount
}

type MockAccount struct {
	Name, Password string
	Address        sdk.AccAddress
}

func NewMockClient() MockClient {
	tc := MockClient{
		user: MockAccount{
			Name:     "test1",
			Password: "11111111",
		},
	}
	fees, err := sdk.ParseDecCoins(Fee)
	if err != nil {
		panic(err)
	}

	c := client.NewSDKClient(sdk.SDKConfig{
		NodeURI:   NodeURI,
		Network:   Network,
		ChainID:   ChainID,
		Gas:       Gas,
		Fee:       fees,
		KeyDAO:    sdk.NewDefaultKeyDAO(&Memory{}),
		Mode:      Mode,
		Online:    Online,
		StoreType: sdk.Key,
		Level:     "debug",
	})

	//init account
	keystore := getKeystore()
	address, err := c.Keys().Import("test1", tc.user.Password, keystore)
	if err != nil {
		panic(err)
	}

	tc.SDKClient = c
	tc.user.Address = sdk.MustAccAddressFromBech32(address)
	return tc
}
func (tc MockClient) Account() MockAccount {
	return tc.user
}

type Memory map[string]sdk.Store

func (m Memory) Write(name string, store sdk.Store) error {
	m[name] = store
	return nil
}

func (m Memory) Read(name string) (sdk.Store, error) {
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
