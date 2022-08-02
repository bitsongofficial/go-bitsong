package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateCandyMachine = "create_candymachine"
	TypeMsgUpdateCandyMachine = "update_candymachine"
	TypeMsgCloseCandyMachine  = "close_candymachine"
	TypeMsgMintNFT            = "mint_nft"
)

var _ sdk.Msg = &MsgCreateCandyMachine{}

func NewMsgCreateCandyMachine(sender sdk.AccAddress, machine CandyMachine,
) *MsgCreateCandyMachine {
	return &MsgCreateCandyMachine{
		Sender:  sender.String(),
		Machine: machine,
	}
}

func (msg MsgCreateCandyMachine) Route() string { return RouterKey }

func (msg MsgCreateCandyMachine) Type() string { return TypeMsgCreateCandyMachine }

func (msg MsgCreateCandyMachine) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateCandyMachine) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateCandyMachine) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgUpdateCandyMachine{}

func NewMsgUpdateCandyMachine(sender sdk.AccAddress, machine CandyMachine,
) *MsgUpdateCandyMachine {
	return &MsgUpdateCandyMachine{
		Sender:  sender.String(),
		Machine: machine,
	}
}

func (msg MsgUpdateCandyMachine) Route() string { return RouterKey }

func (msg MsgUpdateCandyMachine) Type() string { return TypeMsgUpdateCandyMachine }

func (msg MsgUpdateCandyMachine) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateCandyMachine) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgUpdateCandyMachine) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgCloseCandyMachine{}

func NewMsgCloseCandyMachine(sender sdk.AccAddress, collId uint64,
) *MsgCloseCandyMachine {
	return &MsgCloseCandyMachine{
		Sender: sender.String(),
		CollId: collId,
	}
}

func (msg MsgCloseCandyMachine) Route() string { return RouterKey }

func (msg MsgCloseCandyMachine) Type() string { return TypeMsgCloseCandyMachine }

func (msg MsgCloseCandyMachine) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCloseCandyMachine) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCloseCandyMachine) GetSigners() []sdk.AccAddress {
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
