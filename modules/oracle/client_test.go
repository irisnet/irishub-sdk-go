package oracle_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type OracleTestSuite struct {
	suite.Suite
	sim.TestClient
	log.Logger
	serviceName string
	baseTx      sdk.BaseTx
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (ots *OracleTestSuite) SetupTest() {
	ots.TestClient = sim.NewClient()
	ots.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func (ots *OracleTestSuite) SetupService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":[{"denom":"iris-atto","amount":"1000000000000000000"}]}`
	output := `{"last":"100"}`
	serviceName := generateServiceName()

	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	definition := sdk.ServiceDefinitionRequest{
		BaseTx:            baseTx,
		ServiceName:       serviceName,
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := ots.DefineService(definition)
	require.NoError(ots.T(), err)
	require.True(ots.T(), result.IsSuccess())

	deposit, _ := sdk.ParseCoins("20000000000000000000000iris-atto")
	binding := sdk.ServiceBindingRequest{
		BaseTx:       baseTx,
		ServiceName:  definition.ServiceName,
		Deposit:      deposit,
		Pricing:      pricing,
		WithdrawAddr: "",
	}
	result, err = ots.BindService(binding)
	require.NoError(ots.T(), err)
	require.True(ots.T(), result.IsSuccess())

	err = ots.RegisterSingleInvocationListener(serviceName,
		func(input string) (string, string) {
			ots.Info("Service received request", "input", input, "output", output)
			return output, ""
		}, baseTx)

	require.NoError(ots.T(), err)

	ots.serviceName = serviceName
	ots.baseTx = baseTx
}

func (ots *OracleTestSuite) TestFeed() {
	//before
	ots.SetupService()

	input := `{"pair":"iris-usdt"}`
	feedName := generateFeedName(ots.serviceName)
	serviceFeeCap, _ := sdk.ParseCoins("1000000000000000000iris-atto")

	createFeedReq := sdk.FeedCreateRequest{
		BaseTx:            ots.baseTx,
		FeedName:          feedName,
		LatestHistory:     5,
		Description:       "fetch USDT-CNY ",
		ServiceName:       ots.serviceName,
		Providers:         []string{ots.Sender().String()},
		Input:             input,
		Timeout:           3,
		ServiceFeeCap:     serviceFeeCap,
		RepeatedFrequency: 5,
		RepeatedTotal:     -1,
		AggregateFunc:     "avg",
		ValueJsonPath:     "last",
		ResponseThreshold: 1,
	}
	result, err := ots.CreateFeed(createFeedReq)
	require.NoError(ots.T(), err)
	require.True(ots.T(), result.IsSuccess())

	_, err = ots.QueryFeed(feedName)
	require.NoError(ots.T(), err)

	result, err = ots.StartFeed(feedName, ots.baseTx)
	require.NoError(ots.T(), err)
	require.True(ots.T(), result.IsSuccess())

	execTimer(func() {
		result, err := ots.QueryFeedValue(feedName)
		require.NoError(ots.T(), err)
		ots.Info("Query feed value", "feedName", feedName, "result", result)
	})
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}

func generateFeedName(serviceName string) string {
	return fmt.Sprintf("feed-%s", serviceName)
}

func execTimer(call func()) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		call()
	}
}
