package types

type Service interface {
	DefineService(serviceName string,
		description string,
		tags []string,
		authorDescription string,
		schemas string,
		baseTx BaseTx,
	) (Result, error)

	BindService(serviceName string,
		deposit string,
		pricing string,
		withdrawAddr string,
		baseTx BaseTx,
	) (Result, error)

	InvokeService(serviceName string,
		providers []string,
		input string,
		serviceFeeCap string,
		timeout int64,
		superMode bool,
		repeated bool,
		repeatedFrequency uint64,
		repeatedTotal int64,
		baseTx BaseTx,
		callback ServiceInvokeHandler,
	) (requestContextID string, err error)

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
