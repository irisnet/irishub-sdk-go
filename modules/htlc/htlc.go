package htlc

import (
	"encoding/hex"
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
	"github.com/tendermint/iavl/common"
)

type htlcClient struct {
	sdk.BaseClient
	*log.Logger
}

func Create(ac sdk.BaseClient) rpc.Htlc {
	return htlcClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (h htlcClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (h htlcClient) Name() string {
	return ModuleName
}

func (h htlcClient) QueryHtlc(hashLock string) (rpc.HTLC, sdk.Error) {
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
	if err := h.QueryWithResponse("custom/htlc/hltc", params, &htlc); err != nil {
		return rpc.HTLC{}, nil
	}
	return htlc.Convert().(rpc.HTLC), nil
}
