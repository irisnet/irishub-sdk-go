package rpc

import (
	"encoding/hex"
	"time"

	"github.com/irisnet/irishub-sdk-go/types"
)

type ServiceTx interface {
	DefineService(request ServiceDefinitionRequest) (types.Result, error)

	BindService(request ServiceBindingRequest) (types.Result, error)

	UpdateServiceBinding(request UpdateServiceBindingRequest) (types.Result, error)

	InvokeService(request ServiceInvocationRequest,
		callback ServiceInvokeHandler) (requestContextID string, err error)

	SetWithdrawAddress(serviceName string, withdrawAddress string, baseTx types.BaseTx) (types.Result, error)

	DisableService(serviceName string, baseTx types.BaseTx) (types.Result, error)

	EnableService(serviceName string, deposit types.Coins, baseTx types.BaseTx) (types.Result, error)

	RefundServiceDeposit(serviceName string, baseTx types.BaseTx) (types.Result, error)

	PauseRequestContext(requestContextID string, baseTx types.BaseTx) (types.Result, error)

	StartRequestContext(requestContextID string, baseTx types.BaseTx) (types.Result, error)

	KillRequestContext(requestContextID string, baseTx types.BaseTx) (types.Result, error)

	UpdateRequestContext(request UpdateContextRequest) (types.Result, error)

	WithdrawEarnedFees(baseTx types.BaseTx) (types.Result, error)

	WithdrawTax(destAddress string, amount types.Coins, baseTx types.BaseTx) (types.Result, error)

	RegisterServiceListener(serviceRouter ServiceRouter,
		baseTx types.BaseTx) error

	RegisterSingleServiceListener(serviceName string,
		respondHandler ServiceRespondHandler,
		baseTx types.BaseTx) error
}

type ServiceQuery interface {
	QueryDefinition(serviceName string) (ServiceDefinition, error)

	QueryBinding(serviceName string, provider types.AccAddress) (ServiceBinding, error)
	QueryBindings(serviceName string) ([]ServiceBinding, error)

	QueryRequest(requestID string) (RequestService, error)
	QueryRequests(serviceName string, provider types.AccAddress) ([]RequestService, error)
	QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) ([]RequestService, error)

	QueryResponse(requestID string) (ServiceResponse, error)
	QueryResponses(requestContextID string, batchCounter uint64) ([]ServiceResponse, error)

	QueryRequestContext(requestContextID string) (RequestContext, error)
	QueryFees(provider string) (EarnedFees, error)
}

type Service interface {
	types.Module
	ServiceTx
	ServiceQuery
}

type ServiceInvokeHandler func(reqCtxID string, responses string)
type ServiceRespondHandler func(input string) (output string, errMsg string)
type ServiceRouter map[string]ServiceRespondHandler

// RequestService defines a request which contains the detailed request data
type RequestService struct {
	ServiceName                string           `json:"service_name"`
	Provider                   types.AccAddress `json:"provider"`
	Consumer                   types.AccAddress `json:"consumer"`
	Input                      string           `json:"input"`
	ServiceFee                 types.Coins      `json:"service_fee"`
	SuperMode                  bool             `json:"super_mode"`
	RequestHeight              int64            `json:"request_height"`
	ExpirationHeight           int64            `json:"expiration_height"`
	RequestContextID           string           `json:"request_context_id"`
	RequestContextBatchCounter uint64           `json:"request_context_batch_counter"`
}

// ServiceResponse defines a response
type ServiceResponse struct {
	Provider                   types.AccAddress `json:"provider"`
	Consumer                   types.AccAddress `json:"consumer"`
	Output                     string           `json:"output"`
	Error                      string           `json:"error"`
	RequestContextID           string           `json:"request_context_id"`
	RequestContextBatchCounter uint64           `json:"request_context_batch_counter"`
}

type ServiceDefinitionRequest struct {
	types.BaseTx
	ServiceName       string   `json:"service_name"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	AuthorDescription string   `json:"author_description"`
	Schemas           string   `json:"schemas"`
}

// ServiceDefinition represents a service definition
type ServiceDefinition struct {
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Tags              []string         `json:"tags"`
	Author            types.AccAddress `json:"author"`
	AuthorDescription string           `json:"author_description"`
	Schemas           string           `json:"schemas"`
}

type ServiceBindingRequest struct {
	types.BaseTx
	ServiceName string      `json:"service_name"`
	Deposit     types.Coins `json:"deposit"`
	Pricing     string      `json:"pricing"`
}

// UpdateServiceBindingRequest defines a message to update a service binding
type UpdateServiceBindingRequest struct {
	types.BaseTx
	ServiceName string      `json:"service_name"`
	Deposit     types.Coins `json:"deposit"`
	Pricing     string      `json:"pricing"`
}

// ServiceBinding defines a struct for service binding
type ServiceBinding struct {
	ServiceName     string           `json:"service_name"`
	Provider        types.AccAddress `json:"provider"`
	Deposit         types.Coins      `json:"deposit"`
	Pricing         string           `json:"pricing"`
	WithdrawAddress types.AccAddress `json:"withdraw_address"`
	Available       bool             `json:"available"`
	DisabledTime    time.Time        `json:"disabled_time"`
}

type ServiceInvocationRequest struct {
	types.BaseTx
	ServiceName       string      `json:"service_name"`
	Providers         []string    `json:"providers"`
	Input             string      `json:"input"`
	ServiceFeeCap     types.Coins `json:"service_fee_cap"`
	Timeout           int64       `json:"timeout"`
	SuperMode         bool        `json:"super_mode"`
	Repeated          bool        `json:"repeated"`
	RepeatedFrequency uint64      `json:"repeated_frequency"`
	RepeatedTotal     int64       `json:"repeated_total"`
}

// UpdateContextRequest defines a message to update a request context
type UpdateContextRequest struct {
	types.BaseTx
	RequestContextID  string      `json:"request_context_id"`
	Providers         []string    `json:"providers"`
	ServiceFeeCap     types.Coins `json:"service_fee_cap"`
	Timeout           int64       `json:"timeout"`
	RepeatedFrequency uint64      `json:"repeated_frequency"`
	RepeatedTotal     int64       `json:"repeated_total"`
}

// RequestContext defines a context which holds request-related data
type RequestContext struct {
	ServiceName        string             `json:"service_name"`
	Providers          []types.AccAddress `json:"providers"`
	Consumer           types.AccAddress   `json:"consumer"`
	Input              string             `json:"input"`
	ServiceFeeCap      types.Coins        `json:"service_fee_cap"`
	Timeout            int64              `json:"timeout"`
	SuperMode          bool               `json:"super_mode"`
	Repeated           bool               `json:"repeated"`
	RepeatedFrequency  uint64             `json:"repeated_frequency"`
	RepeatedTotal      int64              `json:"repeated_total"`
	BatchCounter       uint64             `json:"batch_counter"`
	BatchRequestCount  uint16             `json:"batch_request_count"`
	BatchResponseCount uint16             `json:"batch_response_count"`
	BatchState         string             `json:"batch_state"`
	State              string             `json:"state"`
	ResponseThreshold  uint16             `json:"response_threshold"`
	ModuleName         string             `json:"module_name"`
}

// EarnedFees defines a struct for the fees earned by the provider
type EarnedFees struct {
	Address types.AccAddress `json:"address"`
	Coins   types.Coins      `json:"coins"`
}

func RequestContextIDToString(reqCtxID []byte) string {
	return hex.EncodeToString(reqCtxID)
}

func RequestContextIDToByte(reqCtxID string) []byte {
	dst, err := hex.DecodeString(reqCtxID)
	if err != nil {
		panic(err)
	}
	return dst
}
