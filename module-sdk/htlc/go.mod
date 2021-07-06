module github.com/irisnet/htlc-sdk-go

go 1.16

require (
    github.com/irisnet/core-sdk-go v0.1.0
    github.com/gogo/protobuf v1.3.3
)

replace (
github.com/irisnet/core-sdk-go => github.com/Nicke-lucky/core-sdk-go 3f1da92ff70f692ab20a6f7326429de1f5c69a34
github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)