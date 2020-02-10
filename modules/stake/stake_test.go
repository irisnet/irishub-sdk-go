package stake_test

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/stake"
	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
)

type StakeTestSuite struct {
	suite.Suite
	stake.Stake
	sender, validator string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(StakeTestSuite))
}

func (sts *StakeTestSuite) SetupTest() {
	tc := sim.NewTestClient()
	sts.Stake = tc.Stake
	sts.sender = tc.GetTestSender()
	sts.validator = tc.GetTestValidator()
}

func (sts StakeTestSuite) TestQueryDelegation() {
	acc, err := sts.QueryDelegation(sts.sender, sts.validator)
	require.NoError(sts.T(), err)
	fmt.Printf("%v", acc)
}

func (sts StakeTestSuite) TestQueryDelegations() {
	acc, err := sts.QueryDelegations(sts.sender)
	require.NoError(sts.T(), err)
	fmt.Printf("%v", acc)
}

func (sts StakeTestSuite) TestQueryValidator() {
	acc, err := sts.QueryValidator(sts.validator)
	require.NoError(sts.T(), err)
	fmt.Printf("%v", acc)
}

func (sts StakeTestSuite) TestQueryAllValidators() {
	val, err := sts.QueryAllValidators()
	require.NoError(sts.T(), err)
	bz, err := json.Marshal(val)
	require.NoError(sts.T(), err)
	fmt.Printf(string(bz))
}

func (sts StakeTestSuite) TestQueryValidators() {
	acc, err := sts.QueryValidators(1, 100)
	require.NoError(sts.T(), err)
	fmt.Printf("%v", acc)
}

func (sts StakeTestSuite) TestDelegate() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	acc, err := sts.Delegate(sts.validator, coin, baseTx)
	require.NoError(sts.T(), err)
	fmt.Printf("%v", acc)
}

func (sts StakeTestSuite) TestUnDelegate() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	acc, err := sts.Undelegate(sts.validator, coin, baseTx)
	require.NoError(sts.T(), err)
	fmt.Printf("%v", acc)
}
