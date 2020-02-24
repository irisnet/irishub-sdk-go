package client

import (
	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
)

type Client struct {
	types.Bank
	types.WSClient
}

func New(cfg types.SDKConfig) Client {
	cdc := makeCodec()
	rpc := net.NewRPCClient(cfg.NodeURI, cdc)
	ctx := &types.TxContext{
		Codec:   cdc,
		ChainID: cfg.ChainID,
		Online:  cfg.Online,
		KeyDAO:  cfg.KeyDAO,
		Network: cfg.Network,
		Mode:    cfg.Mode,
		RPC:     rpc,
	}

	types.SetNetwork(ctx.Network)
	abstClient := abstractClient{
		ctx,
	}
	return Client{
		Bank:     bank.New(abstClient),
		WSClient: rpc,
	}
}

func makeCodec() types.Codec {
	cdc := types.NewAminoCodec()

	types.RegisterCodec(cdc)
	// register msg
	bank.RegisterCodec(cdc)

	return cdc
}
