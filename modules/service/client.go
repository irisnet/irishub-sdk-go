package service

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type serviceClient struct {
	sdk.AbstractClient
}

func New(ac sdk.AbstractClient) sdk.Service {
	return serviceClient{
		AbstractClient: ac,
	}
}

//DefineService is responsible for creating a new service definition
func (s serviceClient) DefineService(request sdk.ServiceDefinitionRequest) (sdk.Result, error) {
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
func (s serviceClient) BindService(request sdk.ServiceBindingRequest) (sdk.Result, error) {
	provider, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	coins, err := sdk.ParseCoins(request.Deposit)
	if err != nil {
		return nil, err
	}

	var withdrawAddress sdk.AccAddress
	if len(request.WithdrawAddr) > 0 {
		withdrawAddress, err = sdk.AccAddressFromBech32(request.WithdrawAddr)
		if err != nil {
			return nil, err
		}
	}

	msg := MsgBindService{
		ServiceName:     request.ServiceName,
		Provider:        provider,
		Deposit:         coins,
		Pricing:         request.Pricing,
		WithdrawAddress: withdrawAddress,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

//InvokeService is responsible for invoke a new service and callback `callback`
func (s serviceClient) InvokeService(request sdk.ServiceInvocationRequest,
	callback sdk.ServiceInvokeHandler) (string, error) {
	consumer, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return "", err
	}

	var ps []sdk.AccAddress
	for _, p := range request.Providers {
		provider, err := sdk.AccAddressFromBech32(p)
		if err != nil {
			return "", err
		}
		ps = append(ps, provider)
	}

	coins, err := sdk.ParseCoins(request.ServiceFeeCap)
	if err != nil {
		return "", err
	}

	msg := MsgRequestService{
		ServiceName:       request.ServiceName,
		Providers:         ps,
		Consumer:          consumer,
		Input:             request.Input,
		ServiceFeeCap:     coins,
		Timeout:           request.Timeout,
		SuperMode:         request.SuperMode,
		Repeated:          request.Repeated,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
	}

	result, err := s.Broadcast(request.BaseTx, []sdk.Msg{msg})
	if err != nil {
		return "", err
	}

	requestContextID := result.GetTags().GetValue(TagRequestContextID)

	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.ActionKey, sdk.EventValue(TagRespondService)).
		AddCondition(sdk.EventKey(TagConsumer), sdk.EventValue(consumer.String()))

	var subscription sdk.Subscription
	subscription, err = s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		var responses []string
		for _, msg := range tx.Tx.Msgs {
			msg, ok := msg.(MsgRespondService)
			request, err := s.QueryRequest(msg.RequestID)
			if err != nil {
				//TODO
				continue
			}
			if ok && request.ServiceName == request.ServiceName {
				responses = append(responses, msg.Output)
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

//RegisterInvocationListener is responsible for registering a group of service handler
func (s serviceClient) RegisterInvocationListener(serviceRouter sdk.ServiceRouter, baseTx sdk.BaseTx) error {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	_, err = s.SubscribeNewBlock(func(block sdk.EventDataNewBlock) {
		reqIDs := block.ResultEndBlock.Tags.GetValues(TagRequestID)
		for _, reqID := range reqIDs {
			request, err := s.QueryRequest(reqID)
			if err != nil {
				//TODO
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

//RegisterSingleInvocationListener is responsible for registering a single service handler
func (s serviceClient) RegisterSingleInvocationListener(serviceName string,
	respondHandler sdk.ServiceRespondHandler,
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
	_, err = s.SubscribeNewBlock(func(block sdk.EventDataNewBlock) {
		reqIDs := block.ResultEndBlock.Tags.GetValues(TagRequestID)
		for _, reqID := range reqIDs {
			request, err := s.QueryRequest(reqID)
			if err != nil {
				//TODO
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
						panic(err)
					}
				}()
			}
		}
	})
	return err
}

// QueryDefinition return a service definition of the specified name
func (s serviceClient) QueryDefinition(serviceName string) (result sdk.ServiceDefinition, err error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}
	err = s.Query("custom/service/definition", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryBinding return the specified service binding
func (s serviceClient) QueryBinding(serviceName string, provider sdk.AccAddress) (result sdk.ServiceBinding, err error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}
	err = s.Query("custom/service/binding", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryBindings returns all bindings of the specified service
func (s serviceClient) QueryBindings(serviceName string) (result []sdk.ServiceBinding, err error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}
	err = s.Query("custom/service/bindings", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryRequest returns  the active request of the specified requestID
func (s serviceClient) QueryRequest(requestID string) (request sdk.Request, err error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}
	err = s.Query("custom/service/request", param, &request)
	if err != nil {
		return request, err
	}
	return
}

// QueryRequest returns all the active requests of the specified service binding
func (s serviceClient) QueryRequests(serviceName string, provider sdk.AccAddress) (result []sdk.Request, err error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}
	err = s.Query("custom/service/requests", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryRequestsByReqCtx returns all requests of the specified request context ID and batch counter
func (s serviceClient) QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) (result []sdk.Request, err error) {
	param := struct {
		RequestContextID []byte
		BatchCounter     uint64
	}{
		RequestContextID: []byte(requestContextID),
		BatchCounter:     batchCounter,
	}
	err = s.Query("custom/service/requests_by_ctx", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryResponse returns a response with the speicified request ID
func (s serviceClient) QueryResponse(requestID string) (result sdk.Response, err error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}
	err = s.Query("custom/service/response", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryResponses returns all responses of the specified request context and batch counter
func (s serviceClient) QueryResponses(requestContextID string, batchCounter uint64) (result []sdk.Response, err error) {
	param := struct {
		RequestContextID []byte
		BatchCounter     uint64
	}{
		RequestContextID: []byte(requestContextID),
		BatchCounter:     batchCounter,
	}
	err = s.Query("custom/service/responses", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryRequestContext return the specified request context
func (s serviceClient) QueryRequestContext(requestContextID string) (result sdk.RequestContext, err error) {
	param := struct {
		RequestContextID []byte
	}{
		RequestContextID: []byte(requestContextID),
	}
	err = s.Query("custom/service/context", param, &result)
	if err != nil {
		return result, err
	}
	return
}

//QueryFees return the earned fees for a provider
func (s serviceClient) QueryFees(provider sdk.AccAddress) (result sdk.EarnedFees, err error) {
	param := struct {
		Address sdk.AccAddress
	}{
		Address: provider,
	}
	err = s.Query("custom/service/fees", param, &result)
	if err != nil {
		return result, err
	}
	return
}
