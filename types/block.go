package types

import (
	"encoding/base64"
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Block struct {
	tmtypes.Header `json:"header"`
	Data           `json:"data"`
	Evidence       tmtypes.EvidenceData `json:"evidence"`
	LastCommit     *tmtypes.Commit      `json:"last_commit"`
}

type Data struct {
	Txs []StdTx `json:"txs"`
}

func ParseBlock(cdc Codec, block *tmtypes.Block) Block {
	var txs []StdTx
	for _, tx := range block.Txs {
		var stdTx StdTx
		if err := cdc.UnmarshalBinaryLengthPrefixed(tx, &stdTx); err == nil {
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
	Height  int64         `json:"height"`
	Results ABCIResponses `json:"results"`
}

type ABCIResponses struct {
	DeliverTx  []TxResult
	EndBlock   ResultEndBlock
	BeginBlock ResultBeginBlock
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
		vUpdates = append(vUpdates, ValidatorUpdate{
			PubKey: PubKey{
				Type:  v.PubKey.Type,
				Value: base64.StdEncoding.EncodeToString(v.PubKey.Data),
			},
			Power: v.Power,
		})
	}
	return vUpdates
}

func ParseBlockResult(res *ctypes.ResultBlockResults) BlockResult {
	var txResults = make([]TxResult, len(res.TxsResults))
	for i, r := range res.TxsResults {
		txResults[i] = TxResult{
			Code:      r.Code,
			Log:       r.Log,
			GasWanted: r.GasWanted,
			GasUsed:   r.GasUsed,
			//Tags:      ParseTags(r.Tags),
		}
	}
	return BlockResult{
		Height: res.Height,
		Results: ABCIResponses{
			DeliverTx: txResults,
			EndBlock: ResultEndBlock{
				//Tags:             ParseTags(res.Results.EndBlock.Tags),
				ValidatorUpdates: ParseValidatorUpdate(res.ValidatorUpdates),
			},
			BeginBlock: ResultBeginBlock{
				//Tags: ParseTags(res.Results.BeginBlock.Tags),
			},
		},
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
