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
func (bc bankClient) GetAccount(address string) (types.BaseAccount, error) {
	return bc.QueryAccount(address)
}

// GetTokenStats return token information specified tokenID
func (bc bankClient) GetTokenStats(tokenID string) (result types.TokenStats, err error) {
	param := QueryTokenParams{TokenId: tokenID}
	err = bc.Query("custom/acc/tokenStats", param, &result)
	if err != nil {
		return result, err
	}
	return
}

func (bc bankClient) Send(to string, amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender := bc.GetSender(baseTx.From)
	outAddr := types.MustAccAddressFromBech32(to)
	in := []types.Input{types.NewInput(sender, amount)}
	out := []types.Output{types.NewOutput(outAddr, amount)}

	msg := types.NewMsgSend(in, out)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return bc.Broadcast(baseTx, []types.Msg{msg})
}
func (bc bankClient) Burn(amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender := bc.GetSender(baseTx.From)
	msg := types.NewMsgBurn(sender, amount)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return bc.Broadcast(baseTx, []types.Msg{msg})
}
func (bc bankClient) SetMemoRegexp(memoRegexp string, baseTx types.BaseTx) (types.Result, error) {
	sender := bc.GetSender(baseTx.From)
	msg := types.NewMsgSetMemoRegexp(sender, memoRegexp)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return bc.Broadcast(baseTx, []types.Msg{msg})
}
