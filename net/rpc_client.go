package net

import (
	"github.com/pkg/errors"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"
)

type RPCClient struct {
	rpc.Client
}

func NewRPCClient(remote string) RPCClient {
	client := rpc.NewHTTP(remote, "/websocket")
	return RPCClient{client}
}

func (r RPCClient) Query(path string, data cmn.HexBytes) (res []byte, err error) {
	result, err := r.ABCIQueryWithOptions(path, data, rpc.DefaultABCIQueryOptions)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.Errorf(resp.Log)
	}

	return resp.Value, nil
}
