package types

import (
	"time"
)

type ServiceTx interface {
	DefineService(request ServiceDefinitionRequest) (Result, error)

	BindService(request ServiceBindingRequest) (Result, error)

	UpdateServiceBinding(request UpdateServiceBindingRequest) (Result, error)

	InvokeService(request ServiceInvocationRequest,
		callback ServiceInvokeHandler) (requestContextID string, err error)

	SetWithdrawAddress(serviceName string, withdrawAddress string, baseTx BaseTx) (Result, error)

	DisableService(serviceName string, baseTx BaseTx) (Result, error)

	EnableService(serviceName string, deposit Coins, baseTx BaseTx) (Result, error)

	RefundServiceDeposit(serviceName string, baseTx BaseTx) (Result, error)

	PauseRequestContext(requestContextID string, baseTx BaseTx) (Result, error)

	StartRequestContext(requestContextID string, baseTx BaseTx) (Result, error)

	KillRequestContext(requestContextID string, baseTx BaseTx) (Result, error)

	UpdateRequestContext(request UpdateContextRequest) (Result, error)

	WithdrawEarnedFees(baseTx BaseTx) (Result, error)

	WithdrawTax(destAddress string, amount Coins, baseTx BaseTx) (Result, error)

	RegisterInvocationListener(serviceRouter ServiceRouter,
		baseTx BaseTx) error

	RegisterSingleInvocationListener(serviceName string,
		respondHandler ServiceRespondHandler,
		baseTx BaseTx) error
}

type ServiceQuery interface {
	QueryDefinition(serviceName string) (ServiceDefinition, error)

	QueryBinding(serviceName string, provider AccAddress) (ServiceBinding, error)
	QueryBindings(serviceName string) ([]ServiceBinding, error)

	QueryRequest(requestID string) (Request, error)
	QueryRequests(serviceName string, provider AccAddress) ([]Request, error)
	QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) ([]Request, error)

	QueryResponse(requestID string) (Response, error)
	QueryResponses(requestContextID string, batchCounter uint64) ([]Response, error)

	QueryRequestContext(requestContextID string) (RequestContext, error)
	QueryFees(provider AccAddress) (EarnedFees, error)
}

type Service interface {
	ServiceTx
	ServiceQuery
}

type ServiceInvokeHandler func(reqCtxID string, responses string)
type ServiceRespondHandler func(input string) (output string, errMsg string)
type ServiceRouter map[string]ServiceRespondHandler

// Request defines a request which contains the detailed request data
type Request struct {
	ServiceName                string     `json:"service_name"`
	Provider                   AccAddress `json:"provider"`
	Consumer                   AccAddress `json:"consumer"`
	Input                      string     `json:"input"`
	ServiceFee                 Coins      `json:"service_fee"`
	SuperMode                  bool       `json:"super_mode"`
	RequestHeight              int64      `json:"request_height"`
	ExpirationHeight           int64      `json:"expiration_height"`
	RequestContextID           []byte     `json:"request_context_id"`
	RequestContextBatchCounter uint64     `json:"request_context_batch_counter"`
}

// Response defines a response
type Response struct {
	Provider                   AccAddress `json:"provider"`
	Consumer                   AccAddress `json:"consumer"`
	Output                     string     `json:"output"`
	Error                      string     `json:"error"`
	RequestContextID           []byte     `json:"request_context_id"`
	RequestContextBatchCounter uint64     `json:"request_context_batch_counter"`
}

type ServiceDefinitionRequest struct {
	BaseTx
	ServiceName       string   `json:"service_name"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	AuthorDescription string   `json:"author_description"`
	Schemas           string   `json:"schemas"`
}

// ServiceDefinition represents a service definition
type ServiceDefinition struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Tags              []string   `json:"tags"`
	Author            AccAddress `json:"author"`
	AuthorDescription string     `json:"author_description"`
	Schemas           string     `json:"schemas"`
}

type ServiceBindingRequest struct {
	BaseTx
	ServiceName  string `json:"service_name"`
	Deposit      Coins  `json:"deposit"`
	Pricing      string `json:"pricing"`
	WithdrawAddr string `json:"withdraw_addr"`
}

// UpdateServiceBindingRequest defines a message to update a service binding
type UpdateServiceBindingRequest struct {
	BaseTx
	ServiceName string `json:"service_name"`
	Deposit     Coins  `json:"deposit"`
	Pricing     string `json:"pricing"`
}

// ServiceBinding defines a struct for service binding
type ServiceBinding struct {
	ServiceName     string     `json:"service_name"`
	Provider        AccAddress `json:"provider"`
	Deposit         Coins      `json:"deposit"`
	Pricing         string     `json:"pricing"`
	WithdrawAddress AccAddress `json:"withdraw_address"`
	Available       bool       `json:"available"`
	DisabledTime    time.Time  `json:"disabled_time"`
}

type ServiceInvocationRequest struct {
	BaseTx
	ServiceName       string   `json:"service_name"`
	Providers         []string `json:"providers"`
	Input             string   `json:"input"`
	ServiceFeeCap     Coins    `json:"service_fee_cap"`
	Timeout           int64    `json:"timeout"`
	SuperMode         bool     `json:"super_mode"`
	Repeated          bool     `json:"repeated"`
	RepeatedFrequency uint64   `json:"repeated_frequency"`
	RepeatedTotal     int64    `json:"repeated_total"`
}

// UpdateContextRequest defines a message to update a request context
type UpdateContextRequest struct {
	BaseTx
	RequestContextID  string   `json:"request_context_id"`
	Providers         []string `json:"providers"`
	ServiceFeeCap     Coins    `json:"service_fee_cap"`
	Timeout           int64    `json:"timeout"`
	RepeatedFrequency uint64   `json:"repeated_frequency"`
	RepeatedTotal     int64    `json:"repeated_total"`
}

// RequestContext defines a context which holds request-related data
type RequestContext struct {
	ServiceName        string       `json:"service_name"`
	Providers          []AccAddress `json:"providers"`
	Consumer           AccAddress   `json:"consumer"`
	Input              string       `json:"input"`
	ServiceFeeCap      Coins        `json:"service_fee_cap"`
	Timeout            int64        `json:"timeout"`
	SuperMode          bool         `json:"super_mode"`
	Repeated           bool         `json:"repeated"`
	RepeatedFrequency  uint64       `json:"repeated_frequency"`
	RepeatedTotal      int64        `json:"repeated_total"`
	BatchCounter       uint64       `json:"batch_counter"`
	BatchRequestCount  uint16       `json:"batch_request_count"`
	BatchResponseCount uint16       `json:"batch_response_count"`
	BatchState         byte         `json:"batch_state"`
	State              byte         `json:"state"`
	ResponseThreshold  uint16       `json:"response_threshold"`
	ModuleName         string       `json:"module_name"`
}

// EarnedFees defines a struct for the fees earned by the provider
type EarnedFees struct {
	Address AccAddress `json:"address"`
	Coins   Coins      `json:"coins"`
}
