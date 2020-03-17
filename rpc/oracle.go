package rpc

import (
	"time"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Oracle interface {
	sdk.Module
	OracleTx
	OracleQuery
}

type OracleTx interface {
	CreateFeed(request FeedCreateRequest) (result sdk.ResultTx, err sdk.Error)
	StartFeed(feedName string, baseTx sdk.BaseTx) (result sdk.ResultTx, err sdk.Error)
	CreateAndStartFeed(request FeedCreateRequest) (result sdk.ResultTx, err sdk.Error)
	PauseFeed(feedName string, baseTx sdk.BaseTx) (result sdk.ResultTx, err sdk.Error)
	EditFeed(request FeedEditRequest) (result sdk.ResultTx, err sdk.Error)
}

type OracleQuery interface {
	QueryFeed(feedName string) (feed FeedContext, err sdk.Error)
	QueryFeeds(state string) (feed []FeedContext, err sdk.Error)
	QueryFeedValue(feedName string) (value []FeedValue, err sdk.Error)
}

// FeedCreateRequest - struct for create a feed
type FeedCreateRequest struct {
	sdk.BaseTx
	FeedName          string       `json:"feed_name"`
	LatestHistory     uint64       `json:"latest_history"`
	Description       string       `json:"description"`
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Input             string       `json:"input"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	AggregateFunc     string       `json:"aggregate_func"`
	ValueJsonPath     string       `json:"value_json_path"`
	ResponseThreshold uint16       `json:"response_threshold"`
}

//______________________________________________________________________

// FeedEditRequest - struct for edit a existed feed
type FeedEditRequest struct {
	sdk.BaseTx
	FeedName          string       `json:"feed_name"`
	Description       string       `json:"description"`
	LatestHistory     uint64       `json:"latest_history"`
	Providers         []string     `json:"providers"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	ResponseThreshold uint16       `json:"response_threshold"`
}

//-----------------------------for query-----------------------------

type Feed struct {
	FeedName         string `json:"feed_name"`
	Description      string `json:"description"`
	AggregateFunc    string `json:"aggregate_func"`
	ValueJsonPath    string `json:"value_json_path"`
	LatestHistory    uint64 `json:"latest_history"`
	RequestContextID string `json:"request_context_id"`
	Creator          string `json:"creator"`
}
type FeedContext struct {
	Feed              Feed      `json:"feed"`
	ServiceName       string    `json:"service_name"`
	Providers         []string  `json:"providers"`
	Input             string    `json:"input"`
	Timeout           int64     `json:"timeout"`
	ServiceFeeCap     sdk.Coins `json:"service_fee_cap"`
	RepeatedFrequency uint64    `json:"repeated_frequency"`
	RepeatedTotal     int64     `json:"repeated_total"`
	ResponseThreshold uint16    `json:"response_threshold"`
	State             string    `json:"state"`
}

type FeedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
