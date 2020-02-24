package client

import (
	"errors"

	"github.com/irisnet/irishub-sdk-go/crypto"
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
		KeyManager: keyManager{
			cfg.KeyDAO,
		},
		Network: cfg.Network,
		Mode:    cfg.Mode,
	}

	types.SetNetwork(ctx.Network)
	abstClient := abstractClient{
		TxContext: ctx,
		RPC:       rpc,
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

type keyManager struct {
	types.KeyDAO
}

func (manager keyManager) Sign(name, password string, data []byte) (signature types.Signature, err error) {
	store := manager.Read(name)

	var mm crypto.KeyManager
	switch store := store.(type) {
	case types.KeyInfo:
		mm, err = crypto.NewPrivateKeyManager(store.PrivKey)
		if err != nil {
			return signature, err
		}
	case types.KeystoreInfo:
		mm, err = crypto.NewKeyStoreKeyManager(store.KeystoreJSON, password)
		if err != nil {
			return signature, err
		}
	}
	signByte, err := mm.Sign(data)

	return types.Signature{
		PubKey:    mm.GetPrivKey().PubKey(),
		Signature: signByte,
	}, nil
}

func (manager keyManager) QueryAddress(name, password string) (addr types.AccAddress, err error) {
	store := manager.Read(name)

	var mm crypto.KeyManager
	switch store := store.(type) {
	case types.KeyInfo:
		mm, err = crypto.NewPrivateKeyManager(store.PrivKey)
		if err != nil {
			return addr, err
		}
		return types.AccAddressFromBech32(store.Address)
	case types.KeystoreInfo:
		mm, err = crypto.NewKeyStoreKeyManager(store.KeystoreJSON, password)
		if err != nil {
			return addr, err
		}
		accAddr := types.AccAddress(mm.GetPrivKey().PubKey().Address())
		return accAddr, nil
	}
	return addr, errors.New("invalid StoreType")
}
