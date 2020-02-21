package event_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/event"
	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EventTestSuite struct {
	suite.Suite
	event.Event
	bank.Bank
	sender, validator string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}

func (ets *EventTestSuite) SetupTest() {
	tc := sim.NewTestClient()
	ets.Event = tc.Event
	ets.Bank = tc.Bank
	ets.sender = tc.GetTestSender()
	ets.validator = tc.GetTestValidator()
}

func (ets EventTestSuite) TestSubscribeNewBlock() {
	subscription, err := ets.SubscribeNewBlock(func(data types.EventDataNewBlock) {
		bz, _ := json.Marshal(data)
		fmt.Println(string(bz))
	})
	require.NoError(ets.T(), err)
	time.Sleep(20 * time.Second)
	subscription.Unsubscribe()
}

func (ets EventTestSuite) TestSubscribeTx() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	coins := types.NewCoins(coin)
	to := "faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Async,
	}
	result, err := ets.Send(to, coins, baseTx)
	require.NoError(ets.T(), err)
	require.True(ets.T(), result.IsSuccess())
	ch := make(chan int)
	builder := types.NewEventQueryBuilder()
	builder.AddCondition(types.SenderKey, types.EventValue("faa1d3mf696gvtwq2dfx03ghe64akf6t5vyz6pe3le"))
	subscription, err := ets.SubscribeTx(builder, func(data types.EventDataTx) {
		require.Equal(ets.T(), result.GetHash(), data.Hash)
		bz, _ := json.Marshal(data)
		fmt.Println(string(bz))
		ch <- 1
	})
	require.NoError(ets.T(), err)
	<-ch
	subscription.Unsubscribe()
}
