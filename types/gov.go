package types

type Gov interface {
	Module
	Deposit(proposalID uint64, amount Coins, baseTx BaseTx) (Result, error)
	Vote(proposalID uint64, option VoteOption, baseTx BaseTx) (Result, error)

	//TODO
	//QueryProposal(proposalID uint64)
}

type VoteOption string

const (
	Yes        VoteOption = "Yes"
	No         VoteOption = "No"
	NoWithVeto VoteOption = "NoWithVeto"
	Abstain    VoteOption = "Abstain"
)
