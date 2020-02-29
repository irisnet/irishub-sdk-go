package oracle

import (
	"errors"
	"strings"

	"github.com/irisnet/irishub-sdk-go/tools/json"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

var (
	_ sdk.Msg = MsgCreateFeed{}
	_ sdk.Msg = MsgStartFeed{}
	_ sdk.Msg = MsgPauseFeed{}
	_ sdk.Msg = MsgEditFeed{}

	cdc = sdk.NewAminoCodec()
)

func init() {
	RegisterCodec(cdc)
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
type Feed struct {
	FeedName         string         `json:"feed_name"`
	Description      string         `json:"description"`
	AggregateFunc    string         `json:"aggregate_func"`
	ValueJsonPath    string         `json:"value_json_path"`
	LatestHistory    uint64         `json:"latest_history"`
	RequestContextID []byte         `json:"request_context_id"`
	Creator          sdk.AccAddress `json:"creator"`
}

type FeedContext struct {
	Feed              Feed             `json:"feed"`
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

func (ctx FeedContext) toSDKFeedContext() sdk.FeedContext {
	var providers []string
	for _, provider := range ctx.Providers {
		providers = append(providers, provider.String())
	}
	return sdk.FeedContext{
		Feed: sdk.Feed{
			FeedName:         ctx.Feed.FeedName,
			Description:      ctx.Feed.Description,
			AggregateFunc:    ctx.Feed.AggregateFunc,
			ValueJsonPath:    ctx.Feed.ValueJsonPath,
			LatestHistory:    ctx.Feed.LatestHistory,
			RequestContextID: sdk.RequestContextIDToString(ctx.Feed.RequestContextID),
			Creator:          ctx.Feed.Creator.String(),
		},
		ServiceName:       ctx.ServiceName,
		Providers:         providers,
		Input:             ctx.Input,
		Timeout:           ctx.Timeout,
		ServiceFeeCap:     ctx.ServiceFeeCap,
		RepeatedFrequency: ctx.RepeatedFrequency,
		RepeatedTotal:     ctx.RepeatedTotal,
		ResponseThreshold: ctx.ResponseThreshold,
		State:             ctx.State,
	}
}

type FeedContexts []FeedContext

func (ctx FeedContexts) toSDKFeedContexts() (result []sdk.FeedContext) {
	for _, feedCtx := range ctx {
		result = append(result, feedCtx.toSDKFeedContext())
	}
	return
}

func RegisterCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgCreateFeed{}, "irishub/oracle/MsgCreateFeed")
	cdc.RegisterConcrete(MsgStartFeed{}, "irishub/oracle/MsgStartFeed")
	cdc.RegisterConcrete(MsgPauseFeed{}, "irishub/oracle/MsgPauseFeed")
	cdc.RegisterConcrete(MsgEditFeed{}, "irishub/oracle/MsgEditFeed")

	cdc.RegisterConcrete(Feed{}, "irishub/oracle/Feed")
	cdc.RegisterConcrete(FeedContext{}, "irishub/oracle/FeedContext")
}
