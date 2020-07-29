package service

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/libs/bytes"
)

// queryRequestContextByTxQuery will query for a single request context via a direct txs tags query.
func (s serviceClient) queryRequestContextByTxQuery(reqCtxID string) (requestContext, error) {
	txHash, msgIndex, err := splitRequestContextID(reqCtxID)
	if err != nil {
		return requestContext{}, err
	}

	txInfo, err := s.QueryTx(hex.EncodeToString(txHash))
	if err != nil {
		return requestContext{}, err
	}

	if int64(len(txInfo.Tx.GetMsgs())) > msgIndex {
		msg := txInfo.Tx.GetMsgs()[msgIndex]
		if msg, ok := msg.(MsgCallService); ok {
			return requestContext{
				ServiceName:        msg.ServiceName,
				Providers:          msg.Providers,
				Consumer:           msg.Consumer,
				Input:              msg.Input,
				ServiceFeeCap:      msg.ServiceFeeCap,
				Timeout:            msg.Timeout,
				SuperMode:          msg.SuperMode,
				Repeated:           msg.Repeated,
				RepeatedFrequency:  msg.RepeatedFrequency,
				RepeatedTotal:      msg.RepeatedTotal,
				BatchCounter:       uint64(msg.RepeatedTotal),
				BatchRequestCount:  0,
				BatchResponseCount: 0,
				State:              0,
				ResponseThreshold:  0,
				ModuleName:         "",
			}, nil
		}
	}
	return requestContext{}, errors.New(fmt.Sprintf("invalid reqCtxID:%s", reqCtxID))
}

// queryRequestByTxQuery will query for a single request via a direct txs tags query.
func (s serviceClient) queryRequestByTxQuery(requestID string) (request, error) {
	//reqCtxID, _, requestHeight, batchRequestIndex, err := splitRequestID(requestID)
	//if err != nil {
	//	return request{}, err
	//}
	//
	//// query request context
	//reqCtx, err := s.QueryRequestContext(hex.EncodeToString(reqCtxID))
	//if err != nil {
	//	return request{}, err
	//}
	//
	//blockResult, err := s.BlockResults(&requestHeight)
	//if err != nil {
	//	return request{}, err
	//}

	//for _, tag := range blockResult.Results.EndBlock.Tags {
	//	key := actionTagKey(actionNewBatchRequest, reqCtxID.String())
	//	if string(tag.Key) == string(key) {
	//		var requests []compactRequest
	//		if err := json.Unmarshal(tag.GetValue(), &requests); err != nil {
	//			return request{}, err
	//		}
	//		if len(requests) > int(batchRequestIndex) {
	//			compactRequest := requests[batchRequestIndex]
	//			return request{
	//				ID:                         requestID,
	//				ServiceName:                reqCtx.ServiceName,
	//				Provider:                   compactRequest.Provider,
	//				Consumer:                   reqCtx.Consumer,
	//				Input:                      reqCtx.Input,
	//				ServiceFee:                 compactRequest.ServiceFee,
	//				SuperMode:                  reqCtx.SuperMode,
	//				RequestHeight:              compactRequest.RequestHeight,
	//				ExpirationHeight:           compactRequest.RequestHeight + reqCtx.Timeout,
	//				RequestContextID:           compactRequest.RequestContextID,
	//				RequestContextBatchCounter: compactRequest.RequestContextBatchCounter,
	//			}, nil
	//		}
	//	}
	//}

	return request{}, errors.New(fmt.Sprintf("invalid requestID:%s", requestID))
}

// queryResponseByTxQuery will query for a single request via a direct txs tags query.
func (s serviceClient) queryResponseByTxQuery(requestID string) (response, error) {
	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(sdk.ActionKey).EQ("respond_service")).
		AddCondition(sdk.Cond(tagRequestID).EQ(sdk.EventValue(requestID)))

	result, err := s.QueryTxs(builder, 1, 1)
	if err != nil {
		return response{}, err
	}

	if len(result.Txs) == 0 {
		return response{}, fmt.Errorf("unknown response: %s", requestID)
	}

	reqCtxID, batchCounter, _, _, err := splitRequestID(requestID)
	if err != nil {
		return response{}, err
	}

	// query request context
	reqCtx, err := s.QueryRequestContext(hex.EncodeToString(reqCtxID))
	if err != nil {
		return response{}, err
	}

	for _, msg := range result.Txs[0].Tx.GetMsgs() {
		if responseMsg, ok := msg.(MsgRespondService); ok {
			//if responseMsg.RequestID.String() != requestID {
			//	continue
			//}
			return response{
				Provider:                   responseMsg.Provider,
				Consumer:                   reqCtx.Consumer,
				Output:                     responseMsg.Output,
				Result:                     responseMsg.Result,
				RequestContextID:           reqCtxID,
				RequestContextBatchCounter: batchCounter,
			}, nil
		}
	}

	return response{}, nil
}

// SplitRequestContextID splits the given contextID to txHash and msgIndex
func splitRequestContextID(reqCtxID string) (bytes.HexBytes, int64, error) {
	contextID, err := hex.DecodeString(reqCtxID)
	if err != nil {
		return nil, 0, errors.New("invalid request context id")
	}

	if len(contextID) != contextIDLen {
		return nil, 0, errors.New(fmt.Sprintf("invalid request context id:%s", reqCtxID))
	}

	txHash := contextID[0:32]
	msgIndex := int64(binary.BigEndian.Uint64(contextID[32:40]))
	return txHash, msgIndex, nil
}

// SplitRequestID splits the given contextID to contextID, batchCounter, requestHeight, batchRequestIndex
func splitRequestID(reqID string) (bytes.HexBytes, uint64, int64, int16, error) {
	requestID, err := hex.DecodeString(reqID)
	if err != nil {
		return nil, 0, 0, 0, errors.New("invalid request id")
	}

	if len(requestID) != requestIDLen {
		return nil, 0, 0, 0, errors.New("invalid request id")
	}

	reqCtxID := requestID[0:40]
	batchCounter := binary.BigEndian.Uint64(requestID[40:48])
	requestHeight := int64(binary.BigEndian.Uint64(requestID[48:56]))
	batchRequestIndex := int16(binary.BigEndian.Uint16(requestID[56:]))
	return reqCtxID, batchCounter, requestHeight, batchRequestIndex, nil
}
