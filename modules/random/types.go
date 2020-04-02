package random

import (
	"errors"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/rpc"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	ModuleName   = "random"
	tagRequestID = "request-id"
)

var (
	_ sdk.Msg = MsgRequestRand{}

	cdc = sdk.NewAminoCodec()

	tagRand = func(reqID string) string {
		return fmt.Sprintf("rand.%s", reqID)
	}
)

func init() {
	registerCodec(cdc)
}

// MsgRequestRand represents a msg for requesting a random number
type MsgRequestRand struct {
	Consumer      sdk.AccAddress `json:"consumer"`        // request address
	BlockInterval uint64         `json:"block_interval"`  // block interval after which the requested random number will be generated
	Oracle        bool           `json:"oracle"`          // oracle method
	ServiceFeeCap sdk.Coins      `json:"service_fee_cap"` // service fee cap
}

func (msg MsgRequestRand) Route() string { return ModuleName }

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
	return rpc.ResponseRandom{
		RequestTxHash: cmn.HexBytes(r.RequestTxHash).String(),
		Height:        r.Height,
		Value:         r.Value,
	}
}

// ServiceRequest represents a request for a random number
type request struct {
	Height        int64          `json:"height"`          // the height of the block in which the request tx is included
	Consumer      sdk.AccAddress `json:"consumer"`        // the request address
	TxHash        []byte         `json:"txhash"`          // the request tx hash
	Oracle        bool           `json:"oracle"`          // oracle method
	ServiceFeeCap sdk.Coins      `json:"service_fee_cap"` // service fee cap
}

func (r request) Convert() interface{} {
	return rpc.RequestRandom{
		Height:        r.Height,
		Consumer:      r.Consumer.String(),
		TxHash:        cmn.HexBytes(r.TxHash).String(),
		Oracle:        r.Oracle,
		ServiceFeeCap: r.ServiceFeeCap,
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
