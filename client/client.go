package client

import (
	"fmt"

	"github.com/irisnet/irishub-sdk-go/modules/random"
	"github.com/irisnet/irishub-sdk-go/modules/slashing"

	"github.com/irisnet/irishub-sdk-go/adapter"
	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/distribution"
	"github.com/irisnet/irishub-sdk-go/modules/gov"
	"github.com/irisnet/irishub-sdk-go/modules/oracle"
	"github.com/irisnet/irishub-sdk-go/modules/service"
	"github.com/irisnet/irishub-sdk-go/modules/staking"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type SDKClient struct {
	cdc     sdk.Codec
	modules map[string]sdk.Module

	sdk.WSClient
}

func NewSDKClient(cfg sdk.SDKConfig) SDKClient {
	client := &SDKClient{
		cdc:     sdk.NewAminoCodec(),
		modules: make(map[string]sdk.Module),
	}

	rpc := net.NewRPCClient(cfg.NodeURI, client.cdc)
	ctx := &sdk.TxContext{
		Codec:      client.cdc,
		ChainID:    cfg.ChainID,
		Online:     cfg.Online,
		KeyManager: adapter.NewDAOAdapter(cfg.KeyDAO, cfg.StoreType),
		Network:    cfg.Network,
		Mode:       cfg.Mode,
	}

	sdk.SetNetwork(ctx.Network)
	abstClient := abstractClient{
		TxContext: ctx,
		RPC:       rpc,
		logger:    log.NewLogger(cfg.Level).With("AbstractClient"),
	}

	client.registerModule(
		bank.New(abstClient),
		service.New(abstClient),
		oracle.New(abstClient),
		staking.New(abstClient),
		distribution.New(abstClient),
		gov.New(abstClient),
		slashing.New(abstClient),
		random.New(abstClient),
	)

	return *client
}

func (s *SDKClient) registerModule(modules ...sdk.Module) {
	s.modules = make(map[string]sdk.Module, len(modules))
	for _, m := range modules {
		if _, existed := s.modules[m.Name()]; existed {
			panic(fmt.Sprintf("module[%s] has existed", m.Name()))
		}
		m.RegisterCodec(s.cdc)
		s.modules[m.Name()] = m
	}
	sdk.RegisterCodec(s.cdc)
}

func (s *SDKClient) Bank() sdk.Bank {
	return s.modules[bank.ModuleName].(sdk.Bank)
}

func (s *SDKClient) Distr() sdk.Distribution {
	return s.modules[distribution.ModuleName].(sdk.Distribution)
}

func (s *SDKClient) Service() sdk.Service {
	return s.modules[service.ModuleName].(sdk.Service)
}

func (s *SDKClient) Oracle() sdk.Oracle {
	return s.modules[oracle.ModuleName].(sdk.Oracle)
}

func (s *SDKClient) Staking() sdk.Staking {
	return s.modules[staking.ModuleName].(sdk.Staking)
}

func (s *SDKClient) Gov() sdk.Gov {
	return s.modules[gov.ModuleName].(sdk.Gov)
}

func (s *SDKClient) Slashing() sdk.Slashing {
	return s.modules[slashing.ModuleName].(sdk.Slashing)
}

func (s *SDKClient) Random() sdk.Random {
	return s.modules[random.ModuleName].(sdk.Random)
}
