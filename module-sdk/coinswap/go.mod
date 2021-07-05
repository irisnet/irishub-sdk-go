module github.com/irisnet/coinswap-sdk-go

go 1.16

require (
github.com/gogo/protobuf v1.3.3
github.com/irisnet/core-sdk-go v0.1.0
)

replace (
github.com/irisnet/core-sdk-go => ../../../core-sdk
github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)