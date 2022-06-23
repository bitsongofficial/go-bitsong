package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAuction       = "create_auction"
	TypeMsgSetAuctionAuthority = "set_auction_authority"
	TypeMsgStartAuction        = "start_auction"
	TypeMsgEndAuction          = "end_auction"
	TypeMsgPlaceBid            = "place_bid"
	TypeMsgCancelBid           = "cancel_bid"
	TypeMsgClaimBid            = "claim_bid"
)

var _ sdk.Msg = &MsgCreateAuction{}

func NewMsgCreateAuction(sender sdk.AccAddress, nftId string, prizeType AuctionPrizeType,
	bidDenom string, duration time.Duration,
	priceFloor, instantSalePrice, tickSize, editionLimit uint64,
) *MsgCreateAuction {
	return &MsgCreateAuction{
		Sender:           sender.String(),
		NftId:            nftId,
		PrizeType:        prizeType,
		BidDenom:         bidDenom,
		Duration:         duration,
		PriceFloor:       priceFloor,
		InstantSalePrice: instantSalePrice,
		TickSize:         tickSize,
		EditionLimit:     editionLimit,
	}
}

func (msg MsgCreateAuction) Route() string { return RouterKey }

func (msg MsgCreateAuction) Type() string { return TypeMsgCreateAuction }

func (msg MsgCreateAuction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.NftId == "" {
		return ErrInvalidNftId
	}

	if msg.InstantSalePrice == 0 {
		return ErrInvalidInstantSalePrice
	}

	if msg.Duration == 0 {
		return ErrInvalidDuration
	}

	if msg.BidDenom == "" {
		return ErrInvalidBidDenom
	}

	if AuctionPrizeType_name[int32(msg.PrizeType)] == "" {
		return ErrInvalidPrizeType
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAuction) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateAuction) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgSetAuctionAuthority{}

func NewMsgSetAuctionAuthority(sender sdk.AccAddress, auctionId uint64, newAuthority string) *MsgSetAuctionAuthority {
	return &MsgSetAuctionAuthority{
		Sender:       sender.String(),
		AuctionId:    auctionId,
		NewAuthority: newAuthority,
	}
}

func (msg MsgSetAuctionAuthority) Route() string { return RouterKey }

func (msg MsgSetAuctionAuthority) Type() string { return TypeMsgSetAuctionAuthority }

func (msg MsgSetAuctionAuthority) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.AuctionId == 0 {
		return ErrInvalidAuctionId
	}
	_, err = sdk.AccAddressFromBech32(msg.NewAuthority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auction authority address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetAuctionAuthority) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgSetAuctionAuthority) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgStartAuction{}

func NewMsgStartAuction(sender sdk.AccAddress, auctionId uint64) *MsgStartAuction {
	return &MsgStartAuction{
		Sender:    sender.String(),
		AuctionId: auctionId,
	}
}

func (msg MsgStartAuction) Route() string { return RouterKey }

func (msg MsgStartAuction) Type() string { return TypeMsgStartAuction }

func (msg MsgStartAuction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.AuctionId == 0 {
		return ErrInvalidAuctionId
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgStartAuction) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgStartAuction) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgEndAuction{}

func NewMsgEndAuction(sender sdk.AccAddress, auctionId uint64) *MsgEndAuction {
	return &MsgEndAuction{
		Sender:    sender.String(),
		AuctionId: auctionId,
	}
}

func (msg MsgEndAuction) Route() string { return RouterKey }

func (msg MsgEndAuction) Type() string { return TypeMsgEndAuction }

func (msg MsgEndAuction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.AuctionId == 0 {
		return ErrInvalidAuctionId
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEndAuction) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgEndAuction) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgPlaceBid{}

func NewMsgPlaceBid(sender sdk.AccAddress, auctionId uint64, amount sdk.Coin) *MsgPlaceBid {
	return &MsgPlaceBid{
		Sender:    sender.String(),
		AuctionId: auctionId,
		Amount:    amount,
	}
}

func (msg MsgPlaceBid) Route() string { return RouterKey }

func (msg MsgPlaceBid) Type() string { return TypeMsgPlaceBid }

func (msg MsgPlaceBid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.AuctionId == 0 {
		return ErrInvalidAuctionId
	}

	if msg.Amount.Denom == "" {
		return ErrInvalidBidDenom
	}

	if msg.Amount.Amount.IsZero() {
		return ErrInvalidBidAmount
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgPlaceBid) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgPlaceBid) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgCancelBid{}

func NewMsgCancelBid(sender sdk.AccAddress, auctionId uint64) *MsgCancelBid {
	return &MsgCancelBid{
		Sender:    sender.String(),
		AuctionId: auctionId,
	}
}

func (msg MsgCancelBid) Route() string { return RouterKey }

func (msg MsgCancelBid) Type() string { return TypeMsgCancelBid }

func (msg MsgCancelBid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.AuctionId == 0 {
		return ErrInvalidAuctionId
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCancelBid) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCancelBid) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgClaimBid{}

func NewMsgClaimBid(sender sdk.AccAddress, auctionId uint64) *MsgClaimBid {
	return &MsgClaimBid{
		Sender:    sender.String(),
		AuctionId: auctionId,
	}
}

func (msg MsgClaimBid) Route() string { return RouterKey }

func (msg MsgClaimBid) Type() string { return TypeMsgClaimBid }

func (msg MsgClaimBid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.AuctionId == 0 {
		return ErrInvalidAuctionId
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClaimBid) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgClaimBid) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
