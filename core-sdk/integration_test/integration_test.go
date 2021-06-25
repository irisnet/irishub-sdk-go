package integration_test

import (
	sdk "github.com/irisnet/irishub-sdk-go"
	"github.com/irisnet/irishub-sdk-go/common/log"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/store"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	nodeURI  = "tcp://localhost:26657"
	grpcAddr = "localhost:9090"
	chainID  = "test"
	charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	addr     = "iaa1w9lvhwlvkwqvg08q84n2k4nn896u9pqx93velx"
)

type IntegrationTestSuite struct {
	suite.Suite
	sdk.IRISHUBClient
	r            *rand.Rand
	rootAccount  MockAccount
	randAccounts []MockAccount
}

type SubTest struct {
	testName string
	testCase func(s IntegrationTestSuite)
}

// MockAccount define a account for test
type MockAccount struct {
	Name, Password string
	Address        types.AccAddress
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {

	options := []types.Option{
		types.KeyDAOOption(store.NewMemory(nil)),
		types.TimeoutOption(10),
		types.TokenManagerOption(TokenManager{}),
	}
	cfg, err := types.NewClientConfig(nodeURI, grpcAddr, chainID, options...)
	if err != nil {
		panic(err)
	}

	s.IRISHUBClient = sdk.NewIRISHUBClient(cfg)
	s.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.rootAccount = MockAccount{
		Name:     "validator",
		Password: "1234567890",
		Address:  types.MustAccAddressFromBech32(addr),
	}
	s.SetLogger(log.NewLogger(log.Config{
		Format: log.FormatJSON,
		Level:  log.DebugLevel,
	}))
	s.initAccount()
}

func (s *IntegrationTestSuite) initAccount() {
	_, err := s.Import(
		s.Account().Name,
		s.Account().Password,
		string(getPrivKeyArmor()),
	)
	if err != nil {
		panic(err)
	}

	//var receipts bank.Receipts
	for i := 0; i < 5; i++ {
		name := s.RandStringOfLength(10)
		pwd := s.RandStringOfLength(16)
		address, _, err := s.Add(name, "11111111")
		if err != nil {
			panic("generate test account failed")
		}

		s.randAccounts = append(s.randAccounts, MockAccount{
			Name:     name,
			Password: pwd,
			Address:  types.MustAccAddressFromBech32(address),
		})
	}
}

// RandStringOfLength return a random string
func (s *IntegrationTestSuite) RandStringOfLength(l int) string {
	var result []byte
	bytes := []byte(charset)
	for i := 0; i < l; i++ {
		result = append(result, bytes[s.r.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandAccount return a random test account
func (s *IntegrationTestSuite) GetRandAccount() MockAccount {
	return s.randAccounts[s.r.Intn(len(s.randAccounts))]
}

// Account return a test account
func (s *IntegrationTestSuite) Account() MockAccount {
	return s.rootAccount
}

func getPrivKeyArmor() []byte {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path = filepath.Dir(path)
	path = filepath.Join(path, "integration_test/scripts/priv.key")
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bz
}

type TokenManager struct{}

func (TokenManager TokenManager) QueryToken(denom string) (types.Token, error) {
	return types.Token{}, nil
}

func (TokenManager TokenManager) SaveTokens(tokens ...types.Token) {
	return
}

func (TokenManager TokenManager) ToMinCoin(coins ...types.DecCoin) (types.Coins, types.Error) {

	for i, _ := range coins {
		if coins[i].Denom == "iris" {
			coins[i].Denom = "uiris"
			coins[i].Amount = coins[i].Amount.MulInt(types.NewIntWithDecimal(1, 6))
		}
	}
	ucoins, _ := types.DecCoins(coins).TruncateDecimal()
	return ucoins, nil
}

func (TokenManager TokenManager) ToMainCoin(coins ...types.Coin) (types.DecCoins, types.Error) {
	decCoins := make(types.DecCoins, len(coins), 0)
	for _, coin := range coins {
		if coin.Denom == "uiris" {
			amtount := types.NewDecFromInt(coin.Amount).Mul(types.NewDecWithPrec(1, 6))
			decCoins = append(decCoins, types.NewDecCoinFromDec("iris", amtount))
		}
	}
	return decCoins, nil
}
