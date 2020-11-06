package integration_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/nft"
	"github.com/irisnet/irishub-sdk-go/types"
	"strings"
)

var (
	denomId   string
	nftId     string
	recipient string

	anotherName, anotherPasswd string
)

func (s IntegrationTestSuite) TestNft() {
	cases := []SubTest{
		{
			"TestIssueDenom",
			issueDenom,
		},
		{
			"TestQueryDenom",
			queryDenom,
		},
		{
			"TestMintNft",
			mintNft,
		},
		{
			"TestQueryNft",
			queryNft,
		},
		{
			"TestEditNft",
			editNft,
		},
		{
			"TestTransferNft",
			transferNft,
		},
		{
			"TestQuerySupply",
			querySupply,
		},
		{
			"TestQueryOwner",
			queryOwner,
		},
		{
			"TestQueryDenoms",
			queryDenoms,
		},
		{
			"TestBurnNft",
			burnNft,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}
}

func issueDenom(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
	}

	schema := `
{
  "$id": "https://example.com/nft.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "description": "nft test",
  "type": "object",
  "properties": {
    "id": {
      "type": "number"
    },
    "name": {
      "type": "string"
    }
  }
}
`

	denomId = strings.ToLower(s.RandStringOfLength(10))
	name := s.RandStringOfLength(4)
	issueReq := nft.IssueDenomRequest{
		ID:     denomId,
		Name:   name,
		Schema: schema,
	}
	denom, err := s.NFT.IssueDenom(issueReq, baseTx)
	s.NoError(err)
	s.NotEmpty(denom.Hash)
}

func queryDenom(s IntegrationTestSuite) {
	denomResp, err := s.NFT.QueryDenom(denomId)
	s.NoError(err)
	s.Equal(denomId, denomResp.ID)
	fmt.Println(denomResp)
}

func mintNft(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
	}

	nftId = strings.ToLower(s.RandStringOfLength(10))
	name := s.RandStringOfLength(4)
	data := `
{
  "id": 1,
  "name": "hello nft"
}
`
	mintNftReq := nft.MintNFTRequest{
		Denom: denomId,
		ID:    nftId,
		Name:  name,
		URI:   fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
		Data:  data,
	}
	res, err := s.NFT.MintNFT(mintNftReq, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
}

func queryNft(s IntegrationTestSuite) {
	nftResp, err := s.NFT.QueryNFT(denomId, nftId)
	s.NoError(err)
	s.NotEmpty(nftResp.Data)
	s.Equal(nftId, nftResp.ID)
}

func editNft(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
	}

	editReq := nft.EditNFTRequest{
		Denom: denomId,
		ID:    nftId,
		URI:   fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
	}

	res, err := s.NFT.EditNFT(editReq, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
	s.Greater(res.GasUsed, int64(0))
}

func transferNft(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
	}
	coins, err := types.ParseDecCoins("100iris")
	s.NoError(err)

	anotherName = s.RandStringOfLength(10)
	anotherPasswd = "11111111"
	recipient, _, err = s.Key.Add(anotherName, anotherPasswd)
	s.NoError(err)
	_, err = s.Bank.Send(recipient, coins, baseTx)
	s.NoError(err)

	transferNftReq := nft.TransferNFTRequest{
		Denom:     denomId,
		ID:        nftId,
		URI:       fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
		Recipient: recipient,
	}
	res, err := s.NFT.TransferNFT(transferNftReq, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
	fmt.Println(res)
}

func querySupply(s IntegrationTestSuite) {
	supplyRes, err := s.NFT.QuerySupply(denomId, s.Account().Address.String())
	s.NoError(err)
	fmt.Println(supplyRes)
}

func queryOwner(s IntegrationTestSuite) {
	creator := s.Account().Address.String()
	owner, err := s.NFT.QueryOwner(creator, denomId)
	s.NoError(err)
	s.Len(owner.IDCs, 1)
	s.Len(owner.IDCs[0].TokenIDs, 1)
}

func queryDenoms(s IntegrationTestSuite) {
	denoms, err := s.NFT.QueryDenoms()
	s.NoError(err)
	s.NotEmpty(denoms)

	var flag bool
	for _, denom := range denoms {
		if denom.ID == denomId {
			flag = true
		}
	}
	s.Equal(true, flag)
}

func burnNft(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     anotherName,
		Password: anotherPasswd,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
	}

	burnReq := nft.BurnNFTRequest{
		Denom: denomId,
		ID:    nftId,
	}

	res, err := s.NFT.BurnNFT(burnReq, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
	s.Greater(res.GasUsed, int64(0))
}
