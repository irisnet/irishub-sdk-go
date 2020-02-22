package bank

import (
	"github.com/irisnet/irishub-sdk-go/types"
)

/**
 * This module is mainly used to transfer coins between accounts,
 * query account balances, and provide common offline transaction signing and broadcasting methods.
 * In addition, the available units of tokens in the IRIShub system are defined using [coin-type](https://www.irisnet.org/docs/concepts/coin-type.html).
 *
 * [More Details](https://www.irisnet.org/docs/features/bank.html)
 *
 * @category Modules
 */
type bankClient struct {
	types.AbstractClient
}

func New(ac types.AbstractClient) types.Bank {
	return bankClient{
		AbstractClient: ac,
	}
}

// QueryAccount return account information specified address
func (b bankClient) QueryAccount(address string) (types.BaseAccount, error) {
	return b.AbstractClient.QueryAccount(address)
}

// GetTokenStats return token statistic, including total loose tokens, total burned tokens and total bonded tokens.
func (b bankClient) QueryTokenStats(tokenID string) (result types.TokenStats, err error) {
	param := QueryTokenParams{TokenId: tokenID}
	err = b.Query("custom/acc/tokenStats", param, &result)
	if err != nil {
		return result, err
	}
	return
}

//Send is responsible for transferring tokens from `From` to `to` account
func (b bankClient) Send(to string, amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender := b.QueryAddress(baseTx.From)
	in := []Input{
		NewInput(sender, amount),
	}

	outAddr := types.MustAccAddressFromBech32(to)
	out := []Output{
		NewOutput(outAddr, amount),
	}

	msg := NewMsgSend(in, out)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return b.Broadcast(baseTx, []types.Msg{msg})
}

//Send is responsible for burning some tokens from `From` account
func (b bankClient) Burn(amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender := b.QueryAddress(baseTx.From)
	msg := NewMsgBurn(sender, amount)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return b.Broadcast(baseTx, []types.Msg{msg})
}

//Send is responsible for setting memo regexp for your own address, so that you can only receive coins from transactions with the corresponding memo.
func (b bankClient) SetMemoRegexp(memoRegexp string, baseTx types.BaseTx) (types.Result, error) {
	sender := b.QueryAddress(baseTx.From)
	msg := NewMsgSetMemoRegexp(sender, memoRegexp)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return b.Broadcast(baseTx, []types.Msg{msg})
}

func (b bankClient) subscribeSendTx() {

}
