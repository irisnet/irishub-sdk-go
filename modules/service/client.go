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
func (s serviceClient) DefineService(definition sdk.ServiceDefinition, baseTx sdk.BaseTx) (sdk.Result, error) {
	author, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgDefineService{
		Name:              definition.ServiceName,
		Description:       definition.Description,
		Tags:              definition.Tags,
		Author:            author,
		AuthorDescription: definition.AuthorDescription,
		Schemas:           definition.Schemas,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//BindService is responsible for binding a new service definition
func (s serviceClient) BindService(binding sdk.ServiceBinding, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	coins, err := sdk.ParseCoins(binding.Deposit)
	if err != nil {
		return nil, err
	}

	withdrawAddress, err := sdk.AccAddressFromBech32(binding.WithdrawAddr)
	if err != nil {
		return nil, err
	}
	msg := MsgBindService{
		ServiceName:     binding.ServiceName,
		Provider:        provider,
		Deposit:         coins,
		Pricing:         binding.Pricing,
		WithdrawAddress: withdrawAddress,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//InvokeService is responsible for invoke a new service and callback `callback`
func (s serviceClient) InvokeService(invocation sdk.ServiceInvocation,
	baseTx sdk.BaseTx,
	callback sdk.ServiceInvokeHandler) (string, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return "", err
	}

	var ps []sdk.AccAddress
	for _, p := range invocation.Providers {
		provider, err := sdk.AccAddressFromBech32(p)
		if err != nil {
			return "", err
		}
		ps = append(ps, provider)
	}

	coins, err := sdk.ParseCoins(invocation.ServiceFeeCap)
	if err != nil {
		return "", err
	}

	msg := MsgRequestService{
		ServiceName:       invocation.ServiceName,
		Providers:         ps,
		Consumer:          consumer,
		Input:             invocation.Input,
		ServiceFeeCap:     coins,
		Timeout:           invocation.Timeout,
		SuperMode:         invocation.SuperMode,
		Repeated:          invocation.Repeated,
		RepeatedFrequency: invocation.RepeatedFrequency,
		RepeatedTotal:     invocation.RepeatedTotal,
	}

	result, err := s.Broadcast(baseTx, []sdk.Msg{msg})
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
			if ok && request.ServiceName == invocation.ServiceName {
				responses = append(responses, msg.Output)
				callback(requestContextID, msg.Output)
			}
		}
		//cancel subscription
		if !invocation.Repeated {
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

func (s serviceClient) QueryRequest(requestID string) (request sdk.Request, err error) {
	param := QueryRequestParams{RequestID: requestID}
	//TODO
	err = s.Query("custom/service/request", param, &request)
	if err != nil {
		return request, err
	}
	return
}
