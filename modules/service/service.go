package service

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type serviceClient struct {
	sdk.AbstractClient
	*log.Logger
}

func (s serviceClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (s serviceClient) Name() string {
	return ModuleName
}

func New(ac sdk.AbstractClient) rpc.Service {
	return serviceClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

//DefineService is responsible for creating a new service definition
func (s serviceClient) DefineService(request rpc.ServiceDefinitionRequest) (sdk.Result, error) {
	author, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgDefineService{
		Name:              request.ServiceName,
		Description:       request.Description,
		Tags:              request.Tags,
		Author:            author,
		AuthorDescription: request.AuthorDescription,
		Schemas:           request.Schemas,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

//BindService is responsible for binding a new service definition
func (s serviceClient) BindService(request rpc.ServiceBindingRequest) (sdk.Result, error) {
	provider, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	msg := MsgBindService{
		ServiceName: request.ServiceName,
		Provider:    provider,
		Deposit:     request.Deposit,
		Pricing:     request.Pricing,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

//UpdateServiceBinding updates the specified service binding
func (s serviceClient) UpdateServiceBinding(request rpc.UpdateServiceBindingRequest) (sdk.Result, error) {
	provider, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgUpdateServiceBinding{
		ServiceName: request.ServiceName,
		Provider:    provider,
		Deposit:     request.Deposit,
		Pricing:     request.Pricing,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

// DisableService disables the specified service binding
func (s serviceClient) DisableService(serviceName string, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgDisableService{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// EnableService enables the specified service binding
func (s serviceClient) EnableService(serviceName string, deposit sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgEnableService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//InvokeService is responsible for invoke a new service and callback `callback`
func (s serviceClient) InvokeService(request rpc.ServiceInvocationRequest,
	callback rpc.ServiceInvokeHandler) (string, error) {
	consumer, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return "", err
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		provider := sdk.MustAccAddressFromBech32(provider)
		providers = append(providers, provider)
	}

	msg := MsgRequestService{
		ServiceName:       request.ServiceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             request.Input,
		ServiceFeeCap:     request.ServiceFeeCap,
		Timeout:           request.Timeout,
		SuperMode:         request.SuperMode,
		Repeated:          request.Repeated,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
	}

	//mode must be set to commit
	request.BaseTx.Mode = sdk.Commit

	result, err := s.Broadcast(request.BaseTx, []sdk.Msg{msg})
	if err != nil {
		return "", err
	}

	requestContextID := result.GetTags().GetValue(tagRequestContextID)
	if callback == nil {
		return requestContextID, nil
	}
	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.ActionKey, sdk.EventValue(tagRespondService)).
		AddCondition(sdk.EventKey(tagConsumer), sdk.EventValue(consumer.String())).
		AddCondition(sdk.EventKey(tagServiceName), sdk.EventValue(request.ServiceName))

	var subscription sdk.Subscription
	subscription, err = s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		s.Debug().Str("tx_hash", tx.Hash).Int64("height", tx.Height).
			Msg("consumer received response transaction sent by provider")
		for _, msg := range tx.Tx.Msgs {
			msg, ok := msg.(MsgRespondService)
			if ok {
				callback(requestContextID, msg.Output)
			}
		}
		//cancel subscription
		if !request.Repeated {
			_ = s.Unscribe(subscription)
		}
	})
	if err != nil {
		return "", err
	}

	return requestContextID, nil
}

// SetWithdrawAddress sets a new withdrawal address for the specified service binding
func (s serviceClient) SetWithdrawAddress(serviceName string, withdrawAddress string, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	withdrawAddr := sdk.MustAccAddressFromBech32(withdrawAddress)
	msg := MsgSetWithdrawAddress{
		ServiceName:     serviceName,
		Provider:        provider,
		WithdrawAddress: withdrawAddr,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// RefundServiceDeposit refunds the deposit from the specified service binding
func (s serviceClient) RefundServiceDeposit(serviceName string, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// StartRequestContext starts the specified request context
func (s serviceClient) StartRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.Result, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgStartRequestContext{
		RequestContextID: rpc.RequestContextIDToByte(requestContextID),
		Consumer:         consumer,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// PauseRequestContext suspends the specified request context
func (s serviceClient) PauseRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.Result, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgPauseRequestContext{
		RequestContextID: rpc.RequestContextIDToByte(requestContextID),
		Consumer:         consumer,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// KillRequestContext terminates the specified request context
func (s serviceClient) KillRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.Result, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgKillRequestContext{
		RequestContextID: rpc.RequestContextIDToByte(requestContextID),
		Consumer:         consumer,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// UpdateRequestContext updates the specified request context
func (s serviceClient) UpdateRequestContext(request rpc.UpdateContextRequest) (sdk.Result, error) {
	consumer, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	var providers []sdk.AccAddress
	for _, p := range request.Providers {
		provider := sdk.MustAccAddressFromBech32(p)
		providers = append(providers, provider)
	}

	msg := MsgUpdateRequestContext{
		RequestContextID:  rpc.RequestContextIDToByte(request.RequestContextID),
		Providers:         providers,
		ServiceFeeCap:     request.ServiceFeeCap,
		Timeout:           request.Timeout,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		Consumer:          consumer,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

// WithdrawEarnedFees withdraws the earned fees to the specified provider
func (s serviceClient) WithdrawEarnedFees(baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	msg := MsgWithdrawEarnedFees{
		Provider: provider,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// WithdrawTax withdraws the service tax to the speicified destination address by the trustee
func (s serviceClient) WithdrawTax(destAddress string, amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error) {
	trustee, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	receipt := sdk.MustAccAddressFromBech32(destAddress)
	msg := MsgWithdrawTax{
		Trustee:     trustee,
		DestAddress: receipt,
		Amount:      amount,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//RegisterServiceListener is responsible for registering a group of service handler
func (s serviceClient) RegisterServiceListener(serviceRouter rpc.ServiceRouter, baseTx sdk.BaseTx) error {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.EventKey(tagProvider), sdk.EventValue(provider.String()))
	_, err = s.SubscribeNewBlockWithQuery(builder, func(block sdk.EventDataNewBlock) {
		s.Debug().Int64("height", block.Block.Height).Msg("received block")
		reqIDs := block.ResultEndBlock.Tags.GetValues(tagRequestID)
		for _, reqID := range reqIDs {
			request, err := s.QueryRequest(reqID)
			if err != nil {
				s.Err(err).Str("requestID", reqID).Msg("service request don't exist")
				continue
			}
			if handler, ok := serviceRouter[request.ServiceName]; ok && provider.Equals(request.Provider) {
				output, errMsg := handler(request.Input)
				msg := MsgRespondService{
					RequestID: reqID,
					Provider:  provider,
					Output:    output,
					Error:     errMsg,
				}
				go func() {
					if _, err = s.Broadcast(baseTx, []sdk.Msg{msg}); err != nil {
						panic(err)
					}
				}()
			}
		}
	})
	return err
}

//RegisterSingleServiceListener is responsible for registering a single service handler
func (s serviceClient) RegisterSingleServiceListener(serviceName string,
	respondHandler rpc.ServiceRespondHandler,
	baseTx sdk.BaseTx) error {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	// TODO user will don't received any event from tendermint when block_result has the same key in tag
	//builder := sdk.NewEventQueryBuilder().
	//	AddCondition(sdk.EventKey(tagProvider), sdk.EventValue(provider.String())).
	//	AddCondition(sdk.EventKey(tagServiceName), sdk.EventValue(serviceName))
	_, err = s.SubscribeNewBlockWithQuery(nil, func(block sdk.EventDataNewBlock) {
		reqIDs := block.ResultEndBlock.Tags.GetValues(tagRequestID)
		s.Debug().Int64("height", block.Block.Height).Msg("received block")
		for _, reqID := range reqIDs {
			request, err := s.QueryRequest(reqID)
			if err != nil {
				s.Err(err).Str("requestID", reqID).Msg("service request don't exist")
				continue
			}
			if provider.Equals(request.Provider) && request.ServiceName == serviceName {
				output, errMsg := respondHandler(request.Input)
				msg := MsgRespondService{
					RequestID: reqID,
					Provider:  provider,
					Output:    output,
					Error:     errMsg,
				}
				go func() {
					if _, err = s.Broadcast(baseTx, []sdk.Msg{msg}); err != nil {
						s.Err(err).Str("requestID", reqID).Msg("provider respond failed")
					}
				}()
			}
		}
	})
	return err
}

// QueryDefinition return a service definition of the specified name
func (s serviceClient) QueryDefinition(serviceName string) (rpc.ServiceDefinition, error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}

	var definition serviceDefinition
	if err := s.QueryWithResponse("custom/service/definition", param, &definition); err != nil {
		return rpc.ServiceDefinition{}, err
	}
	return definition.Convert().(rpc.ServiceDefinition), nil
}

// QueryBinding return the specified service binding
func (s serviceClient) QueryBinding(serviceName string, provider sdk.AccAddress) (rpc.ServiceBinding, error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}

	var binding serviceBinding
	if err := s.QueryWithResponse("custom/service/binding", param, &binding); err != nil {
		return rpc.ServiceBinding{}, err
	}
	return binding.Convert().(rpc.ServiceBinding), nil
}

// QueryBindings returns all bindings of the specified service
func (s serviceClient) QueryBindings(serviceName string) ([]rpc.ServiceBinding, error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}

	var bindings serviceBindings
	if err := s.QueryWithResponse("custom/service/bindings", param, &bindings); err != nil {
		return nil, err
	}
	return bindings.Convert().([]rpc.ServiceBinding), nil
}

// QueryRequest returns  the active request of the specified requestID
func (s serviceClient) QueryRequest(requestID string) (rpc.RequestService, error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}

	var request request
	if err := s.QueryWithResponse("custom/service/request", param, &request); err != nil {
		return rpc.RequestService{}, err
	}
	return request.Convert().(rpc.RequestService), nil
}

// QueryRequest returns all the active requests of the specified service binding
func (s serviceClient) QueryRequests(serviceName string, provider sdk.AccAddress) ([]rpc.RequestService, error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}

	var rs requests
	if err := s.QueryWithResponse("custom/service/requests", param, &rs); err != nil {
		return nil, err
	}
	return rs.Convert().([]rpc.RequestService), nil
}

// QueryRequestsByReqCtx returns all requests of the specified request context ID and batch counter
func (s serviceClient) QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) ([]rpc.RequestService, error) {
	param := struct {
		RequestContextID []byte
		BatchCounter     uint64
	}{
		RequestContextID: rpc.RequestContextIDToByte(requestContextID),
		BatchCounter:     batchCounter,
	}

	var rs requests
	if err := s.QueryWithResponse("custom/service/requests_by_ctx", param, &rs); err != nil {
		return nil, err
	}
	return rs.Convert().([]rpc.RequestService), nil
}

// QueryResponse returns a response with the speicified request ID
func (s serviceClient) QueryResponse(requestID string) (rpc.ServiceResponse, error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}

	var response response
	if err := s.QueryWithResponse("custom/service/response", param, &response); err != nil {
		return rpc.ServiceResponse{}, err
	}
	return response.Convert().(rpc.ServiceResponse), nil
}

// QueryResponses returns all responses of the specified request context and batch counter
func (s serviceClient) QueryResponses(requestContextID string, batchCounter uint64) ([]rpc.ServiceResponse, error) {
	param := struct {
		RequestContextID []byte
		BatchCounter     uint64
	}{
		RequestContextID: rpc.RequestContextIDToByte(requestContextID),
		BatchCounter:     batchCounter,
	}
	var rs responses
	if err := s.QueryWithResponse("custom/service/responses", param, &rs); err != nil {
		return nil, err
	}
	return rs.Convert().([]rpc.ServiceResponse), nil
}

// QueryRequestContext return the specified request context
func (s serviceClient) QueryRequestContext(requestContextID string) (rpc.RequestContext, error) {
	param := struct {
		RequestContextID []byte
	}{
		RequestContextID: rpc.RequestContextIDToByte(requestContextID),
	}

	var reqCtx requestContext
	if err := s.QueryWithResponse("custom/service/context", param, &reqCtx); err != nil {
		return rpc.RequestContext{}, err
	}
	return reqCtx.Convert().(rpc.RequestContext), nil
}

//QueryFees return the earned fees for a provider
func (s serviceClient) QueryFees(provider string) (rpc.EarnedFees, error) {
	address, err := sdk.AccAddressFromBech32(provider)
	if err != nil {
		return rpc.EarnedFees{}, err
	}

	param := struct {
		Address sdk.AccAddress
	}{
		Address: address,
	}

	var fee earnedFees

	if err := s.QueryWithResponse("custom/service/fees", param, &fee); err != nil {
		return rpc.EarnedFees{}, err
	}
	return fee.Convert().(rpc.EarnedFees), nil
}
