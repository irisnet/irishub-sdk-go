package sdk

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/params"
	"io"

	"github.com/irisnet/irishub-sdk-go/modules"
	"github.com/irisnet/irishub-sdk-go/modules/asset"
	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/distribution"
	"github.com/irisnet/irishub-sdk-go/modules/gov"
	"github.com/irisnet/irishub-sdk-go/modules/keys"
	"github.com/irisnet/irishub-sdk-go/modules/oracle"
	"github.com/irisnet/irishub-sdk-go/modules/random"
	"github.com/irisnet/irishub-sdk-go/modules/service"
	"github.com/irisnet/irishub-sdk-go/modules/slashing"
	"github.com/irisnet/irishub-sdk-go/modules/staking"
	"github.com/irisnet/irishub-sdk-go/modules/tendermint"
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type Client struct {
	Cdc     sdk.Codec
	modules map[string]sdk.Module
	logger  *log.Logger

	sdk.WSClient
	sdk.TxManager
	sdk.TokenConvert
}

func NewClient(cfg sdk.ClientConfig) Client {
	cdc := sdk.NewAminoCodec()
	baseClient := modules.NewBaseClient(cdc, cfg)

	client := &Client{
		Cdc:          cdc,
		modules:      make(map[string]sdk.Module),
		logger:       baseClient.Logger(),
		WSClient:     baseClient.TmClient,
		TxManager:    baseClient,
		TokenConvert: baseClient,
	}

	client.registerModule(
		bank.Create(baseClient),
		service.Create(baseClient),
		oracle.Create(baseClient),
		staking.Create(baseClient),
		distribution.Create(baseClient),
		gov.Create(baseClient),
		slashing.Create(baseClient),
		random.Create(baseClient),
		keys.Create(baseClient.KeyManager),
		asset.Create(baseClient),
		tendermint.Create(baseClient, cdc),
		params.Create(baseClient),
	)

	return *client
}

func (s *Client) registerModule(modules ...sdk.Module) {
	for _, m := range modules {
		if _, existed := s.modules[m.Name()]; existed {
			panic(fmt.Sprintf("module[%s] has existed", m.Name()))
		}
		m.RegisterCodec(s.Cdc)
		s.modules[m.Name()] = m
	}
	sdk.RegisterCodec(s.Cdc)
}

func (s *Client) Bank() rpc.Bank {
	return s.modules[bank.ModuleName].(rpc.Bank)
}

func (s *Client) Distr() rpc.Distribution {
	return s.modules[distribution.ModuleName].(rpc.Distribution)
}

func (s *Client) Service() rpc.Service {
	return s.modules[service.ModuleName].(rpc.Service)
}

func (s *Client) Oracle() rpc.Oracle {
	return s.modules[oracle.ModuleName].(rpc.Oracle)
}

func (s *Client) Staking() rpc.Staking {
	return s.modules[staking.ModuleName].(rpc.Staking)
}

func (s *Client) Gov() rpc.Gov {
	return s.modules[gov.ModuleName].(rpc.Gov)
}

func (s *Client) Slashing() rpc.Slashing {
	return s.modules[slashing.ModuleName].(rpc.Slashing)
}

func (s *Client) Random() rpc.Random {
	return s.modules[random.ModuleName].(rpc.Random)
}

func (s *Client) Keys() rpc.Keys {
	return s.modules[keys.ModuleName].(rpc.Keys)
}

func (s *Client) Asset() rpc.Asset {
	return s.modules[asset.ModuleName].(rpc.Asset)
}

func (s *Client) Tendermint() rpc.Tendermint {
	return s.modules[tendermint.ModuleName].(rpc.Tendermint)
}

func (s *Client) Params() rpc.Params {
	return s.modules[params.ModuleName].(rpc.Params)
}

func (s *Client) SetOutput(w io.Writer) {
	s.logger.SetOutput(w)
}
