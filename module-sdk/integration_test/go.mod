module github.com/irisnet/integration-test

go 1.16

require (
	github.com/irisnet/coinswap-sdk-go v0.1.0
	github.com/irisnet/core-sdk-go v0.1.0
	github.com/irisnet/gov-sdk-go v0.1.0
	github.com/irisnet/htlc-sdk-go v0.1.0
	github.com/irisnet/keys-sdk-go v0.1.0
	github.com/irisnet/nft-sdk-go v0.1.0
	github.com/irisnet/oracle-sdk-go v0.1.0
	github.com/irisnet/random-sdk-go v0.1.0
	github.com/irisnet/record-sdk-go v0.1.0
	github.com/irisnet/service-sdk-go v0.1.0
	github.com/irisnet/staking-sdk-go v0.1.0
	github.com/irisnet/token-sdk-go v0.1.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11

)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/irisnet/coinswap-sdk-go => ../../module-sdk/coinswap
	github.com/irisnet/core-sdk-go => ../../core-sdk
	github.com/irisnet/gov-sdk-go => ../../module-sdk/gov
	github.com/irisnet/htlc-sdk-go => ../../module-sdk/htlc
	github.com/irisnet/keys-sdk-go => ../../module-sdk/keys
	github.com/irisnet/nft-sdk-go => ../../module-sdk/nft
	github.com/irisnet/oracle-sdk-go => ../../module-sdk/oracle
	github.com/irisnet/random-sdk-go => ../../module-sdk/random
	github.com/irisnet/record-sdk-go => ../../module-sdk/record
	github.com/irisnet/service-sdk-go => ../../module-sdk/service
	github.com/irisnet/staking-sdk-go => ../../module-sdk/staking
	github.com/irisnet/token-sdk-go => ../../module-sdk/token
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
)
