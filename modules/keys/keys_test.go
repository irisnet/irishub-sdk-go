package keys_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
)

type KeysTestSuite struct {
	suite.Suite
	test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

func (kts *KeysTestSuite) SetupTest() {
	tc := test.NewMockClient()
	kts.MockClient = tc
}

func (kts *KeysTestSuite) TestKeys() {
	name, password := "test2", "1234567890"

	address, mnemonic, err := kts.Keys().Add(name, password)
	kts.NoError(err)
	kts.NotEmpty(address)
	kts.NotEmpty(mnemonic)

	address1, err := kts.Keys().Show(name)
	kts.NoError(err)
	kts.Equal(address, address1)

	newPwd := "01234567891"
	keystore, err := kts.Keys().Export(name, password, newPwd)
	kts.NoError(err)
	fmt.Println(keystore)

	err = kts.Keys().Delete(name)
	kts.NoError(err)

	address2, err := kts.Keys().Import(name, newPwd, keystore)
	kts.NoError(err)
	kts.Equal(address, address2)

	err = kts.Keys().Delete(name)
	kts.NoError(err)

	address3, err := kts.Keys().Recover(name, newPwd, mnemonic)
	kts.NoError(err)
	kts.Equal(address, address3)
}
