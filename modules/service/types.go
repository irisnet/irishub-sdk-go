package service

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/json"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "service"

	tagServiceName      = "service-name"
	tagProvider         = "provider"
	tagConsumer         = "consumer"
	tagRequestID        = "request-id"
	tagRespondService   = "respond_service"
	tagRequestContextID = "request-context-id"

	ActionNewBatchRequest = "new-batch-request."

	requestContextIDLen = 40 // length of the request context ID in bytes
)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = MsgBindService{}
	_ sdk.Msg = MsgUpdateServiceBinding{}
	_ sdk.Msg = MsgSetWithdrawAddress{}
	_ sdk.Msg = MsgDisableService{}
	_ sdk.Msg = MsgEnableService{}
	_ sdk.Msg = MsgRefundServiceDeposit{}
	_ sdk.Msg = MsgRequestService{}
	_ sdk.Msg = MsgRespondService{}
	_ sdk.Msg = MsgPauseRequestContext{}
	_ sdk.Msg = MsgStartRequestContext{}
	_ sdk.Msg = MsgKillRequestContext{}
	_ sdk.Msg = MsgUpdateRequestContext{}
	_ sdk.Msg = MsgWithdrawEarnedFees{}
	_ sdk.Msg = MsgWithdrawTax{}

	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
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
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
	Pricing     string         `json:"pricing"`
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
	Result    string         `json:"result"`
	Output    string         `json:"output"`
}

func (msg MsgRespondService) Type() string {
	return "respond_service"
}

func (msg MsgRespondService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Result) == 0 {
		return errors.New("result missing")
	}

	if err := ValidateResponseResult(msg.Result); err != nil {
		return err
	}

	result, err := ParseResult(msg.Result)
	if err != nil {
		return err
	}

	if result.Code == 200 && len(msg.Output) == 0 {
		return errors.New("output must be specified when the result code is 200")
	}

	if result.Code != 200 && len(msg.Output) != 0 {
		return errors.New("output should not be specified when the result code is not 200")
	}

	if len(msg.Output) > 0 {
		if !json2.Valid([]byte(msg.Output)) {
			return errors.New("output is not valid JSON")
		}
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

//______________________________________________________________________

// MsgUpdateServiceBinding defines a message to update a service binding
type MsgUpdateServiceBinding struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
	Pricing     string         `json:"pricing"`
}

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return "update_service_binding" }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return errors.New(fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSetWithdrawAddress defines a message to set the withdrawal address for a service binding
type MsgSetWithdrawAddress struct {
	Provider        sdk.AccAddress `json:"provider"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
}

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return "set_withdraw_address" }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.WithdrawAddress) == 0 {
		return errors.New("withdrawal address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgDisableService defines a message to disable a service binding
type MsgDisableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// Type implements Msg.
func (msg MsgDisableService) Type() string { return "disable_service" }

// GetSignBytes implements Msg.
func (msg MsgDisableService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgDisableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgEnableService defines a message to enable a service binding
type MsgEnableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
}

// Type implements Msg.
func (msg MsgEnableService) Type() string { return "enable_service" }

// GetSignBytes implements Msg.
func (msg MsgEnableService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return errors.New(fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgRefundServiceDeposit defines a message to refund deposit from a service binding
type MsgRefundServiceDeposit struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return "refund_service_deposit" }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgPauseRequestContext defines a message to suspend a request context
type MsgPauseRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// Type implements Msg.
func (msg MsgPauseRequestContext) Type() string { return "pause_request_context" }

// GetSignBytes implements Msg.
func (msg MsgPauseRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgPauseRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != requestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", requestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgPauseRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgStartRequestContext defines a message to resume a request context
type MsgStartRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// Type implements Msg.
func (msg MsgStartRequestContext) Type() string { return "start_request_context" }

// GetSignBytes implements Msg.
func (msg MsgStartRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgStartRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != requestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", requestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgStartRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgKillRequestContext defines a message to terminate a request context
type MsgKillRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// Type implements Msg.
func (msg MsgKillRequestContext) Type() string { return "kill_request_context" }

// GetSignBytes implements Msg.
func (msg MsgKillRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgKillRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != requestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", requestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgKillRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgUpdateRequestContext defines a message to update a request context
type MsgUpdateRequestContext struct {
	RequestContextID  []byte           `json:"request_context_id"`
	Providers         []sdk.AccAddress `json:"providers"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	Timeout           int64            `json:"timeout"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	Consumer          sdk.AccAddress   `json:"consumer"`
}

// Type implements Msg.
func (msg MsgUpdateRequestContext) Type() string { return "update_request_context" }

// GetSignBytes implements Msg.
func (msg MsgUpdateRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != requestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", requestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgWithdrawEarnedFees defines a message to withdraw the fees earned by the provider
type MsgWithdrawEarnedFees struct {
	Provider sdk.AccAddress `json:"provider"`
}

// Type implements Msg.
func (msg MsgWithdrawEarnedFees) Type() string { return "withdraw_earned_fees" }

// GetSignBytes implements Msg.
func (msg MsgWithdrawEarnedFees) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawEarnedFees) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawEarnedFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgWithdrawTax defines a message to withdraw the service tax
type MsgWithdrawTax struct {
	Trustee     sdk.AccAddress `json:"trustee"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Amount      sdk.Coins      `json:"amount"`
}

// Type implements Msg.
func (msg MsgWithdrawTax) Type() string { return "withdraw_tax" }

// GetSignBytes implements Msg.
func (msg MsgWithdrawTax) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawTax) ValidateBasic() error {
	if len(msg.Trustee) == 0 {
		return errors.New("trustee missing")
	}

	if len(msg.DestAddress) == 0 {
		return errors.New("destination address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawTax) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trustee}
}

//==========================================for QueryWithResponse==========================================

// serviceDefinition represents a service definition
type serviceDefinition struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

func (r serviceDefinition) Convert() interface{} {
	return rpc.ServiceDefinition{
		Name:              r.Name,
		Description:       r.Description,
		Tags:              r.Tags,
		Author:            r.Author,
		AuthorDescription: r.AuthorDescription,
		Schemas:           r.Schemas,
	}
}

// serviceBinding defines a struct for service binding
type serviceBinding struct {
	ServiceName     string         `json:"service_name"`
	Provider        sdk.AccAddress `json:"provider"`
	Deposit         sdk.Coins      `json:"deposit"`
	Pricing         string         `json:"pricing"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
	Available       bool           `json:"available"`
	DisabledTime    time.Time      `json:"disabled_time"`
}

func (b serviceBinding) Convert() interface{} {
	return rpc.ServiceBinding{
		ServiceName:     b.ServiceName,
		Provider:        b.Provider,
		Deposit:         b.Deposit,
		Pricing:         b.Pricing,
		WithdrawAddress: b.WithdrawAddress,
		Available:       b.Available,
		DisabledTime:    b.DisabledTime,
	}

}

type serviceBindings []serviceBinding

func (bs serviceBindings) Convert() interface{} {
	bindings := make([]rpc.ServiceBinding, len(bs))
	for _, binding := range bs {
		bindings = append(bindings, binding.Convert().(rpc.ServiceBinding))
	}
	return bindings
}

// request defines a request which contains the detailed request data
type request struct {
	ServiceName                string         `json:"service_name"`
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Input                      string         `json:"input"`
	ServiceFee                 sdk.Coins      `json:"service_fee"`
	SuperMode                  bool           `json:"super_mode"`
	RequestHeight              int64          `json:"request_height"`
	ExpirationHeight           int64          `json:"expiration_height"`
	RequestContextID           []byte         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

func (r request) Convert() interface{} {
	return rpc.ServiceRequest{
		ServiceName:                r.ServiceName,
		Provider:                   r.Provider,
		Consumer:                   r.Consumer,
		Input:                      r.Input,
		ServiceFee:                 r.ServiceFee,
		SuperMode:                  r.SuperMode,
		RequestHeight:              r.RequestHeight,
		ExpirationHeight:           r.ExpirationHeight,
		RequestContextID:           rpc.RequestContextIDToString(r.RequestContextID),
		RequestContextBatchCounter: r.RequestContextBatchCounter,
	}
}

type requests []request

func (rs requests) Convert() interface{} {
	requests := make([]rpc.ServiceRequest, len(rs))
	for _, request := range rs {
		requests = append(requests, request.Convert().(rpc.ServiceRequest))
	}
	return requests
}

// ServiceResponse defines a response
type response struct {
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Output                     string         `json:"output"`
	Error                      string         `json:"error"`
	RequestContextID           []byte         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

func (r response) Convert() interface{} {
	return rpc.ServiceResponse{
		Provider:                   r.Provider,
		Consumer:                   r.Consumer,
		Output:                     r.Output,
		Error:                      r.Error,
		RequestContextID:           rpc.RequestContextIDToString(r.RequestContextID),
		RequestContextBatchCounter: r.RequestContextBatchCounter,
	}
}

type responses []response

func (rs responses) Convert() interface{} {
	responses := make([]rpc.ServiceResponse, len(rs))
	for _, response := range rs {
		responses = append(responses, response.Convert().(rpc.ServiceResponse))
	}
	return responses
}

// requestContext defines a context which holds request-related data
type requestContext struct {
	ServiceName        string           `json:"service_name"`
	Providers          []sdk.AccAddress `json:"providers"`
	Consumer           sdk.AccAddress   `json:"consumer"`
	Input              string           `json:"input"`
	ServiceFeeCap      sdk.Coins        `json:"service_fee_cap"`
	Timeout            int64            `json:"timeout"`
	SuperMode          bool             `json:"super_mode"`
	Repeated           bool             `json:"repeated"`
	RepeatedFrequency  uint64           `json:"repeated_frequency"`
	RepeatedTotal      int64            `json:"repeated_total"`
	BatchCounter       uint64           `json:"batch_counter"`
	BatchRequestCount  uint16           `json:"batch_request_count"`
	BatchResponseCount uint16           `json:"batch_response_count"`
	BatchState         string           `json:"batch_state"`
	State              string           `json:"state"`
	ResponseThreshold  uint16           `json:"response_threshold"`
	ModuleName         string           `json:"module_name"`
}

func (r requestContext) Convert() interface{} {
	return rpc.RequestContext{
		ServiceName:        r.ServiceName,
		Providers:          r.Providers,
		Consumer:           r.Consumer,
		Input:              r.Input,
		ServiceFeeCap:      r.ServiceFeeCap,
		Timeout:            r.Timeout,
		SuperMode:          r.SuperMode,
		Repeated:           r.Repeated,
		RepeatedFrequency:  r.RepeatedFrequency,
		RepeatedTotal:      r.RepeatedTotal,
		BatchCounter:       r.BatchCounter,
		BatchRequestCount:  r.BatchRequestCount,
		BatchResponseCount: r.BatchResponseCount,
		BatchState:         r.BatchState,
		State:              r.State,
		ResponseThreshold:  r.ResponseThreshold,
		ModuleName:         r.ModuleName,
	}
}

// earnedFees defines a struct for the fees earned by the provider
type earnedFees struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

func (e earnedFees) Convert() interface{} {
	return rpc.EarnedFees{
		Address: e.Address,
		Coins:   e.Coins,
	}
}

// Result defines a struct for the response result
type Result struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
}

// ParseResult parses the given string to Result
func ParseResult(result string) (Result, error) {
	var r Result

	if err := json2.Unmarshal([]byte(result), &r); err != nil {
		return r, fmt.Errorf("failed to unmarshal the result: %s", err)
	}

	return r, nil
}

func registerCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgDefineService{}, "irishub/service/MsgDefineService")
	cdc.RegisterConcrete(MsgBindService{}, "irishub/service/MsgBindService")
	cdc.RegisterConcrete(MsgUpdateServiceBinding{}, "irishub/service/MsgUpdateServiceBinding")
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irishub/service/MsgSetWithdrawAddress")
	cdc.RegisterConcrete(MsgDisableService{}, "irishub/service/MsgDisableService")
	cdc.RegisterConcrete(MsgEnableService{}, "irishub/service/MsgEnableService")
	cdc.RegisterConcrete(MsgRefundServiceDeposit{}, "irishub/service/MsgRefundServiceDeposit")
	cdc.RegisterConcrete(MsgRequestService{}, "irishub/service/MsgRequestService")
	cdc.RegisterConcrete(MsgRespondService{}, "irishub/service/MsgRespondService")
	cdc.RegisterConcrete(MsgPauseRequestContext{}, "irishub/service/MsgPauseRequestContext")
	cdc.RegisterConcrete(MsgStartRequestContext{}, "irishub/service/MsgStartRequestContext")
	cdc.RegisterConcrete(MsgKillRequestContext{}, "irishub/service/MsgKillRequestContext")
	cdc.RegisterConcrete(MsgUpdateRequestContext{}, "irishub/service/MsgUpdateRequestContext")
	cdc.RegisterConcrete(MsgWithdrawEarnedFees{}, "irishub/service/MsgWithdrawEarnedFees")
	cdc.RegisterConcrete(MsgWithdrawTax{}, "irishub/service/MsgWithdrawTax")

	cdc.RegisterConcrete(serviceDefinition{}, "irishub/service/ServiceDefinition")
	cdc.RegisterConcrete(serviceBinding{}, "irishub/service/ServiceBinding")
	cdc.RegisterConcrete(requestContext{}, "irishub/service/RequestContext")
	cdc.RegisterConcrete(request{}, "irishub/service/Request")
	cdc.RegisterConcrete(response{}, "irishub/service/Response")
	cdc.RegisterConcrete(earnedFees{}, "irishub/service/EarnedFees")
}
