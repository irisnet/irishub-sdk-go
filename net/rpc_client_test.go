package net_test

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type WSClientTestSuite struct {
	suite.Suite
	types.WSClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(WSClientTestSuite))
}

func (wts *WSClientTestSuite) SetupTest() {
	c := sim.NewClient()
	wts.WSClient = c
}

func (wts WSClientTestSuite) TestSubscribeNewBlock() {
	subscription, err := wts.SubscribeNewBlock(func(data types.EventDataNewBlock) {
		bz, _ := json.Marshal(data)
		fmt.Println(string(bz))
	})
	require.NoError(wts.T(), err)
	time.Sleep(20 * time.Second)
	err = wts.Unscribe(subscription)
	require.NoError(wts.T(), err)
}

func (wts WSClientTestSuite) TestSubscribeValidatorSetUpdates() {
	subscription, err := wts.SubscribeValidatorSetUpdates(func(data types.EventDataValidatorSetUpdates) {
		bz, _ := json.Marshal(data)
		fmt.Println(string(bz))
	})
	require.NoError(wts.T(), err)
	time.Sleep(20 * time.Second)
	err = wts.Unscribe(subscription)
	require.NoError(wts.T(), err)
}
