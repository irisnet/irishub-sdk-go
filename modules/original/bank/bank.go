package bank

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"strings"

	"github.com/irisnet/irishub-sdk-go/rpc"
	utils "github.com/irisnet/irishub-sdk-go/utils"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type bankClient struct {
	original.BaseClient
	*log.Logger
}

func Create(ac original.BaseClient) rpc.Bank {
	return bankClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (b bankClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (b bankClient) Name() string {
	return ModuleName
}

func (b bankClient) QueryBalances(address, denom string) (original.Balances, original.Error) {
	param := struct {
		Denom   string `json:"denom"`
		Address string `json:"address"`
	}{
		Address: address,
		Denom:   denom,
	}
	var ts original.Balances
	uri := fmt.Sprintf("custom/%s/balance", b.Name())
	if denom == "" {
		uri = fmt.Sprintf("custom/%s/all_balances", b.Name())
	}

	if err := b.QueryWithResponse(uri, param, &ts); err != nil {
		return original.Balances{}, original.Wrap(err)
	}
	return ts, nil
}

// QueryAccount return account information specified address
func (b bankClient) QueryAccount(address string) (original.BaseAccount, original.Error) {
	account, err := b.BaseClient.QueryAccount(address)
	if err != nil {
		return original.BaseAccount{}, original.Wrap(err)
	}
	return account, nil
}

// GetTokenStats return token statistic, including total loose tokens, total burned tokens and total bonded tokens.
func (b bankClient) QueryTokenStats(tokenID string) (rpc.TokenStats, original.Error) {
	param := struct {
		TokenId string
	}{
		TokenId: tokenID,
	}

	var ts tokenStats
	if err := b.QueryWithResponse("custom/acc/tokenStats", param, &ts); err != nil {
		return rpc.TokenStats{}, original.Wrap(err)
	}
	return ts.Convert().(rpc.TokenStats), nil
}

// Query the total supply of coins of the chain
func (b bankClient) QueryTotalSupply() (original.Coins, original.Error) {
	param := struct {
		Page int
	}{
		Page: 1,
	}

	var total original.Coins
	res, err := b.Query("custom/bank/total_supply", param)
	if err != nil {
		return original.Coins{}, nil
	}

	if err := cdc.UnmarshalJSON(res, &total); err != nil {
		return original.Coins{}, nil
	}
	return total, nil
}

//Send is responsible for transferring tokens from `From` to `to` account
func (b bankClient) Send(to string, amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	sender, err := b.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrapf("%s not found", baseTx.From)
	}

	amt, err := b.ToMinCoin(amount...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	in := []Input{
		NewInput(sender, amt),
	}

	outAddr, err := original.AccAddressFromBech32(to)
	if err != nil {
		return original.ResultTx{}, original.Wrapf(fmt.Sprintf("%s invalid address", to))
	}
	out := []Output{
		NewOutput(outAddr, amt),
	}

	msg := NewMsgSend(in, out)
	return b.BuildAndSend([]original.Msg{msg}, baseTx)
}

func (b bankClient) MultiSend(receipts rpc.Receipts, baseTx original.BaseTx) (resTxs []original.ResultTx, err original.Error) {
	sender, err := b.QueryAddress(baseTx.From)
	if err != nil {
		return nil, original.Wrapf("%s not found", baseTx.From)
	}

	if len(receipts) > maxMsgLen {
		return b.SendBatch(sender, receipts, baseTx)
	}

	var inputs = make([]Input, len(receipts))
	var outputs = make([]Output, len(receipts))
	//for i, receipt := range receipts {
	//	amt, err := b.ToMinCoin(receipt.Amount...)
	//	if err != nil {
	//		return nil, sdk.Wrap(err)
	//	}
	//
	//	outAddr, e := sdk.AccAddressFromBech32(receipt.Address)
	//	if e != nil {
	//		return nil, sdk.Wrapf(fmt.Sprintf("%s invalid address", receipt.Address))
	//	}
	//
	//	inputs[i] = NewInput(sender, amt)
	//	outputs[i] = NewOutput(outAddr, amt)
	//}

	msg := NewMsgSend(inputs, outputs)
	res, err := b.BuildAndSend([]original.Msg{msg}, baseTx)
	if err != nil {
		return nil, original.Wrap(err)
	}

	resTxs = append(resTxs, res)
	return
}

func (b bankClient) SendBatch(sender original.AccAddress,
	receipts rpc.Receipts, baseTx original.BaseTx) ([]original.ResultTx, original.Error) {
	batchReceipts := utils.SplitArray(maxMsgLen, receipts)

	var msgs original.Msgs
	for _, receipts := range batchReceipts {
		rs := receipts.(rpc.Receipts)
		var inputs = make([]Input, len(rs))
		var outputs = make([]Output, len(rs))
		//for i, receipt := range rs {
		//	amt, err := b.ToMinCoin(receipt.Amount...)
		//	if err != nil {
		//		return nil, sdk.Wrap(err)
		//	}
		//
		//	outAddr, e := sdk.AccAddressFromBech32(receipt.Address)
		//	if e != nil {
		//		return nil, sdk.Wrapf(fmt.Sprintf("%s invalid address", receipt.Address))
		//	}
		//
		//	inputs[i] = NewInput(sender, amt)
		//	outputs[i] = NewOutput(outAddr, amt)
		//}
		msgs = append(msgs, NewMsgSend(inputs, outputs))
	}
	return b.SendMsgBatch(msgs, baseTx)
}

//Send is responsible for burning some tokens from `From` account
func (b bankClient) Burn(amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	//sender, err := b.QueryAddress(baseTx.From)
	//if err != nil {
	//	return sdk.ResultTx{}, sdk.Wrapf("%s not found", baseTx.From)
	//}
	//
	////amt, err := b.ToMinCoin(amount...)
	//if err != nil {
	//	return sdk.ResultTx{}, sdk.Wrap(err)
	//}
	//msg := NewMsgBurn(sender, amt)
	//return b.BuildAndSend([]sdk.Msg{msg}, baseTx)
	return original.ResultTx{}, nil
}

//Send is responsible for setting memo regexp for your own address, so that you can only receive coins from transactions with the corresponding memo.
func (b bankClient) SetMemoRegexp(memoRegexp string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	sender, err := b.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrapf("%s not found", baseTx.From)
	}
	msg := NewMsgSetMemoRegexp(sender, memoRegexp)
	return b.BuildAndSend([]original.Msg{msg}, baseTx)
}

//SubscribeSendTx Subscribe MsgSend event and return subscription
func (b bankClient) SubscribeSendTx(from, to string, callback rpc.EventMsgSendCallback) original.Subscription {
	var builder = original.NewEventQueryBuilder()

	from = strings.TrimSpace(from)
	if len(from) != 0 {
		builder.AddCondition(original.Cond(original.SenderKey).Contains(original.EventValue(from)))
	}

	to = strings.TrimSpace(to)
	if len(to) != 0 {
		builder.AddCondition(original.Cond(original.RecipientKey).Contains(original.EventValue(to)))
	}

	subscription, _ := b.SubscribeTx(builder, func(data original.EventDataTx) {
		for _, msg := range data.Tx.Msgs {
			if value, ok := msg.(MsgMultiSend); ok {
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
