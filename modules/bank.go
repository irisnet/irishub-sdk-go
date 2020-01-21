package modules

import (
	"github.com/irisnet/irishub-sdk-go/types"
)

type bankClient struct {
	types.TxManager
}

func NewBank(tm types.TxManager) types.Bank {
	return bankClient{
		TxManager: tm,
	}
}

// GetAccount return account information specified address
func (bc bankClient) GetAccount(address string) (types.BaseAccount, error) {
	return bc.QueryAccount(address)
}

// GetTokenStats return token information specified tokenID
func (bc bankClient) GetTokenStats(tokenID string) (result types.TokenStats, err error) {
	param := types.QueryTokenParams{TokenId: tokenID}
	err = bc.Query("custom/acc/tokenStats", param, &result)
	if err != nil {
		return result, err
	}
	return
}

func (bc bankClient) Send(to string, amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	keystore := bc.GetTxContext().KeyDAO.Read(baseTx.From)
	inAddr := types.MustAccAddressFromBech32(keystore.GetAddress())
	outAddr := types.MustAccAddressFromBech32(to)
	in := []types.Input{types.NewInput(inAddr, amount)}
	out := []types.Output{types.NewOutput(outAddr, amount)}

	msgSend := types.NewMsgSend(in, out)
	if err := msgSend.ValidateBasic(); err != nil {
		return nil, err
	}
	return bc.Broadcast(baseTx, []types.Msg{msgSend})
}
func (bc bankClient) Burn(amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	//TODO
	return nil, nil
}
func (bc bankClient) SetMemoRegexp(memoRegexp string, baseTx types.BaseTx) (types.Result, error) {
	//TODO
	return nil, nil
}
