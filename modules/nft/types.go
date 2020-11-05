package nft

import (
	"errors"
	"fmt"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	ModuleName = "nft"

	MinDenomLen = 3
	MaxDenomLen = 64

	MaxTokenURILen = 256
)

var (
	_ sdk.Msg = &MsgIssueDenom{}
	_ sdk.Msg = &MsgTransferNFT{}
	_ sdk.Msg = &MsgEditNFT{}
	_ sdk.Msg = &MsgMintNFT{}
	_ sdk.Msg = &MsgBurnNFT{}
)

var (
	// IsAlphaNumeric only accepts alphanumeric characters
	IsAlphaNumeric   = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	IsBeginWithAlpha = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
)

// Route Implements Msg
func (msg MsgIssueDenom) Route() string { return ModuleName }

// Type Implements Msg
func (msg MsgIssueDenom) Type() string { return "issue_denom" }

// ValidateBasic Implements Msg.
func (msg MsgIssueDenom) ValidateBasic() error {
	if err := ValidateDenomID(msg.Id); err != nil {
		return err
	}

	name := strings.TrimSpace(msg.Name)
	if len(name) > 0 && !utf8.ValidString(name) {
		return errors.New("denom name is invalid")
	}

	if msg.Sender.Empty() {
		return errors.New("missing sender address")
	}
	return nil
}

func (msg MsgIssueDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgIssueDenom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgTransferNFT) Route() string { return ModuleName }

func (msg MsgTransferNFT) Type() string { return "transfer_nft" }

func (msg MsgTransferNFT) ValidateBasic() error {
	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}

	if msg.Sender.Empty() {
		return errors.New("missing sender address")
	}

	if msg.Recipient.Empty() {
		return errors.New("missing recipient address")
	}
	return ValidateTokenID(msg.Id)
}

func (msg MsgTransferNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgEditNFT) Route() string { return ModuleName }

func (msg MsgEditNFT) Type() string { return "edit_nft" }

func (msg MsgEditNFT) ValidateBasic() error {
	if msg.Sender.Empty() {
		return errors.New("missing sender address")
	}

	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}

	if err := ValidateTokenURI(msg.URI); err != nil {
		return err
	}
	return ValidateTokenID(msg.Id)
}

func (msg MsgEditNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgEditNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgMintNFT) Route() string { return ModuleName }

func (msg MsgMintNFT) Type() string { return "mint_nft" }

func (msg MsgMintNFT) ValidateBasic() error {
	if msg.Sender.Empty() {
		return errors.New("missing sender address")
	}
	if msg.Recipient.Empty() {
		return errors.New("missing receipt address")
	}
	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}

	if err := ValidateTokenURI(msg.URI); err != nil {
		return err
	}
	return ValidateTokenID(msg.Id)
}

func (msg MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func NewMsgBurnNFT(sender sdk.AccAddress, id string, denom string) *MsgBurnNFT {
	return &MsgBurnNFT{
		Sender: sender,
		Id:     strings.ToLower(strings.TrimSpace(id)),
		Denom:  strings.TrimSpace(denom),
	}
}

func (msg MsgBurnNFT) Route() string { return ModuleName }

func (msg MsgBurnNFT) Type() string { return "burn_nft" }

func (msg MsgBurnNFT) ValidateBasic() error {
	if msg.Sender.Empty() {
		return errors.New("missing sender address")
	}

	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}
	return ValidateTokenID(msg.Id)
}

func (msg MsgBurnNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func ValidateDenomID(denomID string) error {
	denomID = strings.TrimSpace(denomID)
	if len(denomID) < MinDenomLen || len(denomID) > MaxDenomLen {
		return fmt.Errorf("invalid denom %s, only accepts value [%d, %d]", denomID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(denomID) || !IsAlphaNumeric(denomID) {
		return fmt.Errorf("invalid denom %s, only accepts alphanumeric characters,and begin with an english letter", denomID)
	}
	return nil
}

func ValidateTokenID(tokenID string) error {
	tokenID = strings.TrimSpace(tokenID)
	if len(tokenID) < MinDenomLen || len(tokenID) > MaxDenomLen {
		return fmt.Errorf("invalid tokenID %s, only accepts value [%d, %d]", tokenID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(tokenID) || !IsAlphaNumeric(tokenID) {
		return fmt.Errorf("invalid tokenID %s, only accepts alphanumeric characters,and begin with an english letter", tokenID)
	}
	return nil
}

func ValidateTokenURI(tokenURI string) error {
	if len(tokenURI) > MaxTokenURILen {
		return fmt.Errorf("invalid tokenURI %s, only accepts value [0, %d]", tokenURI, MaxTokenURILen)
	}
	return nil
}

func (o Owner) Convert() interface{} {
	var idcs []IDC
	for _, idc := range o.IDCollections {
		idcs = append(idcs, IDC{
			Denom:    idc.Denom,
			TokenIDs: idc.Ids,
		})
	}
	return QueryOwnerResp{
		Address: o.Address.String(),
		IDCs:    idcs,
	}
}

func (d Denom) Convert() interface{} {
	return QueryDenomResp{
		ID:      d.Id,
		Name:    d.Name,
		Schema:  d.Schema,
		Creator: d.Creator.String(),
	}
}

type denoms []Denom

func (ds denoms) Convert() interface{} {
	var denoms []QueryDenomResp
	for _, denom := range ds {
		denoms = append(denoms, denom.Convert().(QueryDenomResp))
	}
	return denoms
}

func (c Collection) Convert() interface{} {
	var nfts []QueryNFTResp
	for _, nft := range c.NFTs {
		nfts = append(nfts, QueryNFTResp{
			ID:      nft.Id,
			Name:    nft.Name,
			URI:     nft.URI,
			Data:    nft.Data,
			Creator: nft.Owner.String(),
		})
	}
	return QueryCollectionResp{
		Denom: c.Denom.Convert().(QueryDenomResp),
		NFTs:  nfts,
	}
}

func (baseNft BaseNFT) Convert() interface{} {
	return QueryNFTResp{
		ID:      baseNft.Id,
		Name:    baseNft.Name,
		URI:     baseNft.URI,
		Data:    baseNft.Data,
		Creator: baseNft.Owner.String(),
	}
}
