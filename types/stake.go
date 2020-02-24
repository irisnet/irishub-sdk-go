package types

type Validator struct {
	Address          string `json:"address"`
	PubKey           string `json:"pub_key"`
	VotingPower      int64  `json:"voting_power"`
	ProposerPriority int64  `json:"proposer_priority"`
}
