package service

import (
	"errors"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/json"
)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = MsgBindService{}
	//_ sdk.Msg = MsgUpdateServiceBinding{}
	//_ sdk.Msg = MsgSetWithdrawAddress{}
	//_ sdk.Msg = MsgDisableService{}
	//_ sdk.Msg = MsgEnableService{}
	//_ sdk.Msg = MsgRefundServiceDeposit{}
	_ sdk.Msg = MsgRequestService{}
	_ sdk.Msg = MsgRespondService{}
	//_ types.Msg = MsgPauseRequestContext{}
	//_ types.Msg = MsgStartRequestContext{}
	//_ types.Msg = MsgKillRequestContext{}
	//_ types.Msg = MsgUpdateRequestContext{}
	//_ types.Msg = MsgWithdrawEarnedFees{}
	//_ types.Msg = MsgWithdrawTax{}

	TagServiceName      = "service-name"
	TagProvider         = "provider"
	TagConsumer         = "consumer"
	TagRequestID        = "request-id"
	TagRespondService   = "respond_service"
	TagRequestContextID = "request-context-id"

	cdc = sdk.NewAminoCodec()
)

func init() {
	RegisterCodec(cdc)
}

// MsgDefineService defines a message to define a service
type MsgDefineService struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

func (msg MsgDefineService) Type() string {
	return "define_service"
}

func (msg MsgDefineService) ValidateBasic() error {
	if len(msg.Author) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Name) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Schemas) == 0 {
		return errors.New("schemas missing")
	}

	return nil
}

func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgDefineService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

// MsgBindService defines a message to bind a service
type MsgBindService struct {
	ServiceName     string         `json:"service_name"`
	Provider        sdk.AccAddress `json:"provider"`
	Deposit         sdk.Coins      `json:"deposit"`
	Pricing         string         `json:"pricing"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
}

func (msg MsgBindService) Type() string {
	return "bind_service"
}

func (msg MsgBindService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Pricing) == 0 {
		return errors.New("pricing missing")
	}
	return nil
}

func (msg MsgBindService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgBindService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

// MsgRequestService defines a message to request a service
type MsgRequestService struct {
	ServiceName       string           `json:"service_name"`
	Providers         []sdk.AccAddress `json:"providers"`
	Consumer          sdk.AccAddress   `json:"consumer"`
	Input             string           `json:"input"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	Timeout           int64            `json:"timeout"`
	SuperMode         bool             `json:"super_mode"`
	Repeated          bool             `json:"repeated"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
}

func (msg MsgRequestService) Type() string {
	return "request_service"
}

func (msg MsgRequestService) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	if len(msg.Providers) == 0 {
		return errors.New("providers missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Input) == 0 {
		return errors.New("input missing")
	}
	return nil
}

func (msg MsgRequestService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgRequestService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

// MsgRespondService defines a message to respond a service request
type MsgRespondService struct {
	RequestID string         `json:"request_id"`
	Provider  sdk.AccAddress `json:"provider"`
	Output    string         `json:"output"`
	Error     string         `json:"error"`
}

func (msg MsgRespondService) Type() string {
	return "respond_service"
}

func (msg MsgRespondService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Output) == 0 && len(msg.Error) == 0 {
		return errors.New("either output or error should be specified, but neither was provided")
	}

	if len(msg.Output) > 0 && len(msg.Error) > 0 {
		return errors.New("either output or error should be specified, but both were provided")
	}
	return nil
}

func (msg MsgRespondService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgRespondService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//================================for ABCI query
// QueryRequestParams defines the params to query a request for a service binding
type QueryRequestParams struct {
	RequestID string
}

func RegisterCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgDefineService{}, "irishub/service/MsgDefineService")
	cdc.RegisterConcrete(MsgBindService{}, "irishub/service/MsgBindService")
	//cdc.RegisterConcrete(MsgUpdateServiceBinding{}, "irishub/service/MsgUpdateServiceBinding")
	//cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irishub/service/MsgSetWithdrawAddress")
	//cdc.RegisterConcrete(MsgDisableService{}, "irishub/service/MsgDisableService")
	//cdc.RegisterConcrete(MsgEnableService{}, "irishub/service/MsgEnableService")
	//cdc.RegisterConcrete(MsgRefundServiceDeposit{}, "irishub/service/MsgRefundServiceDeposit")
	cdc.RegisterConcrete(MsgRequestService{}, "irishub/service/MsgRequestService")
	cdc.RegisterConcrete(MsgRespondService{}, "irishub/service/MsgRespondService")
	//cdc.RegisterConcrete(MsgPauseRequestContext{}, "irishub/service/MsgPauseRequestContext")
	//cdc.RegisterConcrete(MsgStartRequestContext{}, "irishub/service/MsgStartRequestContext")
	//cdc.RegisterConcrete(MsgKillRequestContext{}, "irishub/service/MsgKillRequestContext")
	//cdc.RegisterConcrete(MsgUpdateRequestContext{}, "irishub/service/MsgUpdateRequestContext")
	//cdc.RegisterConcrete(MsgWithdrawEarnedFees{}, "irishub/service/MsgWithdrawEarnedFees")
	//cdc.RegisterConcrete(MsgWithdrawTax{}, "irishub/service/MsgWithdrawTax")

}
