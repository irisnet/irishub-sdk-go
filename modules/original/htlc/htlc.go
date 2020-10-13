package htlc

import (
	"encoding/hex"
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/irisnet/irishub-sdk-go/utils/log"
	"github.com/tendermint/iavl/common"
)

type htlcClient struct {
	original.BaseClient
	*log.Logger
}

func Create(ac original.BaseClient) rpc.Htlc {
	return htlcClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (h htlcClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (h htlcClient) Name() string {
	return ModuleName
}

func (h htlcClient) QueryHtlc(hashLock string) (rpc.HTLC, original.Error) {
	hash, err := hex.DecodeString(hashLock)
	if err != nil {
		return rpc.HTLC{}, nil
	}

	params := struct {
		HashLock common.HexBytes
	}{
		HashLock: hash,
	}

	var htlc htlc
	if err := h.QueryWithResponse("custom/htlc/htlc", params, &htlc); err != nil {
		return rpc.HTLC{}, nil
	}
	return htlc.Convert().(rpc.HTLC), nil
}
