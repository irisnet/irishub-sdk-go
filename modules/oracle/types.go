package oracle

import (
	"errors"
	"strings"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/json"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "oracle"
)

var (
	_ sdk.Msg = MsgCreateFeed{}
	_ sdk.Msg = MsgStartFeed{}
	_ sdk.Msg = MsgPauseFeed{}
	_ sdk.Msg = MsgEditFeed{}

	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

//______________________________________________________________________

// MsgCreateFeed - struct for create a feed
type MsgCreateFeed struct {
	FeedName          string           `json:"feed_name"`
	LatestHistory     uint64           `json:"latest_history"`
	Description       string           `json:"description"`
	Creator           sdk.AccAddress   `json:"creator"`
	ServiceName       string           `json:"service_name"`
	Providers         []sdk.AccAddress `json:"providers"`
	Input             string           `json:"input"`
	Timeout           int64            `json:"timeout"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	AggregateFunc     string           `json:"aggregate_func"`
	ValueJsonPath     string           `json:"value_json_path"`
	ResponseThreshold uint16           `json:"response_threshold"`
}

// Type implements Msg.
func (msg MsgCreateFeed) Type() string {
	return "create_feed"
}

// ValidateBasic implements Msg.
func (msg MsgCreateFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}

	serviceName := strings.TrimSpace(msg.ServiceName)
	if len(serviceName) == 0 {
		return errors.New("serviceName missed")
	}

	if len(msg.Providers) == 0 {
		return errors.New("providers missed")
	}

	aggregateFunc := strings.TrimSpace(msg.AggregateFunc)
	if len(aggregateFunc) == 0 {
		return errors.New("aggregateFunc missed")
	}

	valueJsonPath := strings.TrimSpace(msg.ValueJsonPath)
	if len(valueJsonPath) == 0 {
		return errors.New("valueJsonPath missed")
	}

	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgCreateFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgCreateFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgStartFeed - struct for start a feed
type MsgStartFeed struct {
	FeedName string         `json:"feed_name"`
	Creator  sdk.AccAddress `json:"creator"`
}

// Type implements Msg.
func (msg MsgStartFeed) Type() string {
	return "start_feed"
}

// ValidateBasic implements Msg.
func (msg MsgStartFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}
	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgStartFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgStartFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgPauseFeed - struct for stop a started feed
type MsgPauseFeed struct {
	FeedName string         `json:"feed_name"`
	Creator  sdk.AccAddress `json:"creator"`
}

// Type implements Msg.
func (msg MsgPauseFeed) Type() string {
	return "pause_feed"
}

// ValidateBasic implements Msg.
func (msg MsgPauseFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}
	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgPauseFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgPauseFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgEditFeed - struct for edit a existed feed
type MsgEditFeed struct {
	FeedName          string           `json:"feed_name"`
	Description       string           `json:"description"`
	LatestHistory     uint64           `json:"latest_history"`
	Providers         []sdk.AccAddress `json:"providers"`
	Timeout           int64            `json:"timeout"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	ResponseThreshold uint16           `json:"response_threshold"`
	Creator           sdk.AccAddress   `json:"creator"`
}

// Type implements Msg.
func (msg MsgEditFeed) Type() string {
	return "edit_feed"
}

// ValidateBasic implements Msg.
func (msg MsgEditFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}

	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgEditFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgEditFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

//-------------------------------for query--------------------------
type feed struct {
	FeedName         string         `json:"feed_name"`
	Description      string         `json:"description"`
	AggregateFunc    string         `json:"aggregate_func"`
	ValueJsonPath    string         `json:"value_json_path"`
	LatestHistory    uint64         `json:"latest_history"`
	RequestContextID []byte         `json:"request_context_id"`
	Creator          sdk.AccAddress `json:"creator"`
}

type feedContext struct {
	Feed              feed             `json:"feed"`
	ServiceName       string           `json:"service_name"`
	Providers         []sdk.AccAddress `json:"providers"`
	Input             string           `json:"input"`
	Timeout           int64            `json:"timeout"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	ResponseThreshold uint16           `json:"response_threshold"`
	State             string           `json:"state"`
}

func (fc feedContext) Convert() interface{} {
	var providers []string
	for _, provider := range fc.Providers {
		providers = append(providers, provider.String())
	}
	return rpc.FeedContext{
		Feed: rpc.Feed{
			FeedName:         fc.Feed.FeedName,
			Description:      fc.Feed.Description,
			AggregateFunc:    fc.Feed.AggregateFunc,
			ValueJsonPath:    fc.Feed.ValueJsonPath,
			LatestHistory:    fc.Feed.LatestHistory,
			RequestContextID: rpc.RequestContextIDToString(fc.Feed.RequestContextID),
			Creator:          fc.Feed.Creator.String(),
		},
		ServiceName:       fc.ServiceName,
		Providers:         providers,
		Input:             fc.Input,
		Timeout:           fc.Timeout,
		ServiceFeeCap:     fc.ServiceFeeCap,
		RepeatedFrequency: fc.RepeatedFrequency,
		RepeatedTotal:     fc.RepeatedTotal,
		ResponseThreshold: fc.ResponseThreshold,
		State:             fc.State,
	}
}

type feedContexts []feedContext

func (fcs feedContexts) Convert() interface{} {
	result := make([]rpc.FeedContext, len(fcs))
	for _, fc := range fcs {
		result = append(result, fc.Convert().(rpc.FeedContext))
	}
	return result
}

type feedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
type feedValues []feedValue

func (fvs feedValues) Convert() interface{} {
	result := make([]rpc.FeedValue, len(fvs))
	for _, fv := range fvs {
		result = append(result, rpc.FeedValue{
			Data:      fv.Data,
			Timestamp: fv.Timestamp,
		})
	}
	return result
}

func registerCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgCreateFeed{}, "irishub/oracle/MsgCreateFeed")
	cdc.RegisterConcrete(MsgStartFeed{}, "irishub/oracle/MsgStartFeed")
	cdc.RegisterConcrete(MsgPauseFeed{}, "irishub/oracle/MsgPauseFeed")
	cdc.RegisterConcrete(MsgEditFeed{}, "irishub/oracle/MsgEditFeed")

	cdc.RegisterConcrete(feed{}, "irishub/oracle/Feed")
	cdc.RegisterConcrete(feedContext{}, "irishub/oracle/FeedContext")
}
