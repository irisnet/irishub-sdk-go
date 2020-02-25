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

	requestContextID := result.GetTags().GetValue(TagRequestContextID)

	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.ActionKey, sdk.EventValue(TagRespondService)).
		AddCondition(sdk.EventKey(TagConsumer), sdk.EventValue(consumer.String()))

	var subscription sdk.Subscription
	subscription, err = s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		var responses []string
		for _, msg := range tx.Tx.Msgs {
			msg, ok := msg.(MsgRespondService)
			request := s.QueryRequest(msg.RequestID)
			if ok && request.ServiceName == serviceName {
				responses = append(responses, msg.Output)
				callback(requestContextID, msg.Output)
			}
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
			request := s.QueryRequest(reqID)
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
			request := s.QueryRequest(reqID)
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

//TODO
func (s serviceClient) QueryRequest(requestID string) sdk.Request {
	return sdk.Request{}
}

func New(ac sdk.AbstractClient) sdk.Service {
	return serviceClient{
		AbstractClient: ac,
	}
}
