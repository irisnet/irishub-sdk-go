module github.com/irisnet/irishub-sdk-go

go 1.13

require (
	github.com/bluele/gcache v0.0.0-20190518031135-bc40bd653833
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/gogo/protobuf v1.3.1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.18.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.14.0
	github.com/tendermint/tendermint v0.33.6
	github.com/tendermint/tm-db v0.5.1
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
	google.golang.org/genproto v0.0.0-20200324203455-a04cca1dde73
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.23.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
