package sdk

import (
	"github.com/irisnet/irishub-sdk-go/bank"
	"github.com/irisnet/irishub-sdk-go/client"
	commoncodec "github.com/irisnet/irishub-sdk-go/common/codec"
	cryptotypes "github.com/irisnet/irishub-sdk-go/common/codec/types"
	commoncryptocodec "github.com/irisnet/irishub-sdk-go/common/crypto/codec"

	"github.com/irisnet/irishub-sdk-go/types"
	//"github.com/irisnet/irishub-sdk-go/types/token"
	txtypes "github.com/irisnet/irishub-sdk-go/types/tx"
	"github.com/tendermint/tendermint/libs/log"
)

type IRISHUBClient struct {
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig
	types.BaseClient
	Bank bank.Client
}

func NewIRISHUBClient(cfg types.ClientConfig) IRISHUBClient {
	encodingConfig := makeEncodingConfig()

	// create a instance of baseClient
	baseClient := client.NewBaseClient(cfg, encodingConfig, nil)
	bankClient := bank.NewClient(baseClient, encodingConfig.Marshaler)
	client := &IRISHUBClient{
		logger:         baseClient.Logger(),
		BaseClient:     baseClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,

		Bank: bankClient,
	}
	client.RegisterModule(
		bankClient,
	)
	return *client
}

func (client *IRISHUBClient) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client *IRISHUBClient) Codec() *commoncodec.LegacyAmino {
	return client.encodingConfig.Amino
}

func (client *IRISHUBClient) AppCodec() commoncodec.Marshaler {
	return client.encodingConfig.Marshaler
}

func (client *IRISHUBClient) EncodingConfig() types.EncodingConfig {
	return client.encodingConfig
}

func (client *IRISHUBClient) Manager() types.BaseClient {
	return client.BaseClient
}

func (client *IRISHUBClient) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
	}
}

func (client *IRISHUBClient) Module(name string) types.Module {
	return client.moduleManager[name]
}

func makeEncodingConfig() types.EncodingConfig {
	amino := commoncodec.NewLegacyAmino()
	interfaceRegistry := cryptotypes.NewInterfaceRegistry()
	marshaler := commoncodec.NewProtoCodec(interfaceRegistry)
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
func RegisterLegacyAminoCodec(cdc *commoncodec.LegacyAmino) {
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterInterface((*types.Tx)(nil), nil)
	commoncryptocodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry cryptotypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.v1beta1.Msg", (*types.Msg)(nil))
	txtypes.RegisterInterfaces(registry)
	commoncryptocodec.RegisterInterfaces(registry)
}
