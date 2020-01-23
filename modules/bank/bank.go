package bank

import (
	"github.com/irisnet/irishub-sdk-go/types"
)

func NewBankClient(tm types.TxCtxManager) Bank {
	return bankClient{
		TxCtxManager: tm,
	}
}

// GetAccount return account information specified address
func (b bankClient) GetAccount(address string) (types.BaseAccount, error) {
	return b.QueryAccount(address)
}

// GetTokenStats return token information specified tokenID
func (b bankClient) GetTokenStats(tokenID string) (result types.TokenStats, err error) {
	param := QueryTokenParams{TokenId: tokenID}
	err = b.Query("custom/acc/tokenStats", param, &result)
	if err != nil {
		return result, err
	}
	return
}

func (b bankClient) Send(to string, amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender := b.GetSender(baseTx.From)
	outAddr := types.MustAccAddressFromBech32(to)
	in := []types.Input{types.NewInput(sender, amount)}
	out := []types.Output{types.NewOutput(outAddr, amount)}

	msg := types.NewMsgSend(in, out)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return b.Broadcast(baseTx, []types.Msg{msg})
}
func (b bankClient) Burn(amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender := b.GetSender(baseTx.From)
	msg := types.NewMsgBurn(sender, amount)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return b.Broadcast(baseTx, []types.Msg{msg})
}
func (b bankClient) SetMemoRegexp(memoRegexp string, baseTx types.BaseTx) (types.Result, error) {
	sender := b.GetSender(baseTx.From)
	msg := types.NewMsgSetMemoRegexp(sender, memoRegexp)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return b.Broadcast(baseTx, []types.Msg{msg})
}
