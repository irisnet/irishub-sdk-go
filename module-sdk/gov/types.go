package gov

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
	yaml "gopkg.in/yaml.v2"
)

const (
	ModuleName             = "gov"
	AttributeKeyProposalId = "proposal_id"
)

var (
	_ sdk.Msg = &MsgSubmitProposal{}
	_ sdk.Msg = &MsgDeposit{}
	_ sdk.Msg = &MsgVote{}
)

// NewMsgSubmitProposal creates a new MsgSubmitProposal.
//nolint:interfacer
func NewMsgSubmitProposal(content Content, initialDeposit sdk.Coins, proposer sdk.AccAddress) (*MsgSubmitProposal, error) {
	m := &MsgSubmitProposal{
		InitialDeposit: initialDeposit,
		Proposer:       proposer.String(),
	}
	err := m.SetContent(content)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MsgSubmitProposal) GetInitialDeposit() sdk.Coins { return m.InitialDeposit }

func (m *MsgSubmitProposal) GetProposer() sdk.AccAddress {
	proposer, _ := sdk.AccAddressFromBech32(m.Proposer)
	return proposer
}

func (m *MsgSubmitProposal) SetContent(content Content) error {
	msg, ok := content.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Content = any
	return nil
}

func (m *MsgSubmitProposal) GetContent() Content {
	content, ok := m.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

func (m MsgSubmitProposal) Route() string { return ModuleName }

// Type implements Msg
func (m MsgSubmitProposal) Type() string { return "submit_proposal" }

// ValidateBasic implements Msg
func (m MsgSubmitProposal) ValidateBasic() error {
	if m.Proposer == "" {
		return sdk.Wrapf("missing Proposer")
	}
	if !m.InitialDeposit.IsValid() {
		return sdk.Wrapf("invalidCoins coins, %s", m.InitialDeposit.String())
	}
	if m.InitialDeposit.IsAnyNegative() {
		return sdk.Wrapf("invalidCoins coins, %s", m.InitialDeposit.String())
	}

	content := m.GetContent()
	if content == nil {
		return sdk.Wrapf("missing content")
	}

	if err := content.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg
func (m MsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (m MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	proposer, _ := sdk.AccAddressFromBech32(m.Proposer)
	return []sdk.AccAddress{proposer}
}

// String implements the Stringer interface
func (m MsgSubmitProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgSubmitProposal) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(m.Content, &content)
}

func (msg MsgDeposit) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgDeposit) Type() string { return "deposit" }

// ValidateBasic implements Msg
func (msg MsgDeposit) ValidateBasic() error {
	if msg.Depositor == "" {
		return sdk.Wrapf("missing Proposer")
	}
	if !msg.Amount.IsValid() {
		return sdk.Wrapf("invalidCoins coins, %s", msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdk.Wrapf("invalidCoins coins, %s", msg.Amount.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgDeposit) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	depositor, _ := sdk.AccAddressFromBech32(msg.Depositor)
	return []sdk.AccAddress{depositor}
}

func (msg MsgVote) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgVote) Type() string { return "vote" }

// ValidateBasic implements Msg
func (msg MsgVote) ValidateBasic() error {
	if msg.Voter == "" {
		return sdk.Wrapf("missing Proposer")
	}

	if !ValidVoteOption(msg.Option) {
		return sdk.Wrapf("invalid vote option %s", msg.Option.String())
	}

	return nil
}

func ValidVoteOption(option VoteOption) bool {
	if option == OptionYes ||
		option == OptionAbstain ||
		option == OptionNo ||
		option == OptionNoWithVeto {
		return true
	}
	return false
}

// String implements the Stringer interface
func (msg MsgVote) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	voter, _ := sdk.AccAddressFromBech32(msg.Voter)
	return []sdk.AccAddress{voter}
}

func (q Proposal) Convert() interface{} {
	return QueryProposalResp{
		ProposalId: q.ProposalId,
		Status:     ProposalStatus_name[int32(q.Status)],
		FinalTallyResult: QueryTallyResultResp{
			Yes:        q.FinalTallyResult.Yes,
			Abstain:    q.FinalTallyResult.Abstain,
			No:         q.FinalTallyResult.No,
			NoWithVeto: q.FinalTallyResult.NoWithVeto,
		},
		SubmitTime:      q.SubmitTime,
		DepositEndTime:  q.DepositEndTime,
		TotalDeposit:    q.TotalDeposit,
		VotingStartTime: q.VotingStartTime,
		VotingEndTime:   q.VotingEndTime,
	}
}

type Proposals []Proposal

func (qs Proposals) Convert() interface{} {
	var res []QueryProposalResp
	for _, q := range qs {
		res = append(res, q.Convert().(QueryProposalResp))
	}
	return res
}

func (v Vote) Convert() interface{} {
	return QueryVoteResp{
		ProposalId: v.ProposalId,
		Voter:      v.Voter,
		Option:     int32(v.Option),
	}
}

type Votes []Vote

func (vs Votes) Convert() interface{} {
	var res []QueryVoteResp
	for _, v := range vs {
		res = append(res, v.Convert().(QueryVoteResp))
	}
	return res
}

func (q QueryParamsResponse) Convert() interface{} {
	return QueryParamsResp{
		VotingParams: votingParams{
			VotingPeriod: q.VotingParams.VotingPeriod,
		},
		DepositParams: depositParams{
			MinDeposit:       q.DepositParams.MinDeposit,
			MaxDepositPeriod: q.DepositParams.MaxDepositPeriod,
		},
		TallyParams: tallyParams{
			Quorum:        q.TallyParams.Quorum,
			Threshold:     q.TallyParams.Threshold,
			VetoThreshold: q.TallyParams.VetoThreshold,
		},
	}
}

func (d Deposit) Convert() interface{} {
	return QueryDepositResp{
		ProposalId: d.ProposalId,
		Depositor:  d.Depositor,
		Amount:     d.Amount,
	}
}

type Deposits []Deposit

func (ds Deposits) Convert() interface{} {
	var res []QueryDepositResp
	for _, d := range ds {
		res = append(res, d.Convert().(QueryDepositResp))
	}
	return res
}

func (t TallyResult) Convert() interface{} {
	return QueryTallyResultResp{
		Yes:        t.Yes,
		Abstain:    t.Abstain,
		No:         t.No,
		NoWithVeto: t.NoWithVeto,
	}
}
