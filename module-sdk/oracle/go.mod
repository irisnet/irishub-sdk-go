module github.com/irisnet/oracle-sdk-go

go 1.16

require (
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/irisnet/core-sdk-go v0.1.0
	github.com/irisnet/service-sdk-go v0.1.0
	google.golang.org/genproto v0.0.0-20201119123407-9b1e624d6bc4
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/irisnet/core-sdk-go => github.com/Nicke-lucky/core-sdk-go v0.0.0-20210706063401-ba48b2920add
	github.com/irisnet/service-sdk-go => ../../../module-sdk/service
)
