package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateLaunchPad = "create_launchpad"
	TypeMsgUpdateLaunchPad = "update_launchpad"
	TypeMsgCloseLaunchPad  = "close_launchpad"
	TypeMsgMintNFT         = "mint_nft"
	TypeMsgMintNFTs        = "mint_nfts"
)

var _ sdk.Msg = &MsgCreateLaunchPad{}

func NewMsgCreateLaunchPad(sender sdk.AccAddress, pad LaunchPad,
) *MsgCreateLaunchPad {
	return &MsgCreateLaunchPad{
		Sender: sender.String(),
		Pad:    pad,
	}
}

func (msg MsgCreateLaunchPad) Route() string { return RouterKey }

func (msg MsgCreateLaunchPad) Type() string { return TypeMsgCreateLaunchPad }

func (msg MsgCreateLaunchPad) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateLaunchPad) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateLaunchPad) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgUpdateLaunchPad{}

func NewMsgUpdateLaunchPad(sender sdk.AccAddress, pad LaunchPad,
) *MsgUpdateLaunchPad {
	return &MsgUpdateLaunchPad{
		Sender: sender.String(),
		Pad:    pad,
	}
}

func (msg MsgUpdateLaunchPad) Route() string { return RouterKey }

func (msg MsgUpdateLaunchPad) Type() string { return TypeMsgUpdateLaunchPad }

func (msg MsgUpdateLaunchPad) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateLaunchPad) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgUpdateLaunchPad) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgCloseLaunchPad{}

func NewMsgCloseLaunchPad(sender sdk.AccAddress, collId uint64,
) *MsgCloseLaunchPad {
	return &MsgCloseLaunchPad{
		Sender: sender.String(),
		CollId: collId,
	}
}

func (msg MsgCloseLaunchPad) Route() string { return RouterKey }

func (msg MsgCloseLaunchPad) Type() string { return TypeMsgCloseLaunchPad }

func (msg MsgCloseLaunchPad) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCloseLaunchPad) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCloseLaunchPad) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgMintNFT{}

func NewMsgMintNFT(sender sdk.AccAddress, collId uint64, name string,
) *MsgMintNFT {
	return &MsgMintNFT{
		Sender:       sender.String(),
		CollectionId: collId,
		Name:         name,
	}
}

func (msg MsgMintNFT) Route() string { return RouterKey }

func (msg MsgMintNFT) Type() string { return TypeMsgMintNFT }

func (msg MsgMintNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintNFT) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgMintNFTs{}

func NewMsgMintNFTs(sender sdk.AccAddress, collId uint64, number uint64,
) *MsgMintNFTs {
	return &MsgMintNFTs{
		Sender:       sender.String(),
		CollectionId: collId,
		Number:       number,
	}
}

func (msg MsgMintNFTs) Route() string { return RouterKey }

func (msg MsgMintNFTs) Type() string { return TypeMsgMintNFTs }

func (msg MsgMintNFTs) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintNFTs) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgMintNFTs) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
