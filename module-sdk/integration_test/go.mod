module  github.com/irisnet/integration-test

go 1.16

require (
github.com/irisnet/module-sdk v0.1.1
github.com/gogo/protobuf v1.3.3
github.com/irisnet/coinswap-sdk-go v0.1.0
github.com/irisnet/core-sdk-go v0.1.0
github.com/irisnet/gov-sdk-go v0.1.0
github.com/irisnet/htlc-sdk-go v0.1.0
github.com/irisnet/nft-sdk-go v0.1.0
github.com/irisnet/random-sdk-go v0.1.0
github.com/irisnet/service-sdk-go v0.1.0
github.com/irisnet/record-sdk-go v0.1.0
github.com/irisnet/staking-sdk-go v0.1.0
github.com/irisnet/token-sdk-go v0.1.0
github.com/irisnet/keys-sdk-go v0.1.0
github.com/irisnet/oracle-sdk-go v0.1.0

)

replace (
github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
github.com/irisnet/core-sdk-go => github.com/Nicke-lucky/core-sdk-go 3f1da92ff70f692ab20a6f7326429de1f5c69a34
github.com/irisnet/module-sdk =>  github.com/Nicke-lucky/module-sdk-go 8c6093a2351a9a83e58715becd7dbd2411cc06fa
github.com/irisnet/gov-sdk-go => ../../module-sdk/gov
github.com/irisnet/htlc-sdk-go => ../../module-sdk/htlc
github.com/irisnet/nft-sdk-go => ../../module-sdk/nft
github.com/irisnet/random-sdk-go => ../../module-sdk/random
github.com/irisnet/service-sdk-go => ../../module-sdk/service
github.com/irisnet/record-sdk-go => ../../module-sdk/record
github.com/irisnet/staking-sdk-go => ../../module-sdk/staking
github.com/irisnet/token-sdk-go => ../../module-sdk/token
github.com/irisnet/keys-sdk-go  =>  ../../module-sdk/keys
github.com/irisnet/oracle-sdk-go =>  ../../module-sdk/oracle
github.com/irisnet/coinswap-sdk-go =>  ../../module-sdk/coinswap
)