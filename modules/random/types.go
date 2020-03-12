package random

import (
	"errors"

	"github.com/irisnet/irishub-sdk-go/rpc"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub-sdk-go/tools/json"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName   = "random"
	TagRequestID = "request-id"
)

var (
	//_ sdk.Msg = MsgUnjail{}

	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

// MsgRequestRand represents a msg for requesting a random number
type MsgRequestRand struct {
	Consumer      sdk.AccAddress `json:"consumer"`       // request address
	BlockInterval uint64         `json:"block-interval"` // block interval after which the requested random number will be generated
}

// Implements Msg.
func (msg MsgRequestRand) Type() string { return "request_rand" }

// Implements Msg.
func (msg MsgRequestRand) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("the consumer address must be specified")
	}

	return nil
}

// Implements Msg.
func (msg MsgRequestRand) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// Implements Msg.
func (msg MsgRequestRand) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//=======================for query=====================================================
// rand represents a random number with related data
type rand struct {
	RequestTxHash []byte `json:"request_tx_hash"` // the original request tx hash
	Height        int64  `json:"height"`          // the height of the block used to generate the random number
	Value         string `json:"value"`           // the actual random number
}

func (r rand) Convert() interface{} {
	return rpc.RandomInfo{
		RequestTxHash: cmn.HexBytes(r.RequestTxHash).String(),
		Height:        r.Height,
		Value:         r.Value,
	}
}

// ServiceRequest represents a request for a random number
type request struct {
	Height   int64          `json:"height"`   // the height of the block in which the request tx is included
	Consumer sdk.AccAddress `json:"consumer"` // the request address
	TxHash   []byte         `json:"txhash"`   // the request tx hash
}

func (r request) Convert() interface{} {
	return rpc.RequestRandom{
		Height:   r.Height,
		Consumer: r.Consumer.String(),
		TxHash:   cmn.HexBytes(r.TxHash).String(),
	}
}

type requests []request

func (rs requests) Convert() interface{} {
	var requests = make([]rpc.RequestRandom, len(rs))
	for _, r := range rs {
		requests = append(requests, r.Convert().(rpc.RequestRandom))
	}
	return requests
}

func registerCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgRequestRand{}, "irishub/rand/MsgRequestRand")

	cdc.RegisterConcrete(&rand{}, "irishub/rand/Rand")
	cdc.RegisterConcrete(&request{}, "irishub/rand/Request")
}
