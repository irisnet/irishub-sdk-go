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
	user *MockAccount
}

type MockAccount struct {
	Name, Password string
	Address        types.AccAddress
}

func NewMockClient() MockClient {
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
		Mode:      Mode,
		StoreType: types.PrivKey,
		Timeout:   10 * time.Second,
		Level:     "info",
		DBRootDir: "/Users/zhangzhiqiang/Downloads/",
	})

	tc := MockClient{
		Client: c,
		user: &MockAccount{
			Name:     "test1",
			Password: "11111111",
		},
	}

	tc.init()
	return tc
}

func (tc MockClient) init() {
	address, err := tc.Keys().Show(tc.user.Name)
	if err != nil {
		keystore := getKeystore()
		address, err = tc.Keys().Import(tc.user.Name, tc.user.Password, keystore)
		if err != nil {
			panic(err)
		}
	}
	tc.user.Address = types.MustAccAddressFromBech32(address)
}

func (tc MockClient) Account() MockAccount {
	return *tc.user
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
