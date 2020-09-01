package types

import (
	abci "github.com/tendermint/tendermint/abci/types"
	//"github.com/tendermint/tendermint/libs/kv"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"encoding/json"
)

type Block struct {
	tmtypes.Header                  `json:"header"`
	Data                            `json:"data"`
	Evidence   tmtypes.EvidenceData `json:"evidence"`
	LastCommit *tmtypes.Commit      `json:"last_commit"`
}

type Data struct {
	Txs []StdTx `json:"txs"`
}

func ParseBlock(cdc Codec, block *tmtypes.Block) Block {
	var txs []StdTx
	for _, tx := range block.Txs {
		var stdTx StdTx
		if err := cdc.UnmarshalBinaryBare(tx, &stdTx); err == nil {
			txs = append(txs, stdTx)
		}
	}
	return Block{
		Header: block.Header,
		Data: Data{
			Txs: txs,
		},
		Evidence:   block.Evidence,
		LastCommit: block.LastCommit,
	}
}

type BlockResult struct {
	Height                int64               `json:"height"`
	TxsResults            []ResponseDeliverTx `json:"txs_results"`
	BeginBlockEvents      []Event             `json:"begin_block_events"`
	EndBlockEvents        []Event             `json:"end_block_events"`
	ValidatorUpdates      []ValidatorUpdate   `json:"validator_updates"`
	ConsensusParamUpdates ConsensusParams     `json:"consensus_param_updates"`
}

type ResponseDeliverTx struct {
	Code      uint32  `json:"code"`
	Data      string  `json:"data"`
	Log       string  `json:"log"`
	Info      string  `json:"info"`
	GasWanted int64   `json:"gas_wanted"`
	GasUsed   int64   `json:"gas_used"`
	Events    []Event `json:"events"`
	Codespace string  `json:"codespace"`
}

type ResultBeginBlock struct {
	Tags Tags `json:"tags"`
}

type ResultEndBlock struct {
	Tags             Tags              `json:"tags"`
	ValidatorUpdates []ValidatorUpdate `json:"validator_updates"`
}

func ParseValidatorUpdate(updates []abci.ValidatorUpdate) []ValidatorUpdate {
	var vUpdates []ValidatorUpdate
	for _, v := range updates {
		data, _ := json.Marshal(v.PubKey.Sum)
		vUpdates = append(vUpdates, ValidatorUpdate{
			PubKey: PubKey{
				Sum: string(data),
				//Type:  v.PubKey.Type,
				//Value: base64.StdEncoding.EncodeToString(v.PubKey.Data),
			},
			Power: v.Power,
		})
	}
	return vUpdates
}

func ParseTxsResults(deliver []*abci.ResponseDeliverTx) []ResponseDeliverTx {
	var rDeliverTxs = make([]ResponseDeliverTx, len(deliver))
	for i, v := range deliver {
		rDeliverTxs[i] = ResponseDeliverTx{
			Code:      v.Code,
			Data:      string(v.Data),
			Log:       v.Log,
			Info:      v.Info,
			GasWanted: v.GasWanted,
			GasUsed:   v.GasUsed,
			Events:    ParseBlockEvent(v.Events),
			Codespace: v.Codespace,
		}
	}
	return rDeliverTxs
}

func ParseBlockEvent(event []abci.Event) []Event {
	var events = make([]Event, len(event))

	for i, v := range event {
		events[i] = Event{
			Type:       v.Type,
			Attributes: ParsePair(v.Attributes),
		}
	}
	return events
}

func ParsePair(kvPairs []abci.EventAttribute) []Pair {
	var pairs = make([]Pair, len(kvPairs))
	for i, v := range kvPairs {
		pairs[i] = Pair{
			Key:   string(v.Key),
			Value: string(v.Value),
		}
	}
	return pairs
}

func ParseConsensusParamUpdates(con *abci.ConsensusParams) ConsensusParams {
	var consensusParams ConsensusParams
	if con != nil {
		consensusParams = ConsensusParams{
			Block: BlockParams{
				MaxBytes: con.Block.MaxBytes,
				MaxGas:   con.Block.MaxGas,
			},
			Evidence: EvidenceParams{
				MaxAgeNumBlocks: con.Evidence.MaxAgeNumBlocks,
				MaxAgeDuration:  con.Evidence.MaxAgeDuration,
			},
			Validator: ValidatorParams{
				PubKeyTypes: con.Validator.PubKeyTypes,
			},
		}
	}
	return consensusParams
}

func ParseBlockResult(res *ctypes.ResultBlockResults) BlockResult {
	return BlockResult{
		Height:                res.Height,
		TxsResults:            ParseTxsResults(res.TxsResults),
		BeginBlockEvents:      ParseBlockEvent(res.BeginBlockEvents),
		EndBlockEvents:        ParseBlockEvent(res.EndBlockEvents),
		ValidatorUpdates:      ParseValidatorUpdate(res.ValidatorUpdates),
		ConsensusParamUpdates: ParseConsensusParamUpdates(res.ConsensusParamUpdates),
	}
}

func ParseValidators(vs []*tmtypes.Validator) []Validator {
	var validators = make([]Validator, len(vs))
	for i, v := range vs {
		var pubKey PubKey
		if bz, err := codec.MarshalJSON(v.PubKey); err == nil {
			_ = codec.UnmarshalJSON(bz, &pubKey)
		}
		validators[i] = Validator{
			Address:          v.Address.String(),
			PubKey:           pubKey,
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		}
	}
	return validators
}
