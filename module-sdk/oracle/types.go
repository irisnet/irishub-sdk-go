package oracle

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/irisnet/core-sdk-go/types"
)

const (
	ModuleName = "oracle"
)

var (
	_ sdk.Msg = &MsgCreateFeed{}
	_ sdk.Msg = &MsgStartFeed{}
	_ sdk.Msg = &MsgPauseFeed{}
	_ sdk.Msg = &MsgEditFeed{}

	// the feed/service name only accepts alphanumeric characters, _ and -
	regPlainText = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
)

// Route implements Msg.
func (msg MsgCreateFeed) Route() string {
	return ModuleName
}

// Type implements Msg.
func (msg MsgCreateFeed) Type() string {
	return "create_feed"
}

// ValidateBasic implements Msg.
func (msg MsgCreateFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return sdk.Wrapf("missing feed name")
	}
	if !regPlainText.MatchString(feedName) {
		return sdk.Wrapf("invalid feed name: %s", feedName)
	}

	if len(msg.Description) == 0 {
		return sdk.Wrapf("missing description")
	}

	if len(msg.ServiceName) == 0 {
		return sdk.Wrapf("missing name")
	}
	if !regPlainText.MatchString(msg.ServiceName) {
		return sdk.Wrapf("invalid service name %s", msg.ServiceName)
	}

	if msg.LatestHistory == 0 {
		return sdk.Wrapf("missing latest history")
	}

	if err := validateTimeout(msg.Timeout, msg.RepeatedFrequency); err != nil {
		return err
	}
	if len(msg.Providers) == 0 {
		return sdk.Wrapf("providers missing")
	}

	if len(msg.AggregateFunc) == 0 {
		return sdk.Wrapf("missing aggregateFunc")
	}

	if len(msg.ValueJsonPath) == 0 {
		return sdk.Wrapf("missing valueJsonPath")
	}

	if !msg.ServiceFeeCap.IsValid() {
		return sdk.Wrapf(msg.ServiceFeeCap.String())
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdk.Wrapf("invalid creator")
	}

	return validateResponseThreshold(msg.ResponseThreshold, len(msg.Providers))
}

// GetSignBytes implements Msg.
func (msg MsgCreateFeed) GetSignBytes() []byte {
	if len(msg.Providers) == 0 {
		msg.Providers = nil
	}
	if msg.ServiceFeeCap.Empty() {
		msg.ServiceFeeCap = nil
	}
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgCreateFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg MsgStartFeed) Route() string {
	return ModuleName
}

func (msg MsgStartFeed) Type() string {
	return "start_feed"
}

func (msg MsgStartFeed) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdk.Wrapf("invalid creator")
	}

	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return sdk.Wrapf("missing feed name")
	}
	if !regPlainText.MatchString(feedName) {
		return sdk.Wrapf("invalid feed name: %s", feedName)
	}
	return nil
}

func (msg MsgStartFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgStartFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg MsgPauseFeed) Route() string {
	return ModuleName
}

func (msg MsgPauseFeed) Type() string {
	return "pause_feed"
}

func (msg MsgPauseFeed) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdk.Wrapf("invalid creator")
	}

	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return sdk.Wrapf("missing feed name")
	}
	if !regPlainText.MatchString(feedName) {
		return sdk.Wrapf("invalid feed name: %s", feedName)
	}
	return nil
}

func (msg MsgPauseFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgPauseFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg MsgEditFeed) Route() string {
	return ModuleName
}

func (msg MsgEditFeed) Type() string {
	return "edit_feed"
}

func (msg MsgEditFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return sdk.Wrapf("missing feed name")
	}
	if !regPlainText.MatchString(feedName) {
		return sdk.Wrapf("invalid feed name: %s", feedName)
	}

	if len(msg.Description) == 0 {
		return sdk.Wrapf("missing description")
	}

	if msg.ServiceFeeCap != nil && !msg.ServiceFeeCap.IsValid() {
		return sdk.Wrapf(msg.ServiceFeeCap.String())
	}
	if msg.Timeout != 0 && msg.RepeatedFrequency != 0 {
		if err := validateTimeout(msg.Timeout, msg.RepeatedFrequency); err != nil {
			return err
		}
	}
	if msg.ResponseThreshold != 0 {
		if err := validateResponseThreshold(msg.ResponseThreshold, len(msg.Providers)); err != nil {
			return err
		}
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdk.Wrapf("invalid creator")
	}
	return nil
}

func (msg MsgEditFeed) GetSignBytes() []byte {
	if len(msg.Providers) == 0 {
		msg.Providers = nil
	}
	if msg.ServiceFeeCap.Empty() {
		msg.ServiceFeeCap = nil
	}
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgEditFeed) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func validateResponseThreshold(responseThreshold uint32, maxCnt int) error {
	if (maxCnt != 0 && int(responseThreshold) > maxCnt) || responseThreshold < 1 {
		return sdk.Wrapf("response threshold should be between 1 and %d", maxCnt)
	}
	return nil
}

func validateTimeout(timeout int64, frequency uint64) error {
	if frequency < uint64(timeout) {
		return sdk.Wrapf("timeout [%d] should be no more than frequency [%d]", timeout, frequency)
	}
	return nil
}

// String implements fmt.Stringer
func (f FeedContext) String() string {
	var bf bytes.Buffer
	for _, addr := range f.Providers {
		bf.WriteString(addr)
		bf.WriteString(",")
	}
	return fmt.Sprintf(` FeedContext:
	%s
	ServiceName:                %s
	Providers:                  %s
	Input:                      %s
	Timeout:                    %d
	ServiceFeeCap:              %s
	RepeatedFrequency:          %d
	ResponseThreshold:          %d
	State:                      %s`,
		f.Feed.String(),
		f.ServiceName,
		bf.String(),
		f.Input,
		f.Timeout,
		f.ServiceFeeCap,
		f.RepeatedFrequency,
		f.ResponseThreshold,
		f.State.String(),
	)
}

func (f FeedContext) Convert() interface{} {
	return QueryFeedResp{
		Feed: struct {
			FeedName         string `json:"feed_name"`
			Description      string `json:"description"`
			AggregateFunc    string `json:"aggregate_func"`
			ValueJsonPath    string `json:"value_json_path"`
			LatestHistory    uint64 `json:"latest_history"`
			RequestContextID string `json:"request_context_id"`
			Creator          string `json:"creator"`
		}{
			f.Feed.FeedName,
			f.Feed.Description,
			f.Feed.AggregateFunc,
			f.Feed.ValueJsonPath,
			f.Feed.LatestHistory,
			f.Feed.RequestContextID,
			f.Feed.Creator,
		},
		ServiceName:       f.ServiceName,
		Providers:         f.Providers,
		Input:             f.Input,
		Timeout:           f.Timeout,
		ServiceFeeCap:     f.ServiceFeeCap,
		RepeatedFrequency: f.RepeatedFrequency,
		ResponseThreshold: f.ResponseThreshold,
		State:             int32(f.State),
	}
}

type feedContexts []FeedContext

func (fs feedContexts) Convert() interface{} {
	var res []QueryFeedResp
	for _, f := range fs {
		res = append(res, f.Convert().(QueryFeedResp))
	}
	return res
}

type feedValues []FeedValue

func (f FeedValue) Convert() interface{} {
	return QueryFeedValueResp{
		Data:      f.Data,
		Timestamp: f.Timestamp,
	}
}

func (fs feedValues) Convert() interface{} {
	var res []QueryFeedValueResp
	for _, f := range fs {
		res = append(res, f.Convert().(QueryFeedValueResp))
	}
	return res
}
