package client

import (
	"os"

	"github.com/irisnet/irishub-sdk-go/modules/oracle"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/irishub-sdk-go/adapter"
	"github.com/irisnet/irishub-sdk-go/modules/service"

	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
)

type client struct {
	types.Bank
	types.Service
	types.Oracle
	types.WSClient
}

func NewSDKClient(cfg types.SDKConfig) types.SDKClient {
	cdc := makeCodec()
	rpc := net.NewRPCClient(cfg.NodeURI, cdc)
	ctx := &types.TxContext{
		Codec:      cdc,
		ChainID:    cfg.ChainID,
		Online:     cfg.Online,
		KeyManager: adapter.NewDAOAdapter(cfg.KeyDAO, cfg.StoreType),
		Network:    cfg.Network,
		Mode:       cfg.Mode,
	}

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	types.SetNetwork(ctx.Network)
	abstClient := abstractClient{
		TxContext: ctx,
		RPC:       rpc,
		logger:    logger,
	}
	return client{
		Bank:     bank.New(abstClient),
		Service:  service.New(abstClient),
		Oracle:   oracle.New(abstClient),
		WSClient: rpc,
	}
}

func makeCodec() types.Codec {
	cdc := types.NewAminoCodec()

	types.RegisterCodec(cdc)
	// register msg
	bank.RegisterCodec(cdc)
	service.RegisterCodec(cdc)
	oracle.RegisterCodec(cdc)

	return cdc
}
