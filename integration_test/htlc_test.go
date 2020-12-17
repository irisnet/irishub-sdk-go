package integration_test

import (
	"encoding/hex"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/htlc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func (s IntegrationTestSuite) TestHTLC() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	amount, err := sdk.ParseDecCoins("10iris")
	require.NoError(s.T(), err)
	secret := s.GetSecret()
	hashLock := s.GetHashLock(secret, 0)
	fmt.Println("hashLock: " + hashLock)
	fmt.Println("secret: " + secret)
	receiverOnOtherChain := "0x" + s.RandStringOfLength(14)

	createHTLCRequest := htlc.CreateHTLCRequest{
		To:                   s.GetRandAccount().Address.String(),
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock,
	}
	res, err := s.HTLC.CreateHTLC(createHTLCRequest, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	queryHTLCResp, err := s.HTLC.QueryHTLC(hashLock)
	require.NoError(s.T(), err)
	require.Equal(s.T(), receiverOnOtherChain, queryHTLCResp.ReceiverOnOtherChain)

	res, err = s.HTLC.ClaimHTLC(hashLock, secret, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
}

// GetHashLock calculates the hash lock from the given secret and timestamp
func (s IntegrationTestSuite) GetHashLock(secret string, timestamp uint64) string {
	secretBz, _ := hex.DecodeString(secret)
	if timestamp > 0 {
		return string(tmhash.Sum(append(secretBz, sdk.Uint64ToBigEndian(timestamp)...)))
	}
	sum := tmhash.Sum(secretBz)
	return hex.EncodeToString(sum)
}

func (s IntegrationTestSuite) GetSecret() string {
	random := s.RandStringOfLength(10)
	sum := tmhash.Sum([]byte(random))
	return hex.EncodeToString(sum)
}
