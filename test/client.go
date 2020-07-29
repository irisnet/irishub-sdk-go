package test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	sdk "github.com/irisnet/irishub-sdk-go"
	"github.com/irisnet/irishub-sdk-go/types"
)

const (
	nodeURI = "http://10.1.4.185:36657"
	//nodeURI = "http://localhost:26657"
	chainID = "test"
	network = types.Mainnet
	mode    = types.Commit
	gas     = 20000

	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	mock *MockClient
	lock sync.Mutex
	rnd  = rand.NewSource(64)
)

type MockClient struct {
	sdk.Client
	rootUser *MockAccount
}

type MockAccount struct {
	Name, Password string
	Address        types.AccAddress
}

func GetMock() *MockClient {
	lock.Lock()
	defer lock.Unlock()

	if mock == nil {
		c := newMockClient()
		mock = &c
	}
	return mock
}

func newMockClient() MockClient {
	path := filepath.Join(getPWD(), "test")
	c := sdk.NewClient(types.ClientConfig{
		NodeURI:   nodeURI,
		Network:   network,
		ChainID:   chainID,
		Gas:       gas,
		KeyDAO:    types.NewMemoryDB(), //default keybase
		Mode:      mode,
		StoreType: types.PrivKey,
		Timeout:   10 * time.Second,
		Level:     "info",
		DBRootDir: path,
	})

	tc := MockClient{
		Client: c,
		rootUser: &MockAccount{
			Name:     "test1",
			Password: "11111111",
		},
	}

	tc.init()
	return tc
}

func (tc MockClient) init() {
	address, err := tc.Keys().Show(tc.rootUser.Name)
	if err != nil {
		keystore := getKeystore()
		address, err = tc.Keys().Import(tc.rootUser.Name, tc.rootUser.Password, keystore)
		if err != nil {
			panic(err)
		}
	}
	tc.rootUser.Address = types.MustAccAddressFromBech32(address)
}

func (tc MockClient) Account() MockAccount {
	return *tc.rootUser
}

func (tc MockClient) RandStringOfLength(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rnd.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rnd.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func getKeystore() string {
	path := filepath.Join(getPWD(), "test/scripts/keystore1.json")
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

func getPWD() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path = filepath.Dir(path)
	path = strings.TrimRight(path, "modules")
	return path
}
