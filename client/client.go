package client

import (
	"github.com/irisnet/irishub-sdk-go/adapter"
	"github.com/irisnet/irishub-sdk-go/modules/service"
	"github.com/tendermint/tendermint/libs/log"
	"os"

	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
)

type Client struct {
	types.Bank
	types.Service
	types.WSClient
}

func New(cfg types.SDKConfig) Client {
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
	return Client{
		Bank:     bank.New(abstClient),
		Service:  service.New(abstClient),
		WSClient: rpc,
	}
}

func makeCodec() types.Codec {
	cdc := types.NewAminoCodec()

	types.RegisterCodec(cdc)
	// register msg
	bank.RegisterCodec(cdc)
	service.RegisterCodec(cdc)

	return cdc
}
