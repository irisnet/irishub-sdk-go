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
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type bankClient struct {
	sdk.AbstractClient
	*log.Logger
}

func Create(ac sdk.AbstractClient) rpc.Bank {
	return bankClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

func (b bankClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (b bankClient) Name() string {
	return ModuleName
}

// QueryAccount return account information specified address
func (b bankClient) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := b.AbstractClient.QueryAccount(address)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}
	return account, sdk.Nil
}

// GetTokenStats return token statistic, including total loose tokens, total burned tokens and total bonded tokens.
func (b bankClient) QueryTokenStats(tokenID string) (rpc.TokenStats, sdk.Error) {
	param := struct {
		TokenId string
	}{
		TokenId: tokenID,
	}

	var ts tokenStats
	if err := b.QueryWithResponse("custom/acc/tokenStats", param, &ts); err != nil {
		return rpc.TokenStats{}, sdk.Wrap(err)
	}
	return ts.Convert().(rpc.TokenStats), sdk.Nil
}

//Send is responsible for transferring tokens from `From` to `to` account
func (b bankClient) Send(to string, amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, sdk.Error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrapf("%s not found", baseTx.From)
	}
	in := []Input{
		NewInput(sender, amount),
	}

	outAddr, err := sdk.AccAddressFromBech32(to)
	if err != nil {
		return nil, sdk.Wrapf(fmt.Sprintf("%s invalid address", to))
	}
	out := []Output{
		NewOutput(outAddr, amount),
	}

	msg := NewMsgSend(in, out)
	return b.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//Send is responsible for burning some tokens from `From` account
func (b bankClient) Burn(amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, sdk.Error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrapf("%s not found", baseTx.From)
	}
	msg := NewMsgBurn(sender, amount)
	return b.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//Send is responsible for setting memo regexp for your own address, so that you can only receive coins from transactions with the corresponding memo.
func (b bankClient) SetMemoRegexp(memoRegexp string, baseTx sdk.BaseTx) (sdk.Result, sdk.Error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrapf("%s not found", baseTx.From)
	}
	msg := NewMsgSetMemoRegexp(sender, memoRegexp)
	return b.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//SubscribeSendTx Subscribe MsgSend event and return subscription
func (b bankClient) SubscribeSendTx(from, to string, callback rpc.EventMsgSendCallback) sdk.Subscription {
	var builder = sdk.NewEventQueryBuilder()

	from = strings.TrimSpace(from)
	if len(from) != 0 {
		builder.AddCondition(sdk.SenderKey, sdk.EventValue(from))
	}

	to = strings.TrimSpace(to)
	if len(to) != 0 {
		builder.AddCondition(sdk.RecipientKey, sdk.EventValue(to))
	}

	subscription, _ := b.SubscribeTx(builder, func(data sdk.EventDataTx) {
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
