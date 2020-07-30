package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"time"
)

type ServiceTx interface {
	DefineService(request ServiceDefinitionRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	BindService(request ServiceBindingRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	UpdateServiceBinding(request ServiceBindingUpdateRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	InvokeService(request ServiceInvocationRequest, baseTx sdk.BaseTx) (requestContextID string, err sdk.Error)

	SetWithdrawAddress(withdrawAddress string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	DisableServiceBinding(serviceName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	EnableServiceBinding(serviceName string,
		deposit sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	RefundServiceDeposit(serviceName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	PauseRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	StartRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	KillRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	UpdateRequestContext(request UpdateContextRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	WithdrawEarnedFees(baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	WithdrawTax(destAddress string,
		amount sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	SubscribeServiceRequest(serviceRegistry ServiceRegistry, baseTx sdk.BaseTx) (sdk.Subscription, sdk.Error)

	SubscribeSingleServiceRequest(serviceName string,
		callback ServiceRespondCallback,
		baseTx sdk.BaseTx) (sdk.Subscription, sdk.Error)

	SubscribeServiceResponse(reqCtxID string,
		callback ServiceInvokeCallback) (sdk.Subscription, sdk.Error)
}

type ServiceQuery interface {
	QueryDefinition(serviceName string) (ServiceDefinition, sdk.Error)

	QueryBinding(serviceName string, provider sdk.AccAddress) (ServiceBinding, sdk.Error)
	QueryBindings(serviceName string) ([]ServiceBinding, sdk.Error)

	QueryRequest(requestID string) (ServiceRequest, sdk.Error)
	QueryRequests(serviceName string, provider sdk.AccAddress) ([]ServiceRequest, sdk.Error)
	QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) ([]ServiceRequest, sdk.Error)

	QueryResponse(requestID string) (ServiceResponse, sdk.Error)
	QueryResponses(requestContextID string, batchCounter uint64) ([]ServiceResponse, sdk.Error)

	QueryRequestContext(requestContextID string) (RequestContext, sdk.Error)
	QueryFees(provider string) (EarnedFees, sdk.Error)
}

type Service interface {
	sdk.Module
	ServiceTx
	ServiceQuery
}

type ServiceInvokeCallback func(reqCtxID, reqID, responses string)
type ServiceRespondCallback func(reqCtxID, reqID, input string) (output string, result string)
type ServiceRegistry map[string]ServiceRespondCallback

// ServiceRequest defines a request which contains the detailed request data
type ServiceRequest struct {
	ID                         string         `json:"id"`
	ServiceName                string         `json:"service_name"`
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Input                      string         `json:"input"`
	ServiceFee                 sdk.Coins      `json:"service_fee"`
	SuperMode                  bool           `json:"super_mode"`
	RequestHeight              int64          `json:"request_height"`
	ExpirationHeight           int64          `json:"expiration_height"`
	RequestContextID           string         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

// ServiceResponse defines a response
type ServiceResponse struct {
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Output                     string         `json:"output"`
	Result                     string         `json:"error"`
	RequestContextID           string         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
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
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

type ServiceBindingRequest struct {
	ServiceName string       `json:"service_name"`
	Deposit     sdk.DecCoins `json:"deposit"`
	Pricing     string       `json:"pricing"`
	MinRespTime uint64       `json:"min_resp_time"`
}

// ServiceBindingUpdateRequest defines a message to update a service binding
type ServiceBindingUpdateRequest struct {
	ServiceName string       `json:"service_name"`
	Deposit     sdk.DecCoins `json:"deposit"`
	Pricing     string       `json:"pricing"`
}

// ServiceBinding defines a struct for service binding
type ServiceBinding struct {
	ServiceName  string    `json:"service_name"`
	Provider     string    `json:"provider"`
	Deposit      sdk.Coins `json:"deposit"`
	Pricing      string    `json:"pricing"`
	Qos          uint64    `json:"qos"`
	Available    bool      `json:"available"`
	DisabledTime time.Time `json:"disabled_time"`
	Owner        string    `json:"owner"`
}

type ServiceInvocationRequest struct {
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Input             string       `json:"input"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	SuperMode         bool         `json:"super_mode"`
	Repeated          bool         `json:"repeated"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	Callback          ServiceInvokeCallback
}

// UpdateContextRequest defines a message to update a request context
type UpdateContextRequest struct {
	RequestContextID  string       `json:"request_context_id"`
	Providers         []string     `json:"providers"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
}

// RequestContext defines a context which holds request-related data
type RequestContext struct {
	ServiceName        string           `json:"service_name"`
	Providers          []sdk.AccAddress `json:"providers"`
	Consumer           sdk.AccAddress   `json:"consumer"`
	Input              string           `json:"input"`
	ServiceFeeCap      sdk.Coins        `json:"service_fee_cap"`
	ModuleName         string           `json:"module_name"`
	Timeout            int64            `json:"timeout"`
	SuperMode          bool             `json:"super_mode"`
	Repeated           bool             `json:"repeated"`
	RepeatedFrequency  uint64           `json:"repeated_frequency"`
	RepeatedTotal      int64            `json:"repeated_total"`
	BatchCounter       uint64           `json:"batch_counter"`
	BatchRequestCount  uint32           `json:"batch_request_count"`
	BatchResponseCount uint32           `json:"batch_response_count"`
	ResponseThreshold  uint32           `json:"response_threshold"`
	BatchState         int32            `json:"batch_state"`
	State              int32            `json:"state"`
}

// EarnedFees defines a struct for the fees earned by the provider
type EarnedFees struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}
