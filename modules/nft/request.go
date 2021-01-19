package nft

import sdk "github.com/irisnet/irishub-sdk-go/types"

type NFTRequestHash interface {
	GetMsgs(sender sdk.AccAddress) ([]sdk.Msg, sdk.Error)
}

func (req EditNFTRequest) GetMsgs(sender sdk.AccAddress) ([]sdk.Msg, sdk.Error) {
	msg := &MsgEditNFT{
		Id:      req.ID,
		Name:    req.Name,
		DenomId: req.Denom,
		URI:     req.URI,
		Data:    req.Data,
		Sender:  sender.String(),
	}
	return []sdk.Msg{msg}, nil
}

func (req MintNFTRequest) GetMsgs(sender sdk.AccAddress) ([]sdk.Msg, sdk.Error) {

	var recipient = sender.String()
	if len(req.Recipient) > 0 {
		if err := sdk.ValidateAccAddress(req.Recipient); err != nil {
			return nil, sdk.Wrap(err)
		}
		recipient = req.Recipient
	}

	msg := &MsgMintNFT{
		Id:        req.ID,
		DenomId:   req.Denom,
		Name:      req.Name,
		URI:       req.URI,
		Data:      req.Data,
		Sender:    sender.String(),
		Recipient: recipient,
	}
	return []sdk.Msg{msg}, nil
}

func (req TransferNFTRequest) GetMsgs(sender sdk.AccAddress) ([]sdk.Msg, sdk.Error) {

	if err := sdk.ValidateAccAddress(req.Recipient); err != nil {
		return nil, sdk.Wrap(err)
	}

	msg := &MsgTransferNFT{
		Id:        req.ID,
		Name:      req.Name,
		DenomId:   req.Denom,
		URI:       req.URI,
		Data:      req.Data,
		Sender:    sender.String(),
		Recipient: req.Recipient,
	}
	return []sdk.Msg{msg}, nil
}

func (req BurnNFTRequest) GetMsgs(sender sdk.AccAddress) ([]sdk.Msg, sdk.Error) {
	msg := &MsgBurnNFT{
		Sender:  sender.String(),
		Id:      req.ID,
		DenomId: req.Denom,
	}
	return []sdk.Msg{msg}, nil
}

func (req IssueDenomRequest) GetMsgs(sender sdk.AccAddress) ([]sdk.Msg, sdk.Error) {
	msg := &MsgIssueDenom{
		Id:     req.ID,
		Name:   req.Name,
		Schema: req.Schema,
		Sender: sender.String(),
	}
	return []sdk.Msg{msg}, nil
}
