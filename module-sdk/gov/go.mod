module github.com/irisnet/gov-sdk-go

go 1.16

require (
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/irisnet/core-sdk-go v0.0.0-20210719031639-9c6ece68d908
	github.com/regen-network/cosmos-proto v0.3.1
	google.golang.org/genproto v0.0.0-20201119123407-9b1e624d6bc4
	google.golang.org/grpc v1.37.0
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
