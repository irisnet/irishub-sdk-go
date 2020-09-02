package types

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"time"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

//=========================================Node Status==============================================================================================

// Node Status
type ResultStatus struct {
	NodeInfo      DefaultNodeInfo `json:"node_info"`
	SyncInfo      SyncInfo        `json:"sync_info"`
	ValidatorInfo ValidatorInfo   `json:"validator_info"`
}

type DefaultNodeInfo struct {
	ProtocolVersion ProtocolVersion      `json:"protocol_version"`
	DefaultNodeID   string               `json:"id"`          // authenticated identifier
	ListenAddr      string               `json:"listen_addr"` // accepting incoming
	Network         string               `json:"network"`     // network/chain ID
	Version         string               `json:"version"`     // major.minor.revision
	Channels        string               `json:"channels"`    // channels this node knows about
	Moniker         string               `json:"moniker"`     // arbitrary moniker
	Other           DefaultNodeInfoOther `json:"other"`       // other application specific data
}

// Info about the node's syncing state
type SyncInfo struct {
	LatestBlockHash   string    `json:"latest_block_hash"`
	LatestAppHash     string    `json:"latest_app_hash"`
	LatestBlockHeight int64     `json:"latest_block_height"`
	LatestBlockTime   time.Time `json:"latest_block_time"`

	EarliestBlockHash   string    `json:"earliest_block_hash"`
	EarliestAppHash     string    `json:"earliest_app_hash"`
	EarliestBlockHeight int64     `json:"earliest_block_height"`
	EarliestBlockTime   time.Time `json:"earliest_block_time"`

	CatchingUp bool `json:"catching_up"`
}

// Info about the node's validator
type ValidatorInfo struct {
	Address     string `json:"address"`
	PubKey      ed25519.PubKey `json:"pub_key"`
	VotingPower int64  `json:"voting_power"`
}

type ProtocolVersion struct {
	P2P   uint64 `json:"p2p"`
	Block uint64 `json:"block"`
	App   uint64 `json:"app"`
}

type DefaultNodeInfoOther struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

func ParseNodeStatus(rs *ctypes.ResultStatus) ResultStatus {
	nodeInfo := DefaultNodeInfo{
		ProtocolVersion: ProtocolVersion{
			P2P:   uint64(rs.NodeInfo.ProtocolVersion.P2P),
			Block: uint64(rs.NodeInfo.ProtocolVersion.Block),
			App:   uint64(rs.NodeInfo.ProtocolVersion.App),
		},
		DefaultNodeID: string(rs.NodeInfo.DefaultNodeID),
		ListenAddr:    rs.NodeInfo.ListenAddr,
		Network:       rs.NodeInfo.Network,
		Version:       rs.NodeInfo.Version,
		Channels:      rs.NodeInfo.Channels.String(),
		Moniker:       rs.NodeInfo.Moniker,
		Other: DefaultNodeInfoOther{
			TxIndex:    rs.NodeInfo.Other.TxIndex,
			RPCAddress: rs.NodeInfo.Other.RPCAddress,
		},
	}
	syncInfo := SyncInfo{
		LatestBlockHash:   rs.SyncInfo.LatestBlockHash.String(),
		LatestAppHash:     rs.SyncInfo.LatestAppHash.String(),
		LatestBlockHeight: rs.SyncInfo.LatestBlockHeight,
		LatestBlockTime:   rs.SyncInfo.LatestBlockTime,

		EarliestBlockHash:   rs.SyncInfo.EarliestBlockHash.String(),
		EarliestAppHash:     rs.SyncInfo.EarliestBlockHash.String(),
		EarliestBlockHeight: rs.SyncInfo.EarliestBlockHeight,
		EarliestBlockTime:   rs.SyncInfo.EarliestBlockTime,

		CatchingUp: rs.SyncInfo.CatchingUp,
	}

	//var pubKey PubKey
	//if bz, err := codec.MarshalJSON(rs.ValidatorInfo.PubKey); err == nil {
	//	_ = codec.UnmarshalJSON(bz, &pubKey)
	//}
	validatorInfo := ValidatorInfo{
		Address:     rs.ValidatorInfo.Address.String(),
		PubKey:      ed25519.PubKey(rs.ValidatorInfo.PubKey.Bytes()),
		VotingPower: rs.ValidatorInfo.VotingPower,
	}
	return ResultStatus{
		NodeInfo:      nodeInfo,
		SyncInfo:      syncInfo,
		ValidatorInfo: validatorInfo,
	}
}

//=========================================Genesis==============================================================================================

// GenesisDoc defines the initial conditions for a tendermint blockchain, in particular its validator set.
type GenesisDoc struct {
	GenesisTime     time.Time          `json:"genesis_time"`
	ChainID         string             `json:"chain_id"`
	ConsensusParams *ConsensusParams   `json:"consensus_params,omitempty"`
	Validators      []GenesisValidator `json:"validators,omitempty"`
	AppHash         string             `json:"app_hash"`
	AppState        string             `json:"app_state,omitempty"`
}

type ConsensusParams struct {
	Block     BlockParams     `json:"block"`
	Evidence  EvidenceParams  `json:"evidence"`
	Validator ValidatorParams `json:"validator"`
}

// ValidatorParams restrict the public key types validators can use.
// NOTE: uses ABCI pubkey naming, not Amino names.
type ValidatorParams struct {
	PubKeyTypes []string `json:"pub_key_types"`
}

// HashedParams is a subset of ConsensusParams.
// It is amino encoded and hashed into
// the Header.ConsensusHash.
type HashedParams struct {
	BlockMaxBytes int64
	BlockMaxGas   int64
}

// BlockParams define limits on the block size and gas plus minimum time
// between blocks.
type BlockParams struct {
	MaxBytes int64 `json:"max_bytes"`
	MaxGas   int64 `json:"max_gas"`
	// Minimum time increment between consecutive blocks (in milliseconds)
	// Not exposed to the application.
	TimeIotaMs int64 `json:"time_iota_ms"`
}

// EvidenceParams determine how we handle evidence of malfeasance.
type EvidenceParams struct {
	MaxAgeNumBlocks int64         `json:"max_age_num_blocks"` // only accept new evidence more recent than this
	MaxAgeDuration  time.Duration `json:"max_age_duration"`
}

// GenesisValidator is an initial validator.
type GenesisValidator struct {
	Address string `json:"address"`
	PubKey  ed25519.PubKey `json:"pub_key"`
	Power   int64  `json:"power"`
	Name    string `json:"name"`
}

func ParseGenesis(g *types.GenesisDoc) GenesisDoc {
	consensusParams := ConsensusParams{
		Block: BlockParams{
			MaxBytes:   g.ConsensusParams.Block.MaxBytes,
			MaxGas:     g.ConsensusParams.Block.MaxGas,
			TimeIotaMs: g.ConsensusParams.Block.TimeIotaMs,
		},
		Evidence: EvidenceParams{
			MaxAgeNumBlocks: g.ConsensusParams.Evidence.MaxAgeNumBlocks,
			MaxAgeDuration:  g.ConsensusParams.Evidence.MaxAgeDuration,
		},
		Validator: ValidatorParams{
			PubKeyTypes: g.ConsensusParams.Validator.PubKeyTypes,
		},
	}

	validators := make([]GenesisValidator, 0)
	for _, v := range g.Validators {
		//var pubKey PubKey
		//if bz, err := codec.MarshalJSON(v.PubKey); err == nil {
		//	_ = codec.UnmarshalJSON(bz, &pubKey)
		//}
		validators = append(validators, GenesisValidator{
			Address: v.Address.String(),
			PubKey:  ed25519.PubKey(v.PubKey.Bytes()),
			Power:   v.Power,
			Name:    v.Name,
		})
	}
	res := GenesisDoc{
		GenesisTime:     g.GenesisTime,
		ChainID:         g.ChainID,
		ConsensusParams: &consensusParams,
		Validators:      validators,
		AppHash:         g.AppHash.String(),
		AppState:        string(g.AppState),
	}
	return res
}
