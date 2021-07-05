package oracle

import (
	"time"

	sdk "github.com/irisnet/core-sdk-go/types"
)

// expose Oracle module api for user
type Client interface {
	sdk.Module

	CreateFeed(request CreateFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	StartFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	PauseFeed(FeedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	EditFeed(request EditFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryFeed(feedName string) (QueryFeedResp, sdk.Error)
	QueryFeeds(state string) ([]QueryFeedResp, sdk.Error)
	QueryFeedValue(feedName string) ([]QueryFeedValueResp, sdk.Error)
}

type CreateFeedRequest struct {
	FeedName          string       `json:"feed_name"`
	LatestHistory     uint64       `json:"latest_history"`
	Description       string       `json:"description"`
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Input             string       `json:"input"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	AggregateFunc     string       `json:"aggregate_func"`
	ValueJsonPath     string       `json:"value_json_path"`
	ResponseThreshold uint32       `json:"response_threshold"`
}

type EditFeedRequest struct {
	FeedName          string       `json:"feed_name"`
	Description       string       `json:"description"`
	LatestHistory     uint64       `json:"latest_history"`
	Providers         []string     `json:"providers"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	ResponseThreshold uint32       `json:"response_threshold"`
}

type QueryFeedResp struct {
	Feed struct {
		FeedName         string `json:"feed_name"`
		Description      string `json:"description"`
		AggregateFunc    string `json:"aggregate_func"`
		ValueJsonPath    string `json:"value_json_path"`
		LatestHistory    uint64 `json:"latest_history"`
		RequestContextID string `json:"request_context_id"`
		Creator          string `json:"creator"`
	} `json:"feed"`
	ServiceName       string    `json:"service_name"`
	Providers         []string  `json:"providers"`
	Input             string    `json:"input"`
	Timeout           int64     `json:"timeout"`
	ServiceFeeCap     sdk.Coins `json:"service_fee_cap"`
	RepeatedFrequency uint64    `json:"repeated_frequency"`
	ResponseThreshold uint32    `json:"response_threshold"`
	State             int32     `json:"state"`
}

type QueryFeedValueResp struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
