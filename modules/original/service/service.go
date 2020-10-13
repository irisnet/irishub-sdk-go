package service

import (
	"encoding/json"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/tendermint/tendermint/libs/bytes"
	"strings"

	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type serviceClient struct {
	original.BaseClient
	*log.Logger
}

func (s serviceClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (s serviceClient) Name() string {
	return ModuleName
}

func Create(ac original.BaseClient) rpc.Service {
	return serviceClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

//DefineService is responsible for creating a new service definition
func (s serviceClient) DefineService(request rpc.ServiceDefinitionRequest, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	author, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}
	msg := MsgDefineService{
		Name:              request.ServiceName,
		Description:       request.Description,
		Tags:              request.Tags,
		Author:            author,
		AuthorDescription: request.AuthorDescription,
		Schemas:           request.Schemas,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

//BindService is responsible for binding a new service definition
func (s serviceClient) BindService(request rpc.ServiceBindingRequest, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	//amt, err := s.ToMinCoin(request.Deposit...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgBindService{
		ServiceName: request.ServiceName,
		Provider:    provider,
		//Deposit:     amt,
		Pricing:     request.Pricing,
		MinRespTime: request.MinRespTime,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

//UpdateServiceBinding updates the specified service binding
func (s serviceClient) UpdateServiceBinding(request rpc.ServiceBindingUpdateRequest, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	//amt, err := s.ToMinCoin(request.Deposit...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgUpdateServiceBinding{
		ServiceName: request.ServiceName,
		Provider:    provider,
		//Deposit:     amt,
		Pricing: request.Pricing,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// DisableServiceBinding disables the specified service binding
func (s serviceClient) DisableServiceBinding(serviceName string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}
	msg := MsgDisableServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// EnableServiceBinding enables the specified service binding
func (s serviceClient) EnableServiceBinding(serviceName string, deposit original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	//amt, err := s.ToMinCoin(deposit...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgEnableServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		//Deposit:     amt,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

//InvokeService is responsible for invoke a new service and callback `handler`
func (s serviceClient) InvokeService(request rpc.ServiceInvocationRequest, baseTx original.BaseTx) (string, original.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return "", original.Wrap(err)
	}

	var providers []original.AccAddress
	for _, provider := range request.Providers {
		p, err := original.AccAddressFromBech32(provider)
		if err != nil {
			return "", original.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return "", original.Wrap(err)
	}

	msg := MsgCallService{
		ServiceName:       request.ServiceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             request.Input,
		ServiceFeeCap:     amt,
		Timeout:           request.Timeout,
		SuperMode:         request.SuperMode,
		Repeated:          request.Repeated,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
	}

	//mode must be set to commit
	baseTx.Mode = original.Commit

	result, err := s.BuildAndSend([]original.Msg{msg}, baseTx)
	if err != nil {
		return "", original.Wrap(err)
	}

	reqCtxID, e := result.Events.GetValue(original.EventTypeMessage, attributeKeyRequestContextID)
	if e != nil {
		return reqCtxID, original.Wrap(e)
	}

	_, err = s.SubscribeServiceResponse(reqCtxID, request.Callback)
	return reqCtxID, original.Wrap(err)
}

func (s serviceClient) SubscribeServiceResponse(reqCtxID string,
	callback rpc.ServiceInvokeCallback) (subscription original.Subscription, err original.Error) {
	if len(reqCtxID) == 0 {
		return subscription, original.Wrapf("reqCtxID %s should not be empty", reqCtxID)
	}

	builder := original.NewEventQueryBuilder().
		AddCondition(original.NewCond(original.EventTypeMessage, attributeKeyRequestContextID).
			EQ(original.EventValue(reqCtxID)))

	return s.SubscribeTx(builder, func(tx original.EventDataTx) {
		s.Debug().
			Str("tx_hash", tx.Hash).
			Int64("height", tx.Height).
			Str("reqCtxID", reqCtxID).
			Msg("consumer received response transaction sent by provider")
		for _, msg := range tx.Tx.Msgs {
			msg, ok := msg.(MsgRespondService)
			if ok {
				reqCtxID2, _, _, _, err := splitRequestID(msg.RequestID.String())
				if err != nil {
					s.Err(err).
						Str("requestID", msg.RequestID.String()).
						Msg("invalid requestID")
					continue
				}
				if reqCtxID2.String() == strings.ToUpper(reqCtxID) {
					callback(reqCtxID, msg.RequestID.String(), msg.Output)
				}
			}
		}
		reqCtx, err := s.QueryRequestContext(reqCtxID)
		if err != nil || reqCtx.State == 0 {
			_ = s.Unsubscribe(subscription)
		}
	})
}

// SetWithdrawAddress sets a new withdrawal address for the specified service binding
func (s serviceClient) SetWithdrawAddress(withdrawAddress string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	withdrawAddr, err := original.AccAddressFromBech32(withdrawAddress)
	if err != nil {
		return original.ResultTx{}, original.Wrapf("%s invalid address", withdrawAddress)
	}
	msg := MsgSetWithdrawAddress{
		Provider:        provider,
		WithdrawAddress: withdrawAddr,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// RefundServiceDeposit refunds the deposit from the specified service binding
func (s serviceClient) RefundServiceDeposit(serviceName string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}
	msg := MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// StartRequestContext starts the specified request context
func (s serviceClient) StartRequestContext(requestContextID string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}
	msg := MsgStartRequestContext{
		RequestContextID: hexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// PauseRequestContext suspends the specified request context
func (s serviceClient) PauseRequestContext(requestContextID string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}
	msg := MsgPauseRequestContext{
		RequestContextID: hexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// KillRequestContext terminates the specified request context
func (s serviceClient) KillRequestContext(requestContextID string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}
	msg := MsgKillRequestContext{
		RequestContextID: hexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// UpdateRequestContext updates the specified request context
func (s serviceClient) UpdateRequestContext(request rpc.UpdateContextRequest, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	var providers []original.AccAddress
	for _, provider := range request.Providers {
		p, err := original.AccAddressFromBech32(provider)
		if err != nil {
			return original.ResultTx{}, original.Wrap(err)
		}
		providers = append(providers, p)
	}

	//amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgUpdateRequestContext{
		RequestContextID: hexBytesFrom(request.RequestContextID),
		Providers:        providers,
		//ServiceFeeCap:     amt,
		Timeout:           request.Timeout,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		Consumer:          consumer,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// WithdrawEarnedFees withdraws the earned fees to the specified provider
func (s serviceClient) WithdrawEarnedFees(baseTx original.BaseTx) (original.ResultTx, original.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgWithdrawEarnedFees{
		Provider: provider,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

// WithdrawTax withdraws the service tax to the speicified destination address by the trustee
func (s serviceClient) WithdrawTax(destAddress string, amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	trustee, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	receipt, err := original.AccAddressFromBech32(destAddress)
	if err != nil {
		return original.ResultTx{}, original.Wrapf("%s invalid address", destAddress)
	}

	//amt, err := s.ToMinCoin(amount...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgWithdrawTax{
		Trustee:     trustee,
		DestAddress: receipt,
		//Amount:      amt,
	}
	return s.BuildAndSend([]original.Msg{msg}, baseTx)
}

//SubscribeServiceRequest is responsible for registering a group of service handler
func (s serviceClient) SubscribeServiceRequest(serviceName string, callback rpc.ServiceRespondCallback,
	baseTx original.BaseTx) (subscription original.Subscription, err original.Error) {
	provider, e := s.QueryAddress(baseTx.From)
	if e != nil {
		return original.Subscription{}, original.Wrap(e)
	}

	builder := original.NewEventQueryBuilder().
		AddCondition(original.NewCond(eventTypeNewBatchRequestProvider, attributeKeyProvider).
			EQ(original.EventValue(provider.String()))).
		AddCondition(original.NewCond(eventTypeNewBatchRequest, attributeKeyServiceName).
			EQ(original.EventValue(serviceName)),
		)
	return s.SubscribeNewBlock(builder, func(block original.EventDataNewBlock) {
		msgs := s.GenServiceResponseMsgs(block.ResultEndBlock.Events, serviceName, provider, callback)
		if _, err = s.SendMsgBatch(msgs, baseTx); err != nil {
			s.Logger.Err(err).Msg("provider respond failed")
		}
	})
}

// QueryDefinition return a service definition of the specified name
func (s serviceClient) QueryDefinition(serviceName string) (rpc.ServiceDefinition, original.Error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}

	var definition serviceDefinition
	if err := s.QueryWithResponse("custom/service/definition", param, &definition); err != nil {
		return rpc.ServiceDefinition{}, original.Wrap(err)
	}
	return definition.Convert().(rpc.ServiceDefinition), nil
}

// QueryBinding return the specified service binding
func (s serviceClient) QueryBinding(serviceName string, provider original.AccAddress) (rpc.ServiceBinding, original.Error) {
	param := struct {
		ServiceName string
		Provider    original.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}

	var binding serviceBinding
	if err := s.QueryWithResponse("custom/service/binding", param, &binding); err != nil {
		return rpc.ServiceBinding{}, original.Wrap(err)
	}
	return binding.Convert().(rpc.ServiceBinding), nil
}

// QueryBindings returns all bindings of the specified service
func (s serviceClient) QueryBindings(serviceName string) ([]rpc.ServiceBinding, original.Error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}

	var bindings serviceBindings
	if err := s.QueryWithResponse("custom/service/bindings", param, &bindings); err != nil {
		return nil, original.Wrap(err)
	}
	return bindings.Convert().([]rpc.ServiceBinding), nil
}

// QueryRequest returns  the active request of the specified requestID
func (s serviceClient) QueryRequest(requestID string) (rpc.ServiceRequest, original.Error) {
	param := struct {
		RequestID []byte
	}{
		RequestID: hexBytesFrom(requestID),
	}

	var request request
	if err := s.QueryWithResponse("custom/service/request", param, &request); request.Empty() {
		request, err = s.queryRequestByTxQuery(requestID)
		if err != nil {
			return rpc.ServiceRequest{}, original.Wrap(err)
		}
	}
	return request.Convert().(rpc.ServiceRequest), nil
}

// QueryRequest returns all the active requests of the specified service binding
func (s serviceClient) QueryRequests(serviceName string, provider original.AccAddress) ([]rpc.ServiceRequest, original.Error) {
	param := struct {
		ServiceName string
		Provider    original.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}

	var rs requests
	if err := s.QueryWithResponse("custom/service/requests", param, &rs); err != nil {
		return nil, original.Wrap(err)
	}
	return rs.Convert().([]rpc.ServiceRequest), nil
}

// QueryRequestsByReqCtx returns all requests of the specified request context ID and batch counter
func (s serviceClient) QueryRequestsByReqCtx(reqCtxID string, batchCounter uint64) ([]rpc.ServiceRequest, original.Error) {
	param := struct {
		RequestContextID bytes.HexBytes
		BatchCounter     uint64
	}{
		RequestContextID: hexBytesFrom(reqCtxID),
		BatchCounter:     batchCounter,
	}

	var rs requests
	if err := s.QueryWithResponse("custom/service/requests_by_ctx", param, &rs); err != nil {
		return nil, original.Wrap(err)
	}
	return rs.Convert().([]rpc.ServiceRequest), nil
}

// QueryResponse returns a response with the speicified request ID
func (s serviceClient) QueryResponse(requestID string) (rpc.ServiceResponse, original.Error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}

	var response response
	if err := s.QueryWithResponse("custom/service/response", param, &response); response.Empty() {
		response, err = s.queryResponseByTxQuery(requestID)
		if err != nil {
			return rpc.ServiceResponse{}, original.Wrap(nil)
		}
	}
	return response.Convert().(rpc.ServiceResponse), nil
}

// QueryResponses returns all responses of the specified request context and batch counter
func (s serviceClient) QueryResponses(reqCtxID string, batchCounter uint64) ([]rpc.ServiceResponse, original.Error) {
	param := struct {
		RequestContextID bytes.HexBytes
		BatchCounter     uint64
	}{
		RequestContextID: hexBytesFrom(reqCtxID),
		BatchCounter:     batchCounter,
	}
	var rs responses
	if err := s.QueryWithResponse("custom/service/responses", param, &rs); err != nil {
		return nil, original.Wrap(err)
	}
	return rs.Convert().([]rpc.ServiceResponse), nil
}

// QueryRequestContext return the specified request context
func (s serviceClient) QueryRequestContext(reqCtxID string) (rpc.RequestContext, original.Error) {
	param := struct {
		RequestContextID bytes.HexBytes
	}{
		RequestContextID: hexBytesFrom(reqCtxID),
	}

	var reqCtx requestContext
	if err := s.QueryWithResponse("custom/service/context", param, &reqCtx); err != nil {
		return rpc.RequestContext{}, original.Wrap(err)
	}
	return reqCtx.Convert().(rpc.RequestContext), nil
}

//QueryFees return the earned fees for a provider
func (s serviceClient) QueryFees(provider string) (original.Coins, original.Error) {
	address, err := original.AccAddressFromBech32(provider)
	if err != nil {
		return nil, original.Wrap(err)
	}

	param := struct {
		Address original.AccAddress
	}{
		Address: address,
	}

	bz, e := s.Query("custom/service/fees", param)
	if e != nil {
		return nil, original.Wrap(err)
	}

	var fee original.Coins
	if err := cdc.UnmarshalJSON(bz, &fee); err != nil {
		return nil, original.Wrap(err)
	}
	return fee, nil
}

func (s serviceClient) GenServiceResponseMsgs(events original.Events, serviceName string,
	provider original.AccAddress,
	handler rpc.ServiceRespondCallback) (msgs []original.Msg) {

	var ids []string
	for _, e := range events.Filter(eventTypeNewBatchRequestProvider) {
		svcName := e.Attributes.GetValue(attributeKeyServiceName)
		prov := e.Attributes.GetValue(attributeKeyProvider)
		if svcName == serviceName && prov == provider.String() {
			reqIDsStr := e.Attributes.GetValue(attributeKeyRequests)
			var idsTemp []string
			if err := json.Unmarshal([]byte(reqIDsStr), &idsTemp); err != nil {
				s.Logger.Err(err).
					Str(attributeKeyRequestID, reqIDsStr).
					Str(attributeKeyServiceName, serviceName).
					Str(attributeKeyProvider, provider.String()).
					Msg("service request don't exist")
				return
			}
			ids = append(ids, idsTemp...)
		}
	}

	for _, reqID := range ids {
		request, err := s.QueryRequest(reqID)
		if err != nil {
			s.Logger.Err(err).
				Str(attributeKeyRequestID, reqID).
				Str(attributeKeyServiceName, serviceName).
				Str(attributeKeyProvider, provider.String()).
				Msg("service request don't exist")
			continue
		}
		//check again
		if provider.Equals(request.Provider) && request.ServiceName == serviceName {
			output, result := handler(request.RequestContextID, reqID, request.Input)
			msgs = append(msgs, MsgRespondService{
				RequestID: hexBytesFrom(reqID),
				Provider:  provider,
				Output:    output,
				Result:    result,
			})
		}
	}
	return msgs
}
