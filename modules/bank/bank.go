// Package bank is mainly used to transfer coins between accounts,
//query account balances, and provide common offline transaction signing and broadcasting methods.
//
// In addition, the available units of tokens in the IRIShub system are defined using [coin-type](https://www.irisnet.org/docs/concepts/coin-type.html).
//
// [More Details](https://www.irisnet.org/docs/features/bank.html)
package bank

import (
	"fmt"
	"strings"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/log"

	"github.com/pkg/errors"

	"github.com/irisnet/irishub-sdk-go/types"
)

type bankClient struct {
	types.AbstractClient
	*log.Logger
}

func Create(ac types.AbstractClient) rpc.Bank {
	return bankClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

func (b bankClient) RegisterCodec(cdc types.Codec) {
	registerCodec(cdc)
}

func (b bankClient) Name() string {
	return ModuleName
}

// QueryAccount return account information specified address
func (b bankClient) QueryAccount(address string) (types.BaseAccount, error) {
	return b.AbstractClient.QueryAccount(address)
}

// GetTokenStats return token statistic, including total loose tokens, total burned tokens and total bonded tokens.
func (b bankClient) QueryTokenStats(tokenID string) (rpc.TokenStats, error) {
	param := struct {
		TokenId string
	}{
		TokenId: tokenID,
	}

	var ts tokenStats
	if err := b.QueryWithResponse("custom/acc/tokenStats", param, &ts); err != nil {
		return rpc.TokenStats{}, err
	}
	return ts.Convert().(rpc.TokenStats), nil
}

//Send is responsible for transferring tokens from `From` to `to` account
func (b bankClient) Send(to string, amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s not found", baseTx.From))
	}
	in := []Input{
		NewInput(sender, amount),
	}

	outAddr, err := types.AccAddressFromBech32(to)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s invalid address", to))
	}
	out := []Output{
		NewOutput(outAddr, amount),
	}

	msg := NewMsgSend(in, out)
	return b.BuildAndSend([]types.Msg{msg}, baseTx)
}

//Send is responsible for burning some tokens from `From` account
func (b bankClient) Burn(amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s not found", baseTx.From))
	}
	msg := NewMsgBurn(sender, amount)
	return b.BuildAndSend([]types.Msg{msg}, baseTx)
}

//Send is responsible for setting memo regexp for your own address, so that you can only receive coins from transactions with the corresponding memo.
func (b bankClient) SetMemoRegexp(memoRegexp string, baseTx types.BaseTx) (types.Result, error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s not found", baseTx.From))
	}
	msg := NewMsgSetMemoRegexp(sender, memoRegexp)
	return b.BuildAndSend([]types.Msg{msg}, baseTx)
}

//SubscribeSendTx Subscribe MsgSend event and return subscription
func (b bankClient) SubscribeSendTx(from, to string, callback rpc.EventMsgSendCallback) types.Subscription {
	var builder = types.NewEventQueryBuilder()

	from = strings.TrimSpace(from)
	if len(from) != 0 {
		builder.AddCondition(types.SenderKey, types.EventValue(from))
	}

	to = strings.TrimSpace(to)
	if len(to) != 0 {
		builder.AddCondition(types.RecipientKey, types.EventValue(to))
	}

	subscription, _ := b.SubscribeTx(builder, func(data types.EventDataTx) {
		for _, msg := range data.Tx.Msgs {
			if value, ok := msg.(MsgSend); ok {
				for i, m := range value.Inputs {
					callback(rpc.EventDataMsgSend{
						Height: data.Height,
						Hash:   data.Hash,
						From:   m.Address.String(),
						To:     value.Outputs[i].Address.String(),
						Amount: m.Coins,
					})
				}
			}
		}
	})
	return subscription
}
