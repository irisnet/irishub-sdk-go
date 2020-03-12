package types

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	TypeKey      EventKey = "tm.event"
	ActionKey    EventKey = "action"
	SenderKey    EventKey = "sender"
	RecipientKey EventKey = "recipient"

	TxValue EventValue = "Tx"
)

type EventKey string
type EventValue string
type Subscription struct {
	Ctx   context.Context
	Query string
	ID    string
}

func NewSubscription(ctx context.Context, query, id string) Subscription {
	return Subscription{
		Ctx:   ctx,
		Query: query,
		ID:    id,
	}
}

//===============EventDataTx for SubscribeTx=================
type EventDataTx struct {
	Hash   string   `json:"hash"`
	Height int64    `json:"height"`
	Index  uint32   `json:"index"`
	Tx     StdTx    `json:"tx"`
	Result TxResult `json:"result"`
}

func (tx EventDataTx) MarshalJson() []byte {
	bz, _ := json.Marshal(tx)
	return bz
}

type TxResult struct {
	Log       string `json:"log,omitempty"`
	GasWanted int64  `json:"gas_wanted"`
	GasUsed   int64  `json:"gas_used"`
	Tags      Tags   `json:"tags"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Tags []Tag

func ParseTags(pairs []cmn.KVPair) (tags []Tag) {
	if pairs == nil || len(pairs) == 0 {
		return tags
	}
	for _, pair := range pairs {
		key := string(pair.Key)
		value := string(pair.Value)
		tags = append(tags, Tag{
			Key:   key,
			Value: value,
		})
	}
	return
}
func (t Tags) GetValues(key string) (values []string) {
	for _, tag := range t {
		if tag.Key == key {
			values = append(values, tag.Value)
		}
	}
	return
}

func (t Tags) GetValue(key string) string {
	for _, tag := range t {
		if tag.Key == key {
			return tag.Value
		}
	}
	return ""
}

func (t Tags) String() string {
	var buf bytes.Buffer
	for _, tag := range t {
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%s=%s", tag.Key, tag.Value))
	}
	return buf.String()
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
	Tags Tags `json:"tags"`
}

type ResultEndBlock struct {
	Tags             Tags              `json:"tags"`
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

//===============EventDataNewBlockHeader for SubscribeNewBlockHeader=================
type EventDataNewBlockHeader struct {
	Header tmtypes.Header `json:"header"`

	ResultBeginBlock ResultBeginBlock `json:"result_begin_block"`
	ResultEndBlock   ResultEndBlock   `json:"result_end_block"`
}

type EventNewBlockHeaderCallback func(EventDataNewBlockHeader)

//===============EventDataValidatorSetUpdates for SubscribeValidatorSetUpdates=================
type Validator struct {
	Address          string `json:"address"`
	PubKey           string `json:"pub_key"`
	VotingPower      int64  `json:"voting_power"`
	ProposerPriority int64  `json:"proposer_priority"`
}
type EventDataValidatorSetUpdates struct {
	ValidatorUpdates []Validator `json:"validator_updates"`
}

type EventValidatorSetUpdatesCallback func(EventDataValidatorSetUpdates)

//===============EventQueryBuilder for build query string=================
type condition struct {
	key   EventKey
	value EventValue
	op    string
}

func (c *condition) LTE(v EventValue) *condition {
	c.value = v
	c.op = "<="
	return c
}

func (c *condition) GTE(v EventValue) *condition {
	c.value = v
	c.op = ">="
	return c
}

func (c *condition) LE(v EventValue) *condition {
	c.value = v
	c.op = "<"
	return c
}

func (c *condition) GE(v EventValue) *condition {
	c.value = v
	c.op = ">"
	return c
}

func (c *condition) EQ(v EventValue) *condition {
	c.value = v
	c.op = "="
	return c
}

func (c *condition) Equal(v EventValue) *condition {
	c.value = v
	c.op = "="
	return c
}

func (c *condition) Contains(v EventValue) *condition {
	c.value = v
	c.op = "CONTAINS"
	return c
}

func (c *condition) String() string {
	if len(c.key) == 0 || len(c.value) == 0 || len(c.op) == 0 {
		panic("invalid condition")
	}
	return fmt.Sprintf("%s %s '%s'", c.key, c.op, c.value)
}

func Cond(key EventKey) *condition {
	return &condition{
		key: key,
	}
}

//EventQueryBuilder is responsible for constructing listening conditions
type EventQueryBuilder struct {
	conditions []string
}

func NewEventQueryBuilder() *EventQueryBuilder {
	return &EventQueryBuilder{
		conditions: []string{},
	}
}

//AddCondition is responsible for adding listening conditions
func (eqb *EventQueryBuilder) AddCondition(c *condition) *EventQueryBuilder {
	if c == nil {
		panic("invalid condition")
	}
	eqb.conditions = append(eqb.conditions, c.String())
	return eqb
}

//Build is responsible for constructing the listening condition into a listening instruction identified by tendermint
func (eqb *EventQueryBuilder) Build() string {
	var buf bytes.Buffer
	for _, condition := range eqb.conditions {
		if buf.Len() > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(condition)
	}
	return buf.String()
}
