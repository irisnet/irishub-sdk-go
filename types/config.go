package types

import "time"

type ClientConfig struct {
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

	//Transaction broadcast timeout
	Timeout time.Duration

	//log level(trace|debug|info|warn|error|fatal|panic)
	Level string
}
