package types

type SDKConfig struct {
	// IRISHub node rpc address
	NodeURI string

	// IRISHub Network type, mainnet / testnet
	Network Network

	// IRISHub chain-id
	ChainID string

	// Default Gas limit
	Gas uint64

	// Default Fee amount of iris-atto
	Fee DecCoins

	// Key DAO Implements
	KeyDAO KeyDAO

	// Transaction broadcast Mode
	Mode BroadcastMode

	//
	StoreType StoreType

	//log level(trace|debug|info|warn|error|fatal|panic)
	Level string
}

//type SDKClient interface {
//	Bank
//	Service
//	Oracle
//	WSClient
//	Staking
//	Distribution
//}

type KeyStore interface {
	GetPrivate() string
	GetAddress() string
}
