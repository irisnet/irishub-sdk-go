package service

import (
	"encoding/json"
	"strings"

	cmn "github.com/tendermint/tendermint/libs/common"

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

func Create(ac sdk.AbstractClient) rpc.Service {
	return serviceClient{
		AbstractClient: ac,
		Logger:         ac.Logger(),
	}
}

//DefineService is responsible for creating a new service definition
func (s serviceClient) DefineService(request rpc.ServiceDefinitionRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	author, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := MsgDefineService{
		Name:              request.ServiceName,
		Description:       request.Description,
		Tags:              request.Tags,
		Author:            author,
		AuthorDescription: request.AuthorDescription,
		Schemas:           request.Schemas,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//BindService is responsible for binding a new service definition
func (s serviceClient) BindService(request rpc.ServiceBindingRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	amt, err := s.ToMinCoin(request.Deposit...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgBindService{
		ServiceName: request.ServiceName,
		Provider:    provider,
		Deposit:     amt,
		Pricing:     request.Pricing,
		MinRespTime: request.MinRespTime,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//UpdateServiceBinding updates the specified service binding
func (s serviceClient) UpdateServiceBinding(request rpc.UpdateServiceBindingRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	amt, err := s.ToMinCoin(request.Deposit...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgUpdateServiceBinding{
		ServiceName: request.ServiceName,
		Provider:    provider,
		Deposit:     amt,
		Pricing:     request.Pricing,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// DisableServiceBinding disables the specified service binding
func (s serviceClient) DisableServiceBinding(serviceName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := MsgDisableService{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// EnableServiceBinding enables the specified service binding
func (s serviceClient) EnableServiceBinding(serviceName string, deposit sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	amt, err := s.ToMinCoin(deposit...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgEnableService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     amt,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//InvokeService is responsible for invoke a new service and callback `handler`
func (s serviceClient) InvokeService(request rpc.ServiceInvocationRequest, baseTx sdk.BaseTx) (string, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		p, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return "", sdk.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	msg := MsgRequestService{
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
	baseTx.Mode = sdk.Commit

	result, err := s.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	reqCtxID := result.Tags.GetValue(tagRequestContextID)
	if request.Handler == nil {
		return reqCtxID, nil
	}
	_, err = s.RegisterServiceResponseListener(reqCtxID, request.Handler)
	return reqCtxID, sdk.Wrap(err)
}

func (s serviceClient) RegisterServiceResponseListener(reqCtxID string,
	callback rpc.ServiceInvokeHandler) (subscription sdk.Subscription, err sdk.Error) {
	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(sdk.ActionKey).EQ(tagRespondService)).
		AddCondition(sdk.Cond(tagRequestContextID).EQ(sdk.EventValue(reqCtxID)))

	return s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
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
		if err != nil || reqCtx.State == "completed" {
			_ = s.Unsubscribe(subscription)
		}
	})
}

// SetWithdrawAddress sets a new withdrawal address for the specified service binding
func (s serviceClient) SetWithdrawAddress(withdrawAddress string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	withdrawAddr, err := sdk.AccAddressFromBech32(withdrawAddress)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s invalid address", withdrawAddress)
	}
	msg := MsgSetWithdrawAddress{
		Provider:        provider,
		WithdrawAddress: withdrawAddr,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// RefundServiceDeposit refunds the deposit from the specified service binding
func (s serviceClient) RefundServiceDeposit(serviceName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// StartRequestContext starts the specified request context
func (s serviceClient) StartRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := MsgStartRequestContext{
		RequestContextID: hexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// PauseRequestContext suspends the specified request context
func (s serviceClient) PauseRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := MsgPauseRequestContext{
		RequestContextID: hexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// KillRequestContext terminates the specified request context
func (s serviceClient) KillRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := MsgKillRequestContext{
		RequestContextID: hexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// UpdateRequestContext updates the specified request context
func (s serviceClient) UpdateRequestContext(request rpc.UpdateContextRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		p, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		providers = append(providers, p)
	}

	amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgUpdateRequestContext{
		RequestContextID:  hexBytesFrom(request.RequestContextID),
		Providers:         providers,
		ServiceFeeCap:     amt,
		Timeout:           request.Timeout,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		Consumer:          consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// WithdrawEarnedFees withdraws the earned fees to the specified provider
func (s serviceClient) WithdrawEarnedFees(baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgWithdrawEarnedFees{
		Provider: provider,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// WithdrawTax withdraws the service tax to the speicified destination address by the trustee
func (s serviceClient) WithdrawTax(destAddress string, amount sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	trustee, err := s.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	receipt, err := sdk.AccAddressFromBech32(destAddress)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s invalid address", destAddress)
	}

	amt, err := s.ToMinCoin(amount...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgWithdrawTax{
		Trustee:     trustee,
		DestAddress: receipt,
		Amount:      amt,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//RegisterServiceRequestListener is responsible for registering a group of service handler
func (s serviceClient) RegisterServiceRequestListener(serviceRouter rpc.ServiceRouter,
	baseTx sdk.BaseTx) (subscription sdk.Subscription, err sdk.Error) {
	provider, e := s.QueryAddress(baseTx.From)
	if e != nil {
		return sdk.Subscription{}, sdk.Wrap(e)
	}
	var serviceNames []string
	for name, _ := range serviceRouter {
		serviceNames = append(serviceNames, name)
	}
	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(
			actionTagKey(actionNewBatchRequest, tagProvider)).
			Contains(sdk.EventValue(provider.String())))
	return s.SubscribeNewBlock(builder, func(block sdk.EventDataNewBlock) {
		var msgs []sdk.Msg
		for _, serviceName := range serviceNames {
			msgs = append(msgs,
				s.GenServiceResponseMsgs(block.ResultEndBlock.Tags,
					serviceName,
					provider,
					serviceRouter[serviceName])...)

		}
		if msgs == nil || len(msgs) == 0 {
			return
		}
		if _, err = s.SendMsgBatch(5, msgs, baseTx); err != nil {
			s.Err(err).Msg("provider respond failed")
		}
	})
}

//RegisterSingleServiceRequestListener is responsible for registering a single service handler
func (s serviceClient) RegisterSingleServiceRequestListener(serviceName string,
	respondHandler rpc.ServiceRespondHandler,
	baseTx sdk.BaseTx) (subscription sdk.Subscription, err sdk.Error) {
	provider, e := s.QueryAddress(baseTx.From)
	if e != nil {
		return sdk.Subscription{}, sdk.Wrap(e)
	}

	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(
			actionTagKey(actionNewBatchRequest, tagProvider)).
			Contains(sdk.EventValue(provider.String()))).
		AddCondition(sdk.Cond(
			actionTagKey(actionNewBatchRequest, tagServiceName)).
			Contains(sdk.EventValue(serviceName)),
		)
	return s.SubscribeNewBlock(builder, func(block sdk.EventDataNewBlock) {
		msgs := s.GenServiceResponseMsgs(block.ResultEndBlock.Tags, serviceName, provider, respondHandler)
		if _, err = s.SendMsgBatch(5, msgs, baseTx); err != nil {
			s.Err(err).Msg("provider respond failed")
		}
	})
}

// QueryDefinition return a service definition of the specified name
func (s serviceClient) QueryDefinition(serviceName string) (rpc.ServiceDefinition, sdk.Error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}

	var definition serviceDefinition
	if err := s.QueryWithResponse("custom/service/definition", param, &definition); err != nil {
		return rpc.ServiceDefinition{}, sdk.Wrap(err)
	}
	return definition.Convert().(rpc.ServiceDefinition), nil
}

// QueryBinding return the specified service binding
func (s serviceClient) QueryBinding(serviceName string, provider sdk.AccAddress) (rpc.ServiceBinding, sdk.Error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}

	var binding serviceBinding
	if err := s.QueryWithResponse("custom/service/binding", param, &binding); err != nil {
		return rpc.ServiceBinding{}, sdk.Wrap(err)
	}
	return binding.Convert().(rpc.ServiceBinding), nil
}

// QueryBindings returns all bindings of the specified service
func (s serviceClient) QueryBindings(serviceName string) ([]rpc.ServiceBinding, sdk.Error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}

	var bindings serviceBindings
	if err := s.QueryWithResponse("custom/service/bindings", param, &bindings); err != nil {
		return nil, sdk.Wrap(err)
	}
	return bindings.Convert().([]rpc.ServiceBinding), nil
}

// QueryRequest returns  the active request of the specified requestID
func (s serviceClient) QueryRequest(requestID string) (rpc.ServiceRequest, sdk.Error) {
	param := struct {
		RequestID []byte
	}{
		RequestID: hexBytesFrom(requestID),
	}

	var request request
	if err := s.QueryWithResponse("custom/service/request", param, &request); request.Empty() {
		request, err = s.queryRequestByTxQuery(requestID)
		if err != nil {
			return rpc.ServiceRequest{}, sdk.Wrap(err)
		}
	}
	return request.Convert().(rpc.ServiceRequest), nil
}

// QueryRequest returns all the active requests of the specified service binding
func (s serviceClient) QueryRequests(serviceName string, provider sdk.AccAddress) ([]rpc.ServiceRequest, sdk.Error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}

	var rs requests
	if err := s.QueryWithResponse("custom/service/requests", param, &rs); err != nil {
		return nil, sdk.Wrap(err)
	}
	return rs.Convert().([]rpc.ServiceRequest), nil
}

// QueryRequestsByReqCtx returns all requests of the specified request context ID and batch counter
func (s serviceClient) QueryRequestsByReqCtx(reqCtxID string, batchCounter uint64) ([]rpc.ServiceRequest, sdk.Error) {
	param := struct {
		RequestContextID cmn.HexBytes
		BatchCounter     uint64
	}{
		RequestContextID: hexBytesFrom(reqCtxID),
		BatchCounter:     batchCounter,
	}

	var rs requests
	if err := s.QueryWithResponse("custom/service/requests_by_ctx", param, &rs); err != nil {
		return nil, sdk.Wrap(err)
	}
	return rs.Convert().([]rpc.ServiceRequest), nil
}

// QueryResponse returns a response with the speicified request ID
func (s serviceClient) QueryResponse(requestID string) (rpc.ServiceResponse, sdk.Error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}

	var response response
	if err := s.QueryWithResponse("custom/service/response", param, &response); response.Empty() {
		response, err = s.queryResponseByTxQuery(requestID)
		if err != nil {
			return rpc.ServiceResponse{}, sdk.Wrap(nil)
		}
	}
	return response.Convert().(rpc.ServiceResponse), nil
}

// QueryResponses returns all responses of the specified request context and batch counter
func (s serviceClient) QueryResponses(reqCtxID string, batchCounter uint64) ([]rpc.ServiceResponse, sdk.Error) {
	param := struct {
		RequestContextID cmn.HexBytes
		BatchCounter     uint64
	}{
		RequestContextID: hexBytesFrom(reqCtxID),
		BatchCounter:     batchCounter,
	}
	var rs responses
	if err := s.QueryWithResponse("custom/service/responses", param, &rs); err != nil {
		return nil, sdk.Wrap(err)
	}
	return rs.Convert().([]rpc.ServiceResponse), nil
}

// QueryRequestContext return the specified request context
func (s serviceClient) QueryRequestContext(reqCtxID string) (rpc.RequestContext, sdk.Error) {
	param := struct {
		RequestContextID cmn.HexBytes
	}{
		RequestContextID: hexBytesFrom(reqCtxID),
	}

	var reqCtx requestContext
	if err := s.QueryWithResponse("custom/service/context", param, &reqCtx); reqCtx.Empty() {
		reqCtx, err = s.queryRequestContextByTxQuery(reqCtxID)
		if err != nil {
			return rpc.RequestContext{}, sdk.Wrap(err)
		}
	}
	return reqCtx.Convert().(rpc.RequestContext), nil
}

//QueryFees return the earned fees for a provider
func (s serviceClient) QueryFees(provider string) (rpc.EarnedFees, sdk.Error) {
	address, err := sdk.AccAddressFromBech32(provider)
	if err != nil {
		return rpc.EarnedFees{}, sdk.Wrap(err)
	}

	param := struct {
		Address sdk.AccAddress
	}{
		Address: address,
	}

	var fee earnedFees

	if err := s.QueryWithResponse("custom/service/fees", param, &fee); err != nil {
		return rpc.EarnedFees{}, sdk.Wrap(err)
	}
	return fee.Convert().(rpc.EarnedFees), nil
}

func (s serviceClient) GenServiceResponseMsgs(tags sdk.Tags, serviceName string, provider sdk.AccAddress, handler rpc.ServiceRespondHandler) (msgs []sdk.Msg) {
	idsKey := actionTagKey(actionNewBatchRequest, serviceName, provider.String())
	idsStr := tags.GetValue(string(idsKey))
	if len(idsStr) == 0 {
		return
	}

	s.Debug().
		Str(tagServiceName, serviceName).
		Str(tagProvider, provider.String()).
		Str(tagRequestID, idsStr).
		Msg("received service request")

	var ids []string
	if err := json.Unmarshal([]byte(idsStr), &ids); err != nil {
		s.Err(err).
			Str(tagRequestID, idsStr).
			Str(tagServiceName, serviceName).
			Str(tagProvider, provider.String()).
			Msg("service request don't exist")
		return
	}

	for _, reqID := range ids {
		request, err := s.QueryRequest(reqID)
		if err != nil {
			s.Err(err).
				Str(tagRequestID, reqID).
				Str(tagServiceName, serviceName).
				Str(tagProvider, provider.String()).
				Msg("service request don't exist")
			continue
		}
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
