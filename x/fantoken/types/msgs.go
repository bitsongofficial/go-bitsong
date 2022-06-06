package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "fantoken"

	TypeMsgIssue             = "issue"
	TypeMsgEdit              = "edit"
	TypeMsgMint              = "mint"
	TypeMsgBurn              = "burn"
	TypeMsgTransferOwnership = "transfer_ownership"
)

var (
	_ sdk.Msg = &MsgIssue{}
	_ sdk.Msg = &MsgEdit{}
	_ sdk.Msg = &MsgMint{}
	_ sdk.Msg = &MsgBurn{}
	_ sdk.Msg = &MsgTransferOwnership{}
)

// NewMsgIssue - construct token issue msg.
func NewMsgIssue(name, symbol, uri string, maxSupply sdk.Int, owner string) *MsgIssue {
	return &MsgIssue{
		Name:      name,
		Symbol:    symbol,
		URI:       uri,
		MaxSupply: maxSupply,
		Owner:     owner,
	}
}

// Route Implements Msg.
func (msg MsgIssue) Route() string { return MsgRoute }

// Type Implements Msg.
func (msg MsgIssue) Type() string { return TypeMsgIssue }

// ValidateBasic Implements Msg.
func (msg MsgIssue) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	fantoken := &FanToken{
		MaxSupply: msg.MaxSupply,
		Mintable:  true,
		Owner:     owner.String(),
		MetaData: Metadata{
			Name:   msg.Name,
			Symbol: msg.Symbol,
			URI:    msg.URI,
		},
	}

	return ValidateFanToken(fantoken)
}

// GetSignBytes Implements Msg.
func (msg MsgIssue) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgIssue) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferOwnership return a instance of MsgTransferOwnership
func NewMsgTransferOwnership(denom, srcOwner, dstOwner string) *MsgTransferOwnership {
	return &MsgTransferOwnership{
		Denom:    denom,
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
	}
}

// GetSignBytes implements Msg
func (msg MsgTransferOwnership) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferOwnership) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.SrcOwner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgTransferOwnership) ValidateBasic() error {
	srcOwner, err := sdk.AccAddressFromBech32(msg.SrcOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid source owner address (%s)", err)
	}

	dstOwner, err := sdk.AccAddressFromBech32(msg.DstOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination owner address (%s)", err)
	}

	// check if the `DstOwner` is same as the original owner
	if srcOwner.Equals(dstOwner) {
		return ErrInvalidToAddress
	}

	// check the symbol
	if err := ValidateDenom(msg.Denom); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgTransferOwnership) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferOwnership) Type() string { return TypeMsgTransferOwnership }

// NewMsgEdit creates a MsgEdit
func NewMsgEdit(denom string, mintable bool, owner string) *MsgEdit {
	return &MsgEdit{
		Denom:    denom,
		Mintable: mintable,
		Owner:    owner,
	}
}

// Route implements Msg
func (msg MsgEdit) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEdit) Type() string { return TypeMsgEdit }

// GetSignBytes implements Msg
func (msg MsgEdit) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEdit) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgEdit) ValidateBasic() error {
	// check owner
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return ValidateDenom(msg.Denom)
}

// NewMsgMint creates a MsgMint
func NewMsgMint(recipient, denom, owner string, amount sdk.Int) *MsgMint {
	return &MsgMint{
		Recipient: recipient,
		Denom:     denom,
		Owner:     owner,
		Amount:    amount,
	}
}

// Route implements Msg
func (msg MsgMint) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgMint) Type() string { return TypeMsgMint }

// GetSignBytes implements Msg
func (msg MsgMint) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgMint) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgMint) ValidateBasic() error {
	// check the owner
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// check the reception
	if len(msg.Recipient) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid mint reception address (%s)", err)
		}
	}

	if err := ValidateAmount(msg.Amount); err != nil {
		return err
	}

	return ValidateDenom(msg.Denom)
}

// NewMsgBurn creates a MsgBurn
func NewMsgBurn(denom, owner string, amount sdk.Int) *MsgBurn {
	return &MsgBurn{
		Denom:  denom,
		Amount: amount,
		Sender: owner,
	}
}

// Route implements Msg
func (msg MsgBurn) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgBurn) Type() string { return TypeMsgBurn }

// GetSignBytes implements Msg
func (msg MsgBurn) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgBurn) ValidateBasic() error {
	// check the owner
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err := ValidateAmount(msg.Amount); err != nil {
		return err
	}

	return ValidateDenom(msg.Denom)
}
