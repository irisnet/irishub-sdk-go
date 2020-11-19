package integration_test

import (
	"github.com/irisnet/irishub-sdk-go/modules/token"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"strings"
)

type TestToken struct {
	symbol    string
	scale     uint32
	minUnit   string
	recipient string
}

var testToken TestToken

func (s IntegrationTestSuite) TestToken() {
	cases := []SubTest{
		{
			"TestIssueToken",
			issueToken,
		},
		{
			"TestMintToken",
			mintToken,
		},
		{
			"TestEditToken",
			editToken,
		},
		{
			"TestQueryTokens",
			queryTokens,
		},
		{
			"TestTransferToken",
			transferToken,
		},
		{
			"TestQueryToken",
			queryToken,
		},
		{
			"TestQueryFees",
			queryFees,
		},
		{
			"TestQueryParams",
			queryParams,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}
}

func issueToken(s IntegrationTestSuite) {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
	}

	// init testToken
	testToken.symbol = strings.ToLower(s.RandStringOfLength(8))
	testToken.scale = 9
	testToken.minUnit = strings.ToLower(s.RandStringOfLength(3))

	issueTokenReq := token.IssueTokenRequest{
		Symbol:        testToken.symbol,
		Name:          strings.ToLower(s.RandStringOfLength(8)),
		Scale:         testToken.scale,
		MinUnit:       testToken.minUnit,
		InitialSupply: 10000000,
		MaxSupply:     20000000,
		Mintable:      true,
	}
	res, err := s.Token.IssueToken(issueTokenReq, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
}

func mintToken(s IntegrationTestSuite) {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
	}

	recipient := s.GetRandAccount().Address.String()
	res, err := s.Token.MintToken(testToken.symbol, 1000, recipient, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)

	account, err := s.Bank.QueryAccount(recipient)
	s.NoError(err)

	amt := sdk.NewIntWithDecimal(1000, int(testToken.scale))
	s.Equal(amt, account.Coins.AmountOf(testToken.minUnit))
}

func editToken(s IntegrationTestSuite) {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
	}

	editTokenReq := token.EditTokenRequest{
		Symbol:    testToken.symbol,
		Name:      "my network",
		MaxSupply: 2000000000,
		Mintable:  false,
	}
	res, err := s.Token.EditToken(editTokenReq, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
}

func queryTokens(s IntegrationTestSuite) {
	tokens, err := s.Token.QueryTokens(s.Account().Address.String())
	s.NoError(err)

	var flag bool
	for _, token := range tokens {
		if testToken.symbol == token.Symbol {
			flag = true
		}
	}
	s.True(flag)
}

func transferToken(s IntegrationTestSuite) {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Password: s.Account().Password,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
	}

	testToken.recipient = s.GetRandAccount().Address.String()
	res, err := s.Token.TransferToken(testToken.recipient, testToken.symbol, baseTx)
	s.NoError(err)
	s.NotEmpty(res.Hash)
}

func queryToken(s IntegrationTestSuite) {
	token, err := s.Token.QueryToken(testToken.symbol)
	s.NoError(err)
	s.Equal(testToken.symbol, token.Symbol)
	s.Equal(testToken.recipient, token.Owner)
	s.NotEmpty(token.Owner)
	s.Greater(token.InitialSupply, uint64(1))
	s.Greater(token.MaxSupply, uint64(1))
}

func queryFees(s IntegrationTestSuite) {
	fees, err := s.Token.QueryFees(testToken.symbol)
	s.NoError(err)
	s.Greater(fees.IssueFee.Amount.Int64(), int64(1))
	s.Greater(fees.MintFee.Amount.Int64(), int64(1))
}

func queryParams(s IntegrationTestSuite) {
	params, err := s.Token.QueryParams()
	s.NoError(err)
	s.NotEmpty(params.IssueTokenBaseFee)
	s.NotEmpty(params.IssueTokenBaseFee)
	s.NotEmpty(params.TokenTaxRate)
}
