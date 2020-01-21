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
}

type KeyDAO interface {
	/**
	 * Save the keystore to app, throws error if the save fails.
	 *
	 * @param keystore The keystore object
	 */
	Write(name string, keystore KeyStore) error

	/**
	 * Get the keystore by address
	 *
	 * @param name Name of the key
	 * @returns The keystore object
	 */
	Read(name string) KeyStore

	/**
	 * Delete keystore by address
	 * @param name Name of the key
	 */
	Delete(name string) error
}

type KeyStore interface {
	GetPrivate() string
	GetAddress() string
}
