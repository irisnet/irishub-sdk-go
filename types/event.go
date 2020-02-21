package types

import (
	"bytes"
	"context"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/net"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Subscription struct {
	ws           *net.RPCClient
	ctx          context.Context
	query        string
	subscriberID string
}

func NewSubscription(ctx context.Context, wsClient *net.RPCClient, query, subscriberID string) Subscription {
	return Subscription{
		ws:           wsClient,
		ctx:          ctx,
		query:        query,
		subscriberID: subscriberID,
	}
}
func (s Subscription) Unsubscribe() {
	_ = s.ws.Unsubscribe(s.ctx, s.subscriberID, s.query)
}

//===============EventDataTx for SubscribeTx=================
type EventDataTx struct {
	Hash   string   `json:"hash"`
	Height int64    `json:"height"`
	Index  uint32   `json:"index"`
	Tx     StdTx    `json:"tx"`
	Result TxResult `json:"result"`
}
type TxResult struct {
	Log       string `json:"log,omitempty"`
	GasWanted int64  `json:"gas_wanted"`
	GasUsed   int64  `json:"gas_used"`
	Tags      []Tag  `json:"tags"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type EventTxCallback func(EventDataTx)

//===============EventDataNewBlock for SubscribeNewBlock=================
type EventDataNewBlock struct {
	Block            Block            `json:"block"`
	ResultBeginBlock ResultBeginBlock `json:"result_begin_block"`
	ResultEndBlock   ResultEndBlock   `json:"result_end_block"`
}

type Block struct {
	tmtypes.Header `json:"header"`
	Data           `json:"data"`
	Evidence       tmtypes.EvidenceData `json:"evidence"`
	LastCommit     *tmtypes.Commit      `json:"last_commit"`
}

type Data struct {
	Txs []StdTx `json:"txs"`
}

type ResultBeginBlock struct {
	Tags []Tag `json:"tags"`
}

type ResultEndBlock struct {
	Tags             []Tag             `json:"tags"`
	ValidatorUpdates []ValidatorUpdate `json:"validator_updates"`
}

type ValidatorUpdate struct {
	PubKey EventPubKey `json:"pub_key"`
	Power  int64       `json:"power"`
}

type EventPubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type EventNewBlockCallback func(EventDataNewBlock)

//===============EventQueryBuilder for build query string=================
type EventQueryBuilder struct {
	params map[EventKey]EventValue
}

func NewEventQueryBuilder() *EventQueryBuilder {
	return &EventQueryBuilder{
		params: make(map[EventKey]EventValue),
	}
}

func (eqb *EventQueryBuilder) AddCondition(eventKey EventKey,
	value EventValue) *EventQueryBuilder {
	eqb.params[eventKey] = value
	return eqb
}

func (eqb *EventQueryBuilder) Build() string {
	var buf bytes.Buffer
	for k, v := range eqb.params {
		if buf.Len() > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("%s='%s'", k, v))
	}
	return buf.String()
}

type EventKey string

const (
	TypeKey        EventKey = "tm.event"
	ActionKey      EventKey = "action"
	SenderKey      EventKey = "sender"
	RecipientKey   EventKey = "recipient"
	TxHashKeyKey   EventKey = "tx.hash"
	TxHeightKeyKey EventKey = "tx.height"
)

type EventValue string

const (
	SendValue          EventValue = "send"
	BurnValue          EventValue = "burn"
	SetMemoRegexpValue EventValue = "set-memo-regexp"
)
