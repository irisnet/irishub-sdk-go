package htlc

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type paramsClient struct {
	sdk.BaseClient
	*log.Logger
}

func Create(ac sdk.BaseClient) rpc.Params {
	return paramsClient{}
}

func (p paramsClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (p paramsClient) Name() string {
	return ModuleName
}

//func (p paramsClient) QueryHTLC(hashLock string) (interface{}, sdk.Error) {
//	//hLock := []byte(hashLock)
//	//
//	//params := struct {
//	//	HashLock common.HexBytes
//	//}{
//	//	HashLock: hLock,
//	//}
//	//
//	//
//	//p.QueryWithResponse("custom/htlc/hltc",)
//	//return nil, err
//}
