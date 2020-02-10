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
	ets.Bank = tc.Client
	ets.sender = tc.GetTestSender()
	ets.validator = tc.GetTestValidator()
}

func (ets EventTestSuite) TestSubscribeNewBlock() {
	err := ets.SubscribeNewBlock(func(sub types.Subscription) {
		bz, _ := json.Marshal(sub.GetData())
		fmt.Println(string(bz))
		sub.Unsubscribe()
	})
	require.NoError(ets.T(), err)
	time.Sleep(20 * time.Second)
}

func (ets EventTestSuite) TestSubscribeTx() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	coins := types.NewCoins(coin)
	to := "iaa120v5ev44cwft687l0jcr5ec3vh2626vsschv7e"
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Async,
	}
	result, err := ets.Send(to, coins, baseTx)
	require.NoError(ets.T(), err)
	require.True(ets.T(), result.IsSuccess())
	ch := make(chan int)
	query := types.EventQueryTxFor(result.GetHash())
	err = ets.SubscribeTx(query, func(sub types.Subscription) {
		tx := sub.GetData().(types.EventDataTx)
		require.Equal(ets.T(), result.GetHash(), tx.Hash)
		sub.Unsubscribe()
		ch <- 1
	})
	require.NoError(ets.T(), err)
	<-ch
}
