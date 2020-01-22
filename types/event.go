package types

import (
	"bytes"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Subscription interface {
	Unsubscribe()
	GetData() EventData
}

type EventCallback func(Subscription)
type EventData interface{}
type EventDataTx struct {
	Hash   string                 `json:"hash"`
	Height int64                  `json:"height"`
	Index  uint32                 `json:"index"`
	Tx     StdTx                  `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
}
type EventDataNewBlock = tmtypes.EventDataNewBlock

type KVPair map[string]string

func NewKVPair() KVPair {
	return make(map[string]string)
}
func (kv KVPair) Put(k, v string) {
	kv[k] = v
}
func (kv KVPair) ToQueryString() string {
	var buf bytes.Buffer
	for k, v := range kv {
		if buf.Len() > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("%s='%s'", k, v))
	}
	return buf.String()
}
func EventQueryTxFor(txHash string) KVPair {
	kv := NewKVPair()
	kv.Put(tmtypes.TxHashKey, txHash)
	return EventQueryTx(kv)
}

func EventQueryTx(kv KVPair) KVPair {
	kv.Put(tmtypes.EventTypeKey, tmtypes.EventTx)
	return kv
}
