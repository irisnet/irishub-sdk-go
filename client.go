package sdk

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/original"
	"github.com/irisnet/irishub-sdk-go/modules/original/htlc"
	"github.com/irisnet/irishub-sdk-go/modules/original/params"
	original2 "github.com/irisnet/irishub-sdk-go/types/original"
	"io"

	"github.com/irisnet/irishub-sdk-go/modules/original/asset"
	"github.com/irisnet/irishub-sdk-go/modules/original/bank"
	"github.com/irisnet/irishub-sdk-go/modules/original/distribution"
	"github.com/irisnet/irishub-sdk-go/modules/original/gov"
	"github.com/irisnet/irishub-sdk-go/modules/original/keys"
	"github.com/irisnet/irishub-sdk-go/modules/original/oracle"
	"github.com/irisnet/irishub-sdk-go/modules/original/random"
	"github.com/irisnet/irishub-sdk-go/modules/original/service"
	"github.com/irisnet/irishub-sdk-go/modules/original/slashing"
	"github.com/irisnet/irishub-sdk-go/modules/original/staking"
	"github.com/irisnet/irishub-sdk-go/modules/original/tendermint"
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type Client struct {
	Cdc     original2.Codec
	modules map[string]original2.Module
	logger  *log.Logger

	original2.WSClient
	original2.TxManager
	original2.TokenConvert
}

func NewClient(cfg original2.ClientConfig) Client {
	cdc := original2.NewAminoCodec()
	baseClient := original.NewBaseClient(cdc, cfg)

	client := &Client{
		Cdc:          cdc,
		modules:      make(map[string]original2.Module),
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
		htlc.Create(baseClient),
	)

	return *client
}

func (s *Client) registerModule(modules ...original2.Module) {
	for _, m := range modules {
		if _, existed := s.modules[m.Name()]; existed {
			panic(fmt.Sprintf("module[%s] has existed", m.Name()))
		}
		m.RegisterCodec(s.Cdc)
		s.modules[m.Name()] = m
	}
	original2.RegisterCodec(s.Cdc)
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

func (s *Client) Htlc() rpc.Htlc {
	return s.modules[htlc.ModuleName].(rpc.Htlc)
}

func (s *Client) SetOutput(w io.Writer) {
	s.logger.SetOutput(w)
}
