package rpc

import (
	"time"

	"github.com/irisnet/irishub-sdk-go/types"
)

type Oracle interface {
	types.Module
	OracleTx
	OracleQuery
}

type OracleTx interface {
	CreateFeed(request FeedCreateRequest) (result types.Result, err error)
	StartFeed(feedName string, baseTx types.BaseTx) (result types.Result, err error)
	CreateAndStartFeed(request FeedCreateRequest) (result types.Result, err error)
	PauseFeed(feedName string, baseTx types.BaseTx) (result types.Result, err error)
	EditFeed(request FeedEditRequest) (result types.Result, err error)
}

type OracleQuery interface {
	QueryFeed(feedName string) (feed FeedContext, err error)
	QueryFeeds(state string) (feed []FeedContext, err error)
	QueryFeedValue(feedName string) (value []FeedValue, err error)
}

// FeedCreateRequest - struct for create a feed
type FeedCreateRequest struct {
	types.BaseTx
	FeedName          string      `json:"feed_name"`
	LatestHistory     uint64      `json:"latest_history"`
	Description       string      `json:"description"`
	ServiceName       string      `json:"service_name"`
	Providers         []string    `json:"providers"`
	Input             string      `json:"input"`
	Timeout           int64       `json:"timeout"`
	ServiceFeeCap     types.Coins `json:"service_fee_cap"`
	RepeatedFrequency uint64      `json:"repeated_frequency"`
	RepeatedTotal     int64       `json:"repeated_total"`
	AggregateFunc     string      `json:"aggregate_func"`
	ValueJsonPath     string      `json:"value_json_path"`
	ResponseThreshold uint16      `json:"response_threshold"`
}

//______________________________________________________________________

// FeedEditRequest - struct for edit a existed feed
type FeedEditRequest struct {
	types.BaseTx
	FeedName          string      `json:"feed_name"`
	Description       string      `json:"description"`
	LatestHistory     uint64      `json:"latest_history"`
	Providers         []string    `json:"providers"`
	Timeout           int64       `json:"timeout"`
	ServiceFeeCap     types.Coins `json:"service_fee_cap"`
	RepeatedFrequency uint64      `json:"repeated_frequency"`
	RepeatedTotal     int64       `json:"repeated_total"`
	ResponseThreshold uint16      `json:"response_threshold"`
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
	Feed              Feed        `json:"feed"`
	ServiceName       string      `json:"service_name"`
	Providers         []string    `json:"providers"`
	Input             string      `json:"input"`
	Timeout           int64       `json:"timeout"`
	ServiceFeeCap     types.Coins `json:"service_fee_cap"`
	RepeatedFrequency uint64      `json:"repeated_frequency"`
	RepeatedTotal     int64       `json:"repeated_total"`
	ResponseThreshold uint16      `json:"response_threshold"`
	State             string      `json:"state"`
}

type FeedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
