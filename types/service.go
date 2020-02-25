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
	) (string, error)

	RegisterInvocationListener(subscriptions []ServiceRespSubscription,
		baseTx BaseTx) error

	RegisterSingleInvocationListener(subscription ServiceRespSubscription,
		baseTx BaseTx) error
}

type ServiceInvokeHandler func(reqCtxID string, responses []string)
type ServiceRespondHandler func(input string) (output string, errMsg string)
type ServiceRespSubscription struct {
	ServiceName    string
	RespondHandler ServiceRespondHandler
}
