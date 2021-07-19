module github.com/irisnet/integration-test

go 1.16

require (
	github.com/irisnet/coinswap-sdk-go v0.1.0
	github.com/irisnet/core-sdk-go v0.0.0-20210719031639-9c6ece68d908
	github.com/irisnet/gov-sdk-go v0.1.0
	github.com/irisnet/htlc-sdk-go v0.1.0
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
	github.com/irisnet/coinswap-sdk-go => ../coinswap
	github.com/irisnet/gov-sdk-go => ../gov
	github.com/irisnet/htlc-sdk-go => ../htlc
	github.com/irisnet/nft-sdk-go => ../nft
	github.com/irisnet/oracle-sdk-go => ../oracle
	github.com/irisnet/random-sdk-go => ../random
	github.com/irisnet/record-sdk-go => ../record
	github.com/irisnet/service-sdk-go => ../service
	github.com/irisnet/staking-sdk-go => ../staking
	github.com/irisnet/token-sdk-go => ../token
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
)
