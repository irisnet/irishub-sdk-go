package integrationtest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/irisnet/htlc-sdk-go"
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

	to := s.GetRandAccount().Address
	createHTLCRequest := htlc.CreateHTLCRequest{
		To:                   to.String(),
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock,
	}
	res, err := s.HTLC.CreateHTLC(createHTLCRequest, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	hashLockBytes, _ := hex.DecodeString(hashLock)
	minCoins, _ := s.ToMinCoin(amount...)
	htlcId := hex.EncodeToString(tmhash.Sum(append(append(append(hashLockBytes, s.Account().Address...), to...), []byte(minCoins.Sort().String())...)))

	queryHTLCResp, err := s.HTLC.QueryHTLC(htlcId)
	require.NoError(s.T(), err)
	require.Equal(s.T(), receiverOnOtherChain, queryHTLCResp.ReceiverOnOtherChain)

	res, err = s.HTLC.ClaimHTLC(htlcId, secret, baseTx)
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

func (s IntegrationTestSuite) TestQueryParams() {
	res, err := s.HTLC.QueryParams()
	require.NoError(s.T(), err)
	data, _ := json.Marshal(res)
	fmt.Println(string(data))
}
