package oracle_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/suite"
)

type OracleTestSuite struct {
	suite.Suite
	test.MockClient
	*log.Logger
	serviceName string
	baseTx      sdk.BaseTx
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (ots *OracleTestSuite) SetupTest() {
	ots.MockClient = test.NewMockClient()
	ots.Logger = log.NewLogger("info")
}

func (ots *OracleTestSuite) SetupService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":"1iris"}`
	output := `{"last":"100"}`
	testResult := `{"code":200,"message":""}`
	serviceName := generateServiceName()

	baseTx := sdk.BaseTx{
		From:     ots.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: ots.Account().Password,
	}

	definition := rpc.ServiceDefinitionRequest{
		ServiceName:       serviceName,
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := ots.Service().DefineService(definition, baseTx)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	deposit, e := sdk.ParseDecCoins("20000iris")
	require.NoError(ots.T(), e)
	binding := rpc.ServiceBindingRequest{
		ServiceName: definition.ServiceName,
		Deposit:     deposit,
		Pricing:     pricing,
	}
	result, err = ots.Service().BindService(binding, baseTx)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	_, err = ots.Service().RegisterSingleServiceRequestListener(serviceName,
		func(reqCtxID, reqID, input string) (string, string) {
			ots.Info().Str("input", input).
				Str("reqCtxID", reqCtxID).
				Str("output", output).
				Msg("Service received request")
			return output, testResult
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
	serviceFeeCap, _ := sdk.ParseDecCoins("1iris")

	createFeedReq := rpc.FeedCreateRequest{
		BaseTx:            ots.baseTx,
		FeedName:          feedName,
		LatestHistory:     2,
		Description:       "fetch USDT-CNY ",
		ServiceName:       ots.serviceName,
		Providers:         []string{ots.Account().Address.String()},
		Input:             input,
		Timeout:           3,
		ServiceFeeCap:     serviceFeeCap,
		RepeatedFrequency: 5,
		AggregateFunc:     "avg",
		ValueJsonPath:     "last",
		ResponseThreshold: 1,
	}
	result, err := ots.Oracle().CreateFeed(createFeedReq)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	_, err = ots.Oracle().QueryFeed(feedName)
	require.NoError(ots.T(), err)

	result, err = ots.Oracle().StartFeed(feedName, ots.baseTx)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	ch := make(chan rpc.FeedValue)
	err = ots.Oracle().RegisterFeedListener(feedName, func(value rpc.FeedValue) {
		ots.Info().
			Str("feedName", feedName).
			Str("feedValue", value.Data).
			Msg("received feed value")
		ch <- value
	})

	ots.NoError(err)
	for {
		select {
		case v := <-ch:
			result, err := ots.Oracle().QueryFeedValue(feedName)
			require.NoError(ots.T(), err)
			require.EqualValues(ots.T(), v, result[0])

			if len(result) == int(createFeedReq.LatestHistory) {
				goto stop
			}
		case <-time.After(1 * time.Minute):
			require.Panics(ots.T(), func() {}, "test oracle timeout")
		}
	}
stop:
	close(ch)
	_, err = ots.Oracle().PauseFeed(feedName, ots.baseTx)
	require.NoError(ots.T(), err)

	feed, err := ots.Oracle().QueryFeed(feedName)
	require.NoError(ots.T(), err)
	require.Equal(ots.T(), "paused", feed.State)

}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}

func generateFeedName(serviceName string) string {
	return fmt.Sprintf("feed-%s", serviceName)
}
