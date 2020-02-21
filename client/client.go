package client

import (
	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
)

type Client struct {
	types.Bank
}

func New(cfg types.SDKConfig) Client {
	ctx := &types.TxContext{
		Codec:   makeCodec(),
		ChainID: cfg.ChainID,
		Online:  cfg.Online,
		KeyDAO:  cfg.KeyDAO,
		Network: cfg.Network,
		Mode:    cfg.Mode,
		RPC:     net.NewRPCClient(cfg.NodeURI),
	}

	types.SetNetwork(ctx.Network)
	abstClient := abstractClient{ctx}
	return Client{
		Bank: bank.New(abstClient),
	}
}

func makeCodec() types.Codec {
	cdc := types.NewAminoCodec()

	// register msg
	bank.RegisterCodec(cdc)

	return cdc
}
