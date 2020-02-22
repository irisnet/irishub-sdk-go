package types

const (
	Testnet Network = "testnet"
	Mainnet Network = "mainnet"
)

type Network string

// Can be configured through environment variables
var (
	defaultNetwork = Testnet
)

var (
	testnetEnv = &AddrPrefixCfg{
		bech32AddressPrefix: map[string]string{
			"account_addr":   "faa",
			"validator_addr": "fva",
			"consensus_addr": "fca",
			"account_pub":    "fap",
			"validator_pub":  "fvp",
			"consensus_pub":  "fcp",
		},
	}
	mainnetEnv = &AddrPrefixCfg{
		bech32AddressPrefix: map[string]string{
			"account_addr":   "iaa",
			"validator_addr": "iva",
			"consensus_addr": "ica",
			"account_pub":    "iap",
			"validator_pub":  "ivp",
			"consensus_pub":  "icp",
		},
	}
)

type AddrPrefixCfg struct {
	bech32AddressPrefix map[string]string
}

func SetNetwork(network Network) {
	defaultNetwork = network
}

// GetAddrPrefixCfg returns the config instance for the corresponding Network type
func GetAddrPrefixCfg() *AddrPrefixCfg {
	if defaultNetwork == Mainnet {
		return mainnetEnv
	}
	return testnetEnv
}

// GetBech32AccountAddrPrefix returns the Bech32 prefix for account address
func (config *AddrPrefixCfg) GetBech32AccountAddrPrefix() string {
	return config.bech32AddressPrefix["account_addr"]
}

// GetBech32ValidatorAddrPrefix returns the Bech32 prefix for validator address
func (config *AddrPrefixCfg) GetBech32ValidatorAddrPrefix() string {
	return config.bech32AddressPrefix["validator_addr"]
}

// GetBech32ConsensusAddrPrefix returns the Bech32 prefix for consensus node address
func (config *AddrPrefixCfg) GetBech32ConsensusAddrPrefix() string {
	return config.bech32AddressPrefix["consensus_addr"]
}

// GetBech32AccountPubPrefix returns the Bech32 prefix for account public key
func (config *AddrPrefixCfg) GetBech32AccountPubPrefix() string {
	return config.bech32AddressPrefix["account_pub"]
}

// GetBech32ValidatorPubPrefix returns the Bech32 prefix for validator public key
func (config *AddrPrefixCfg) GetBech32ValidatorPubPrefix() string {
	return config.bech32AddressPrefix["validator_pub"]
}

// GetBech32ConsensusPubPrefix returns the Bech32 prefix for consensus node public key
func (config *AddrPrefixCfg) GetBech32ConsensusPubPrefix() string {
	return config.bech32AddressPrefix["consensus_pub"]
}
