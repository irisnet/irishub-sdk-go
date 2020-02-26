package types

type Service interface {
	DefineService(definition ServiceDefinition, baseTx BaseTx) (Result, error)

	BindService(binding ServiceBinding, baseTx BaseTx) (Result, error)

	InvokeService(invocation ServiceInvocation, baseTx BaseTx,
		callback ServiceInvokeHandler) (requestContextID string, err error)

	RegisterInvocationListener(serviceRouter ServiceRouter,
		baseTx BaseTx) error

	RegisterSingleInvocationListener(serviceName string,
		respondHandler ServiceRespondHandler,
		baseTx BaseTx) error
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

type ServiceDefinition struct {
	ServiceName       string   `json:"service_name"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	AuthorDescription string   `json:"author_description"`
	Schemas           string   `json:"schemas"`
}

type ServiceBinding struct {
	ServiceName  string `json:"service_name"`
	Deposit      string `json:"deposit"`
	Pricing      string `json:"pricing"`
	WithdrawAddr string `json:"withdraw_addr"`
}

type ServiceInvocation struct {
	ServiceName       string   `json:"service_name"`
	Providers         []string `json:"providers"`
	Input             string   `json:"input"`
	ServiceFeeCap     string   `json:"service_fee_cap"`
	Timeout           int64    `json:"timeout"`
	SuperMode         bool     `json:"super_mode"`
	Repeated          bool     `json:"repeated"`
	RepeatedFrequency uint64   `json:"repeated_frequency"`
	RepeatedTotal     int64    `json:"repeated_total"`
}
