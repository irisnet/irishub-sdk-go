package client

import (
	"fmt"

	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/distribution"
	"github.com/irisnet/irishub-sdk-go/modules/gov"
	"github.com/irisnet/irishub-sdk-go/modules/oracle"
	"github.com/irisnet/irishub-sdk-go/modules/random"
	"github.com/irisnet/irishub-sdk-go/modules/service"
	"github.com/irisnet/irishub-sdk-go/modules/slashing"
	"github.com/irisnet/irishub-sdk-go/modules/staking"
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type SDKClient struct {
	cdc     sdk.Codec
	modules map[string]sdk.Module

	sdk.WSClient
	sdk.TxManager
}

func NewSDKClient(cfg sdk.SDKConfig) SDKClient {
	cdc := sdk.NewAminoCodec()
	sdk.SetNetwork(cfg.Network)

	abstClient := createAbstractClient(cdc, cfg)
	client := &SDKClient{
		cdc:       cdc,
		modules:   make(map[string]sdk.Module),
		WSClient:  abstClient.TmClient,
		TxManager: abstClient,
	}

	client.registerModule(
		bank.Create(abstClient),
		service.Create(abstClient),
		oracle.Create(abstClient),
		staking.Create(abstClient),
		distribution.Create(abstClient),
		gov.Create(abstClient),
		slashing.Create(abstClient),
		random.Create(abstClient),
	)

	return *client
}

func (s *SDKClient) registerModule(modules ...sdk.Module) {
	for _, m := range modules {
		if _, existed := s.modules[m.Name()]; existed {
			panic(fmt.Sprintf("module[%s] has existed", m.Name()))
		}
		m.RegisterCodec(s.cdc)
		s.modules[m.Name()] = m
	}
	sdk.RegisterCodec(s.cdc)
}

func (s *SDKClient) Bank() rpc.Bank {
	return s.modules[bank.ModuleName].(rpc.Bank)
}

func (s *SDKClient) Distr() rpc.Distribution {
	return s.modules[distribution.ModuleName].(rpc.Distribution)
}

func (s *SDKClient) Service() rpc.Service {
	return s.modules[service.ModuleName].(rpc.Service)
}

func (s *SDKClient) Oracle() rpc.Oracle {
	return s.modules[oracle.ModuleName].(rpc.Oracle)
}

func (s *SDKClient) Staking() rpc.Staking {
	return s.modules[staking.ModuleName].(rpc.Staking)
}

func (s *SDKClient) Gov() rpc.Gov {
	return s.modules[gov.ModuleName].(rpc.Gov)
}

func (s *SDKClient) Slashing() rpc.Slashing {
	return s.modules[slashing.ModuleName].(rpc.Slashing)
}

func (s *SDKClient) Random() rpc.Random {
	return s.modules[random.ModuleName].(rpc.Random)
}
