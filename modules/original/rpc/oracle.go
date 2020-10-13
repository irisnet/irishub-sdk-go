package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
	"time"
)

type Oracle interface {
	original.Module
	OracleTx
	OracleQuery
}

type OracleTx interface {
	CreateFeed(request FeedCreateRequest) (result original.ResultTx, err original.Error)
	StartFeed(feedName string, baseTx original.BaseTx) (result original.ResultTx, err original.Error)
	CreateAndStartFeed(request FeedCreateRequest) (result original.ResultTx, err original.Error)
	PauseFeed(feedName string, baseTx original.BaseTx) (result original.ResultTx, err original.Error)
	EditFeed(request FeedEditRequest) (result original.ResultTx, err original.Error)
	SubscribeFeedValue(feedName string, callback func(value FeedValue)) original.Error
}

type OracleQuery interface {
	QueryFeed(feedName string) (feed FeedContext, err original.Error)
	QueryFeeds(state string) (feed []FeedContext, err original.Error)
	QueryFeedValue(feedName string) (value []FeedValue, err original.Error)
}

// FeedCreateRequest - struct for create a feed
type FeedCreateRequest struct {
	original.BaseTx
	FeedName          string            `json:"feed_name"`
	LatestHistory     uint64            `json:"latest_history"`
	Description       string            `json:"description"`
	ServiceName       string            `json:"service_name"`
	Providers         []string          `json:"providers"`
	Input             string            `json:"input"`
	Timeout           int64             `json:"timeout"`
	ServiceFeeCap     original.DecCoins `json:"service_fee_cap"`
	RepeatedFrequency uint64            `json:"repeated_frequency"`
	AggregateFunc     string            `json:"aggregate_func"`
	ValueJsonPath     string            `json:"value_json_path"`
	ResponseThreshold uint16            `json:"response_threshold"`
}

//______________________________________________________________________

// FeedEditRequest - struct for edit a existed feed
type FeedEditRequest struct {
	original.BaseTx
	FeedName          string            `json:"feed_name"`
	Description       string            `json:"description"`
	LatestHistory     uint64            `json:"latest_history"`
	Providers         []string          `json:"providers"`
	Timeout           int64             `json:"timeout"`
	ServiceFeeCap     original.DecCoins `json:"service_fee_cap"`
	RepeatedFrequency uint64            `json:"repeated_frequency"`
	ResponseThreshold uint16            `json:"response_threshold"`
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
	Feed              Feed           `json:"feed"`
	ServiceName       string         `json:"service_name"`
	Providers         []string       `json:"providers"`
	Input             string         `json:"input"`
	Timeout           int64          `json:"timeout"`
	ServiceFeeCap     original.Coins `json:"service_fee_cap"`
	RepeatedFrequency uint64         `json:"repeated_frequency"`
	ResponseThreshold uint16         `json:"response_threshold"`
	State             string         `json:"state"`
}

type FeedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
