package integrationtest

import (
	"github.com/irisnet/core-sdk-go/bank"
	"github.com/irisnet/core-sdk-go/client"
	"github.com/irisnet/core-sdk-go/common/codec"
	cdctypes "github.com/irisnet/core-sdk-go/common/codec/types"
	cryptocodec "github.com/irisnet/core-sdk-go/common/crypto/codec"
	"github.com/irisnet/core-sdk-go/types"
	txtypes "github.com/irisnet/core-sdk-go/types/tx"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/token-sdk-go"

	"github.com/irisnet/coinswap-sdk-go"

	"github.com/irisnet/staking-sdk-go"

	"github.com/irisnet/service-sdk-go"

	"github.com/irisnet/record-sdk-go"

	"github.com/irisnet/random-sdk-go"

	keys "github.com/irisnet/core-sdk-go/client"

	"github.com/irisnet/oracle-sdk-go"

	"github.com/irisnet/nft-sdk-go"

	"github.com/irisnet/htlc-sdk-go"

	"github.com/irisnet/gov-sdk-go"
)

type Client struct {
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig

	types.BaseClient
	Key     keys.Client
	Bank    bank.Client
	Token   token.Client
	Staking staking.Client
	Gov     gov.Client
	Service service.Client
	Record  record.Client
	Random  random.Client
	NFT     nft.Client
	Oracle  oracle.Client
	HTLC    htlc.Client
	Swap    coinswap.Client
}

func NewClient(cfg types.ClientConfig) Client {
	encodingConfig := makeEncodingConfig()

	// create a instance of baseClient
	baseClient := client.NewBaseClient(cfg, encodingConfig, nil)
	keysClient := keys.NewClient(baseClient)
	bankClient := bank.NewClient(baseClient, encodingConfig.Marshaler)
	tokenClient := token.NewClient(baseClient, encodingConfig.Marshaler)
	stakingClient := staking.NewClient(baseClient, encodingConfig.Marshaler)
	govClient := gov.NewClient(baseClient, encodingConfig.Marshaler)

	serviceClient := service.NewClient(baseClient, encodingConfig.Marshaler)
	recordClient := record.NewClient(baseClient, encodingConfig.Marshaler)
	nftClient := nft.NewClient(baseClient, encodingConfig.Marshaler)
	randomClient := random.NewClient(baseClient, encodingConfig.Marshaler)
	oracleClient := oracle.NewClient(baseClient, encodingConfig.Marshaler)
	htlcClient := htlc.NewClient(baseClient, encodingConfig.Marshaler)
	swapClient := coinswap.NewClient(baseClient, encodingConfig.Marshaler, bankClient.TotalSupply)

	client := &Client{
		logger:         baseClient.Logger(),
		BaseClient:     baseClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,
		Key:            keysClient,
		Bank:           bankClient,
		Token:          tokenClient,
		Staking:        stakingClient,
		Gov:            govClient,
		Service:        serviceClient,
		Record:         recordClient,
		Random:         randomClient,
		NFT:            nftClient,
		Oracle:         oracleClient,
		HTLC:           htlcClient,
		Swap:           swapClient,
	}

	client.RegisterModule(
		bankClient,
		tokenClient,
		stakingClient,
		govClient,
		serviceClient,
		recordClient,
		nftClient,
		randomClient,
		oracleClient,
		htlcClient,
		swapClient,
	)
	return *client
}

func (client Client) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client Client) Codec() *codec.LegacyAmino {
	return client.encodingConfig.Amino
}

func (client Client) AppCodec() codec.Marshaler {
	return client.encodingConfig.Marshaler
}

func (client Client) EncodingConfig() types.EncodingConfig {
	return client.encodingConfig
}

func (client Client) Manager() types.BaseClient {
	return client.BaseClient
}

func (client Client) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
	}
}

func (client Client) Module(name string) types.Module {
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
	cryptocodec.RegisterInterfaces(registry)
}
