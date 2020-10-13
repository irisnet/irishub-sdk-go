package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

type Random interface {
	original.Module
	Request(request RandomRequest, baseTx original.BaseTx) (reqID string, err original.Error)
	QueryRandom(reqID string) (ResponseRandom, original.Error)
	QueryRequests(height int64) ([]RequestRandom, original.Error)
	SubscribeRandom(requestID string, callback EventRequestRandomCallback) (original.Subscription, original.Error)
}

type RandomRequest struct {
	BlockInterval uint64            `json:"block_interval"`  // block interval after which the requested random number will be generated
	Oracle        bool              `json:"oracle"`          // oracle method
	ServiceFeeCap original.DecCoins `json:"service_fee_cap"` // service fee cap
	Callback      EventRequestRandomCallback
}

type EventRequestRandomCallback func(reqID, randomNum string, err original.Error)

// ResponseRandom represents a random number with related data
type ResponseRandom struct {
	RequestTxHash string `json:"request_tx_hash"` // the original request tx hash
	Height        int64  `json:"height"`          // the height of the block used to generate the random number
	Value         string `json:"value"`           // the actual random number
}

// RequestRandom represents a request for a random number
type RequestRandom struct {
	Height        int64          `json:"height"`          // the height of the block in which the request tx is included
	Consumer      string         `json:"consumer"`        // the request address
	TxHash        string         `json:"tx_hash"`         // the request tx hash
	Oracle        bool           `json:"oracle"`          // oracle method
	ServiceFeeCap original.Coins `json:"service_fee_cap"` // service fee cap
}
