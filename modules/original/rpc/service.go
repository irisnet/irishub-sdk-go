package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
	"time"
)

type ServiceTx interface {
	DefineService(request ServiceDefinitionRequest, baseTx original.BaseTx) (original.ResultTx, original.Error)

	BindService(request ServiceBindingRequest, baseTx original.BaseTx) (original.ResultTx, original.Error)

	UpdateServiceBinding(request ServiceBindingUpdateRequest, baseTx original.BaseTx) (original.ResultTx, original.Error)

	InvokeService(request ServiceInvocationRequest, baseTx original.BaseTx) (requestContextID string, err original.Error)

	SetWithdrawAddress(withdrawAddress string, baseTx original.BaseTx) (original.ResultTx, original.Error)

	DisableServiceBinding(serviceName string, baseTx original.BaseTx) (original.ResultTx, original.Error)

	EnableServiceBinding(serviceName string,
		deposit original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error)

	RefundServiceDeposit(serviceName string, baseTx original.BaseTx) (original.ResultTx, original.Error)

	PauseRequestContext(requestContextID string, baseTx original.BaseTx) (original.ResultTx, original.Error)

	StartRequestContext(requestContextID string, baseTx original.BaseTx) (original.ResultTx, original.Error)

	KillRequestContext(requestContextID string, baseTx original.BaseTx) (original.ResultTx, original.Error)

	UpdateRequestContext(request UpdateContextRequest, baseTx original.BaseTx) (original.ResultTx, original.Error)

	WithdrawEarnedFees(baseTx original.BaseTx) (original.ResultTx, original.Error)

	WithdrawTax(destAddress string,
		amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error)

	SubscribeServiceRequest(serviceName string,
		callback ServiceRespondCallback,
		baseTx original.BaseTx) (original.Subscription, original.Error)

	SubscribeServiceResponse(reqCtxID string,
		callback ServiceInvokeCallback) (original.Subscription, original.Error)
}

type ServiceQuery interface {
	QueryDefinition(serviceName string) (ServiceDefinition, original.Error)

	QueryBinding(serviceName string, provider original.AccAddress) (ServiceBinding, original.Error)
	QueryBindings(serviceName string) ([]ServiceBinding, original.Error)

	QueryRequest(requestID string) (ServiceRequest, original.Error)
	QueryRequests(serviceName string, provider original.AccAddress) ([]ServiceRequest, original.Error)
	QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) ([]ServiceRequest, original.Error)

	QueryResponse(requestID string) (ServiceResponse, original.Error)
	QueryResponses(requestContextID string, batchCounter uint64) ([]ServiceResponse, original.Error)

	QueryRequestContext(requestContextID string) (RequestContext, original.Error)
	QueryFees(provider string) (original.Coins, original.Error)
}

type Service interface {
	original.Module
	ServiceTx
	ServiceQuery
}

type ServiceInvokeCallback func(reqCtxID, reqID, responses string)
type ServiceRespondCallback func(reqCtxID, reqID, input string) (output string, result string)
type ServiceRegistry map[string]ServiceRespondCallback

// ServiceRequest defines a request which contains the detailed request data
type ServiceRequest struct {
	ID                         string              `json:"id"`
	ServiceName                string              `json:"service_name"`
	Provider                   original.AccAddress `json:"provider"`
	Consumer                   original.AccAddress `json:"consumer"`
	Input                      string              `json:"input"`
	ServiceFee                 original.Coins      `json:"service_fee"`
	SuperMode                  bool                `json:"super_mode"`
	RequestHeight              int64               `json:"request_height"`
	ExpirationHeight           int64               `json:"expiration_height"`
	RequestContextID           string              `json:"request_context_id"`
	RequestContextBatchCounter uint64              `json:"request_context_batch_counter"`
}

// ServiceResponse defines a response
type ServiceResponse struct {
	Provider                   original.AccAddress `json:"provider"`
	Consumer                   original.AccAddress `json:"consumer"`
	Output                     string              `json:"output"`
	Result                     string              `json:"result"`
	RequestContextID           string              `json:"request_context_id"`
	RequestContextBatchCounter uint64              `json:"request_context_batch_counter"`
}

// ServiceDefinitionRequest defines the request parameters of the service definition
type ServiceDefinitionRequest struct {
	ServiceName       string   `json:"service_name"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	AuthorDescription string   `json:"author_description"`
	Schemas           string   `json:"schemas"`
}

// ServiceDefinition represents a service definition
type ServiceDefinition struct {
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Tags              []string            `json:"tags"`
	Author            original.AccAddress `json:"author"`
	AuthorDescription string              `json:"author_description"`
	Schemas           string              `json:"schemas"`
}

type ServiceBindingRequest struct {
	ServiceName string            `json:"service_name"`
	Deposit     original.DecCoins `json:"deposit"`
	Pricing     string            `json:"pricing"`
	MinRespTime uint64            `json:"min_resp_time"`
	Qos         uint64            `json:"Qos"`
}

// ServiceBindingUpdateRequest defines a message to update a service binding
type ServiceBindingUpdateRequest struct {
	ServiceName string            `json:"service_name"`
	Deposit     original.DecCoins `json:"deposit"`
	Pricing     string            `json:"pricing"`
}

// ServiceBinding defines a struct for service binding
type ServiceBinding struct {
	ServiceName  string         `json:"service_name"`
	Provider     string         `json:"provider"`
	Deposit      original.Coins `json:"deposit"`
	Pricing      string         `json:"pricing"`
	Qos          uint64         `json:"qos"`
	Owner        string         `json:"owner"`
	Available    bool           `json:"available"`
	DisabledTime time.Time      `json:"disabled_time"`
}

type ServiceInvocationRequest struct {
	ServiceName       string            `json:"service_name"`
	Providers         []string          `json:"providers"`
	Input             string            `json:"input"`
	ServiceFeeCap     original.DecCoins `json:"service_fee_cap"`
	Timeout           int64             `json:"timeout"`
	SuperMode         bool              `json:"super_mode"`
	Repeated          bool              `json:"repeated"`
	RepeatedFrequency uint64            `json:"repeated_frequency"`
	RepeatedTotal     int64             `json:"repeated_total"`
	Callback          ServiceInvokeCallback
}

// UpdateContextRequest defines a message to update a request context
type UpdateContextRequest struct {
	RequestContextID  string            `json:"request_context_id"`
	Providers         []string          `json:"providers"`
	ServiceFeeCap     original.DecCoins `json:"service_fee_cap"`
	Timeout           int64             `json:"timeout"`
	RepeatedFrequency uint64            `json:"repeated_frequency"`
	RepeatedTotal     int64             `json:"repeated_total"`
}

// RequestContext defines a context which holds request-related data
type RequestContext struct {
	ServiceName        string                `json:"service_name"`
	Providers          []original.AccAddress `json:"providers"`
	Consumer           original.AccAddress   `json:"consumer"`
	Input              string                `json:"input"`
	ServiceFeeCap      original.Coins        `json:"service_fee_cap"`
	ModuleName         string                `json:"module_name"`
	Timeout            int64                 `json:"timeout"`
	SuperMode          bool                  `json:"super_mode"`
	Repeated           bool                  `json:"repeated"`
	RepeatedFrequency  uint64                `json:"repeated_frequency"`
	RepeatedTotal      int64                 `json:"repeated_total"`
	BatchCounter       uint64                `json:"batch_counter"`
	BatchRequestCount  uint32                `json:"batch_request_count"`
	BatchResponseCount uint32                `json:"batch_response_count"`
	ResponseThreshold  uint32                `json:"response_threshold"`
	BatchState         int32                 `json:"batch_state"`
	State              int32                 `json:"state"`
}

// EarnedFees defines a struct for the fees earned by the provider
type EarnedFees struct {
	Address original.AccAddress `json:"address"`
	Coins   original.Coins      `json:"coins"`
}
