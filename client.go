package sdk

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules"
	"github.com/irisnet/irishub-sdk-go/modules/bank"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/irishub-sdk-go/codec"
	cdctypes "github.com/irisnet/irishub-sdk-go/codec/types"
	cryptocodec "github.com/irisnet/irishub-sdk-go/crypto/codec"
	//"github.com/irisnet/irishub-sdk-go/modules"
	//"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/keys"
	//"github.com/irisnet/irishub-sdk-go/modules/nft"
	//"github.com/irisnet/irishub-sdk-go/modules/record"
	//"github.com/irisnet/irishub-sdk-go/modules/service"
	"github.com/irisnet/irishub-sdk-go/types"
	txtypes "github.com/irisnet/irishub-sdk-go/types/tx"
)

type IRISHUBClient struct {
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig

	types.BaseClient
	Bank bank.BankI
	//Token   token.TokenI
	//Record  record.RecordI
	//NFT     nft.NFTI
	//Service service.ServiceI
	Key keys.KeyI
}

func NewIRISHUBClient(cfg types.ClientConfig) IRISHUBClient {
	encodingConfig := makeEncodingConfig()
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//create a instance of baseClient
	baseClient := modules.NewBaseClient(cfg, encodingConfig, nil)

	bankClient := bank.NewClient(baseClient, encodingConfig.Marshaler)
	//tokenClient := token.NewClient(baseClient, encodingConfig.Marshaler)
	keysClient := keys.NewClient(baseClient)
	//recordClient := record.NewClient(baseClient, encodingConfig.Marshaler)
	//nftClient := nft.NewClient(baseClient, encodingConfig.Marshaler)
	//serviceClient := service.NewClient(baseClient, encodingConfig.Marshaler)

	client := &IRISHUBClient{
		logger:     baseClient.Logger(),
		BaseClient: baseClient,
		Bank:       bankClient,
		//Token:          tokenClient,
		Key: keysClient,
		//Record:         recordClient,
		//NFT:            nftClient,
		//Service:        serviceClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,
	}

	client.RegisterModule(
		bankClient,
		//tokenClient,
		//recordClient,
		//nftClient,
		//serviceClient,
	)
	return *client
}

func (client *IRISHUBClient) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client *IRISHUBClient) Codec() *codec.LegacyAmino {
	return client.encodingConfig.Amino
}

func (client *IRISHUBClient) AppCodec() codec.Marshaler {
	return client.encodingConfig.Marshaler
}

func (client *IRISHUBClient) Manager() types.BaseClient {
	return client.BaseClient
}

func (client *IRISHUBClient) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		_, ok := client.moduleManager[m.Name()]
		if ok {
			panic(fmt.Sprintf("%s has register", m.Name()))
		}

		// m.RegisterCodec(client.encodingConfig.Amino)
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
		client.moduleManager[m.Name()] = m
	}
}

func (client *IRISHUBClient) Module(name string) types.Module {
	return client.moduleManager[name]
}

func makeEncodingConfig() types.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := txtypes.NewTxConfig(marshaler, txtypes.DefaultSignModes)

	encodingConfig := types.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	RegisterLegacyAminoCodec(encodingConfig.Amino)
	RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterInterface((*types.Tx)(nil), nil)
	cryptocodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.v1beta1.Msg", (*types.Msg)(nil))
	txtypes.RegisterInterfaces(registry)
}
