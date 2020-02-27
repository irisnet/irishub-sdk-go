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
	Fee string

	// Key DAO Implements
	KeyDAO KeyDAO

	// Transaction broadcast Mode
	Mode BroadcastMode

	//
	Online bool

	//
	StoreType StoreType
}

type SDKClient interface {
	Bank
	Service
	Oracle
	WSClient
}

type KeyStore interface {
	GetPrivate() string
	GetAddress() string
}
