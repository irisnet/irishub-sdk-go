package service

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type serviceClient struct {
	sdk.AbstractClient
}

func (s serviceClient) DefineService(serviceName string,
	description string,
	tags []string,
	authorDescription string,
	schemas string,
	baseTx sdk.BaseTx,
) (sdk.Result, error) {
	author, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgDefineService{
		Name:              serviceName,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

func (s serviceClient) BindService(serviceName string, deposit string, pricing string,
	withdrawAddr string, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	coins, err := sdk.ParseCoins(deposit)
	if err != nil {
		return nil, err
	}

	withdrawAddress, err := sdk.AccAddressFromBech32(withdrawAddr)
	if err != nil {
		return nil, err
	}
	msg := MsgBindService{
		ServiceName:     serviceName,
		Provider:        provider,
		Deposit:         coins,
		Pricing:         pricing,
		WithdrawAddress: withdrawAddress,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

func (s serviceClient) InvokeService(serviceName string,
	providers []string,
	input string,
	serviceFeeCap string,
	timeout int64,
	superMode bool,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	baseTx sdk.BaseTx,
	callback sdk.ServiceInvokeHandler) (string, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return "", err
	}

	var ps []sdk.AccAddress
	for _, p := range providers {
		provider, err := sdk.AccAddressFromBech32(p)
		if err != nil {
			return "", err
		}
		ps = append(ps, provider)
	}

	coins, err := sdk.ParseCoins(serviceFeeCap)
	if err != nil {
		return "", err
	}

	msg := MsgRequestService{
		ServiceName:       serviceName,
		Providers:         ps,
		Consumer:          consumer,
		Input:             input,
		ServiceFeeCap:     coins,
		Timeout:           timeout,
		SuperMode:         superMode,
		Repeated:          repeated,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
	}

	result, err := s.Broadcast(baseTx, []sdk.Msg{msg})
	if err != nil {
		return "", err
	}

	requestContextID := result.GetTags().ToMap()[TagRequestContextID]

	builder := sdk.NewEventQueryBuilder()
	builder.AddCondition(sdk.ActionKey, sdk.EventValue(TagRespondService)).
		AddCondition(sdk.EventKey(TagConsumer), sdk.EventValue(consumer.String()))

	var subscription sdk.Subscription
	subscription, err = s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		var responses []string
		for _, msg := range tx.Tx.Msgs {
			msg, ok := msg.(MsgRespondService)
			if ok {
				responses = append(responses, msg.Output)
			}
		}
		//TODO responseThreshold ?
		if len(responses) > 0 {
			go func() {
				callback(requestContextID, responses)
			}()
		}
		//cancel subscription
		if !repeated {
			_ = s.Unscribe(subscription)
		}
	})
	if err != nil {
		return "", err
	}

	return requestContextID, nil
}

func (s serviceClient) RegisterInvocationListener(subscriptions []sdk.ServiceRespSubscription, baseTx sdk.BaseTx) error {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return err
	}

	handlerMap := make(map[string]sdk.ServiceRespondHandler)
	for _, subscription := range subscriptions {
		handlerMap[subscription.ServiceName] = subscription.RespondHandler
	}
	_, err = s.SubscribeNewBlock(func(block sdk.EventDataNewBlock) {
		tags := block.ResultEndBlock.Tags.ToMap()
		//TODO add ServiceName to Tags
		reqID := tags[TagRequestID]
		serviceName := tags[TagServiceName]
		if tags[TagProvider] == provider.String() &&
			reqID != "" {
			if handler, ok := handlerMap[serviceName]; ok {
				input := s.QueryRequest(reqID)
				output, errMsg := handler(input)
				msg := MsgRespondService{
					RequestID: reqID,
					Provider:  provider,
					Output:    output,
					Error:     errMsg,
				}
				_, _ = s.Broadcast(baseTx, []sdk.Msg{msg})
			}
		}
	})
	return err
}

func (s serviceClient) RegisterSingleInvocationListener(subscription sdk.ServiceRespSubscription, baseTx sdk.BaseTx) error {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return err
	}
	_, err = s.SubscribeNewBlock(func(block sdk.EventDataNewBlock) {
		tags := block.ResultEndBlock.Tags.ToMap()
		//TODO add ServiceName to Tags
		reqID := tags[TagRequestID]
		serviceName := tags[TagServiceName]
		if tags[TagProvider] == provider.String() &&
			reqID != "" {
			if len(subscription.ServiceName) != 0 &&
				subscription.ServiceName == serviceName {
				input := s.QueryRequest(reqID)
				output, errMsg := subscription.RespondHandler(input)
				msg := MsgRespondService{
					RequestID: reqID,
					Provider:  provider,
					Output:    output,
					Error:     errMsg,
				}
				_, _ = s.Broadcast(baseTx, []sdk.Msg{msg})
			}
		}
	})
	return err
}

//TODO
func (s serviceClient) QueryRequest(requestID string) string {
	return ""
}

func New(ac sdk.AbstractClient) sdk.Service {
	return serviceClient{
		AbstractClient: ac,
	}
}
