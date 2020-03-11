package keys_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/sim"
)

type KeysTestSuite struct {
	suite.Suite
	sim.TestClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

func (kts *KeysTestSuite) SetupTest() {
	tc := sim.NewClient()
	kts.TestClient = tc
}

func (kts *KeysTestSuite) TestKeys() {
	name, password := "test2", "1234567890"

	address, mnemonic, err := kts.Keys().Add(name, password)
	require.NoError(kts.T(), err)
	require.NotEmpty(kts.T(), address)
	require.NotEmpty(kts.T(), mnemonic)

	address1, err := kts.Keys().Show(name)
	require.NoError(kts.T(), err)
	require.Equal(kts.T(), address, address1)

	newPwd := "0123456789"
	keystore, err := kts.Keys().Export(name, password, newPwd)
	require.NoError(kts.T(), err)
	fmt.Println(keystore)

	err = kts.Keys().Delete(name, newPwd)
	require.NoError(kts.T(), err)

	address2, err := kts.Keys().Import(name, newPwd, keystore)
	require.NoError(kts.T(), err)
	require.Equal(kts.T(), address, address2)

	err = kts.Keys().Delete(name, newPwd)
	require.NoError(kts.T(), err)

	address3, err := kts.Keys().Recover(name, newPwd, mnemonic)
	require.NoError(kts.T(), err)
	require.Equal(kts.T(), address, address3)
}
