package service_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type ServiceTestSuite struct {
	suite.Suite
	sim.TestClient
	log.Logger
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (sts *ServiceTestSuite) SetupTest() {
	sts.TestClient = sim.NewClient()
	sts.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func (sts *ServiceTestSuite) TestService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":[{"denom":"iris-atto","amount":"1000000000000000000"}]}`

	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	definition := sdk.ServiceDefinitionRequest{
		BaseTx:            baseTx,
		ServiceName:       generateServiceName(),
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := sts.DefineService(definition)
	require.NoError(sts.T(), err)
	require.True(sts.T(), result.IsSuccess())

	defi, err := sts.QueryDefinition(definition.ServiceName)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), definition.ServiceName, defi.Name)
	require.Equal(sts.T(), definition.Description, defi.Description)
	require.EqualValues(sts.T(), definition.Tags, defi.Tags)
	require.Equal(sts.T(), definition.AuthorDescription, defi.AuthorDescription)
	require.Equal(sts.T(), definition.Schemas, defi.Schemas)
	require.Equal(sts.T(), sts.Sender(), defi.Author)

	deposit, _ := sdk.ParseCoins("20000000000000000000000iris-atto")
	binding := sdk.ServiceBindingRequest{
		BaseTx:       baseTx,
		ServiceName:  definition.ServiceName,
		Deposit:      deposit,
		Pricing:      pricing,
		WithdrawAddr: "",
	}
	result, err = sts.BindService(binding)
	require.NoError(sts.T(), err)
	require.True(sts.T(), result.IsSuccess())

	bindResp, err := sts.QueryBinding(definition.ServiceName, sts.Sender())
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), binding.ServiceName, bindResp.ServiceName)
	require.Equal(sts.T(), sts.Sender(), bindResp.Provider)
	require.Equal(sts.T(), binding.Deposit.String(), bindResp.Deposit.String())
	require.Equal(sts.T(), binding.Pricing, bindResp.Pricing)

	input := `{"pair":"iris-usdt"}`
	output := `{"last":"1:100"}`

	err = sts.RegisterSingleInvocationListener(definition.ServiceName,
		func(input string) (string, string) {
			sts.Info("Provider received request", "input", input, "output", output)
			return output, ""
		}, baseTx)
	require.NoError(sts.T(), err)

	serviceFeeCap, _ := sdk.ParseCoins("1000000000000000000iris-atto")
	invocation := sdk.ServiceInvocationRequest{
		BaseTx:            baseTx,
		ServiceName:       definition.ServiceName,
		Providers:         []string{sts.Sender().String()},
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           3,
		SuperMode:         false,
		Repeated:          true,
		RepeatedFrequency: 5,
		RepeatedTotal:     -1,
	}
	var requestContextID string
	requestContextID, err = sts.InvokeService(invocation, func(reqCtxID string, response string) {
		require.Equal(sts.T(), reqCtxID, requestContextID)
		require.Equal(sts.T(), output, response)
		sts.Info("Consumer received response", "RequestContextID", requestContextID, "response", response)
	})

	sts.Info("Request success", "RequestContextID", requestContextID)
	require.NoError(sts.T(), err)

	request, err := sts.QueryRequestContext(requestContextID)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), request.ServiceName, invocation.ServiceName)
	require.Equal(sts.T(), request.Input, invocation.Input)

	time.Sleep(10 * time.Minute)
}

func (sts *ServiceTestSuite) TestQueryDefinition() {
	serviceName := "fbfbbc12-5872-11ea-a1dc-186590e06183"
	definition, err := sts.QueryDefinition(serviceName)
	require.NoError(sts.T(), err)
	fmt.Println(definition)
}

func (sts *ServiceTestSuite) TestQueryRequestContext() {
	reqCtxID := "0ef2acf4e1002f1e38c7f0df36be2ad11250aeb6343c6c7d9294f0424deecd96"
	definition, err := sts.QueryRequestContext(reqCtxID)
	require.NoError(sts.T(), err)
	fmt.Println(definition)
}

func (sts *ServiceTestSuite) TestQueryRequest() {
	//request, err := sts.QueryRequests("service-949455000", sts.Sender())
	//require.NoError(sts.T(), err)
	//fmt.Println(request)

	//str :="PehEueM827vYaZfPm7dodnz2/xv2cnNR1NavUkpX0bE="
	//strByte,_ := base64.StdEncoding.DecodeString(str)
	//fmt.Println(hex.EncodeToString(strByte))

	type Tag struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	jsonStr := `[
          {
            "key": "aGVpZ2h0",
            "value": "Mjk5Mw=="
          },
          {
            "key": "cmVxdWVzdC1pZA==",
            "value": "MDYxMjE0NThlNWNjODFkYWExMzBjMjEyNjNmNmVjMmNhODc5OWY2OGViNzA4MmJlZWFmMTQ1YWEwYzNkZGE5NzAwMDAwMDAwMDAwMDAwODkwMDAw"
          },
          {
            "key": "cHJvdmlkZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "Y29uc3VtZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "c2xhc2hlZC1jb2lucw=="
          },
          {
            "key": "cmVxdWVzdC1pZA==",
            "value": "MTkxYjYxNGEyZDk2MzJjODY1MDM2YThkMDZmZjNkNjJlYjk3NWE4MGZmMWFlNWQzZWUwM2QzMDIzNmY1ZGQ0NjAwMDAwMDAwMDAwMDAwNGIwMDAw"
          },
          {
            "key": "cHJvdmlkZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "Y29uc3VtZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "c2xhc2hlZC1jb2lucw=="
          },
          {
            "key": "cmVxdWVzdC1pZA==",
            "value": "NjM2ZjI1ZmRmZjUwZTFjY2MxNzA1OGIxNDMzNzdiOTFkYTQ2OTI1ZjJhNDExYTc3ZGUzZGJlNTQ1YzBjODliNzAwMDAwMDAwMDAwMDAwMjAwMDAw"
          },
          {
            "key": "cHJvdmlkZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "Y29uc3VtZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "c2xhc2hlZC1jb2lucw=="
          },
          {
            "key": "cmVxdWVzdC1pZA==",
            "value": "MjlkOTYwZjNjZWM5NWZkZDFjMTM1YWJmMjUxZjQ1YmFiMzgwNTkwMDlhZDQxMDg4ZDJmYzYxNDZjNTQwNzczMzAwMDAwMDAwMDAwMDAwNWYwMDAw"
          },
          {
            "key": "cHJvdmlkZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "Y29uc3VtZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "c2VydmljZS1uYW1l",
            "value": "c2VydmljZS03Mzk1MTYwMDA="
          },
          {
            "key": "c2VydmljZS1mZWU=",
            "value": "MTAwMDAwMDAwMDAwMDAwMDAwMGlyaXMtYXR0bw=="
          },
          {
            "key": "cmVxdWVzdC1oZWlnaHQ=",
            "value": "Mjk5Mw=="
          },
          {
            "key": "ZXhwaXJhdGlvbi1oZWlnaHQ=",
            "value": "Mjk5Ng=="
          },
          {
            "key": "cmVxdWVzdC1pZA==",
            "value": "M2RlODQ0YjllMzNjZGJiYmQ4Njk5N2NmOWJiNzY4NzY3Y2Y2ZmYxYmY2NzI3MzUxZDRkNmFmNTI0YTU3ZDFiMTAwMDAwMDAwMDAwMDAwMzEwMDAw"
          },
          {
            "key": "cHJvdmlkZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "Y29uc3VtZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "c2VydmljZS1uYW1l",
            "value": "c2VydmljZS0zOTgxMjQwMDA="
          },
          {
            "key": "c2VydmljZS1mZWU=",
            "value": "MTAwMDAwMDAwMDAwMDAwMDAwMGlyaXMtYXR0bw=="
          },
          {
            "key": "cmVxdWVzdC1oZWlnaHQ=",
            "value": "Mjk5Mw=="
          },
          {
            "key": "ZXhwaXJhdGlvbi1oZWlnaHQ=",
            "value": "Mjk5Ng=="
          },
          {
            "key": "cmVxdWVzdC1pZA==",
            "value": "ZjgyZDc4Yzc2ODg4MmJkMGJkYjBiMjM0ZWM0OGQyZjJiOGNlNTcyZmJhYTczMjE4ZTMxYWQ5MDdmNDlkNzg0MTAwMDAwMDAwMDAwMDAwYjEwMDAw"
          },
          {
            "key": "cHJvdmlkZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "Y29uc3VtZXI=",
            "value": "ZmFhMWZ0YzBrc3lkdXljNXJ4YThmc2d5aHpodGxoNnYwajBnbjZ1aDA0"
          },
          {
            "key": "c2VydmljZS1uYW1l",
            "value": "c2VydmljZS0yMzQ5NzEwMDA="
          },
          {
            "key": "c2VydmljZS1mZWU=",
            "value": "MTAwMDAwMDAwMDAwMDAwMDAwMGlyaXMtYXR0bw=="
          },
          {
            "key": "cmVxdWVzdC1oZWlnaHQ=",
            "value": "Mjk5Mw=="
          },
          {
            "key": "ZXhwaXJhdGlvbi1oZWlnaHQ=",
            "value": "Mjk5Ng=="
          },
          {
            "key": "YXBwX3ZlcnNpb24=",
            "value": "Mw=="
          }
        ]`
	var tags []Tag
	_ = json.Unmarshal([]byte(jsonStr), &tags)
	for _, tag := range tags {
		key, _ := base64.StdEncoding.DecodeString(tag.Key)
		value, _ := base64.StdEncoding.DecodeString(tag.Value)
		fmt.Println(fmt.Sprintf("key=%s value=%s", key, value))
	}
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}
