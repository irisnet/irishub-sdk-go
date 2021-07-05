module github.com/irisnet/oracle-sdk-go

go 1.16

require (
    github.com/irisnet/core-sdk-go v0.1.0
        github.com/irisnet/service-sdk-go v0.1.0
    github.com/gogo/protobuf v1.3.3
)

replace (
github.com/irisnet/core-sdk-go => ../../../core-sdk
   github.com/irisnet/service-sdk-go => ../../../module-sdk/service
github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)