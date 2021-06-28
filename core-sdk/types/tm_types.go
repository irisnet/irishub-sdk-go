package types

import (
	cryptoAmino "github.com/irisnet/irishub-sdk-go/common/crypto/codec"
	"github.com/irisnet/irishub-sdk-go/types/kv"
	"github.com/tendermint/tendermint/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

type (
	HexBytes      = tmbytes.HexBytes
	ABCIClient    = tmclient.ABCIClient
	SignClient    = tmclient.SignClient
	StatusClient  = tmclient.StatusClient
	NetworkClient = tmclient.NetworkClient
	Header        = tmtypes.Header
	Pair          = kv.Pair

	TmPubKey = crypto.PubKey
)

var (
	PubKeyFromBytes = cryptoAmino.PubKeyFromBytes
)
