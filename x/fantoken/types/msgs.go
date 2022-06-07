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
	_ sdk.Msg = &MsgDisableMint{}
	_ sdk.Msg = &MsgMint{}
	_ sdk.Msg = &MsgBurn{}
	_ sdk.Msg = &MsgTransferAuthority{}
)

// NewMsgIssue - construct token issue msg.
func NewMsgIssue(name, symbol, uri string, maxSupply sdk.Int, authority string) *MsgIssue {
	return &MsgIssue{
		Name:      name,
		Symbol:    symbol,
		URI:       uri,
		MaxSupply: maxSupply,
		Authority: authority,
	}
}

// Route Implements Msg.
func (msg MsgIssue) Route() string { return MsgRoute }

// Type Implements Msg.
func (msg MsgIssue) Type() string { return TypeMsgIssue }

// ValidateBasic Implements Msg.
func (msg MsgIssue) ValidateBasic() error {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	fantoken := &FanToken{
		MaxSupply: msg.MaxSupply,
		Mintable:  true,
		Authority: authority.String(),
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
	from, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferAuthority return a instance of MsgTransferAuthority
func NewMsgTransferAuthority(denom, srcAuthority, dstAuthority string) *MsgTransferAuthority {
	return &MsgTransferAuthority{
		Denom:        denom,
		SrcAuthority: srcAuthority,
		DstAuthority: dstAuthority,
	}
}

// GetSignBytes implements Msg
func (msg MsgTransferAuthority) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferAuthority) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.SrcAuthority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgTransferAuthority) ValidateBasic() error {
	srcAuthority, err := sdk.AccAddressFromBech32(msg.SrcAuthority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid source authority address (%s)", err)
	}

	dstAuthority, err := sdk.AccAddressFromBech32(msg.DstAuthority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination authority address (%s)", err)
	}

	// check if the `DstOwner` is same as the original owner
	if srcAuthority.Equals(dstAuthority) {
		return ErrInvalidToAddress
	}

	// check the symbol
	if err := ValidateDenom(msg.Denom); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgTransferAuthority) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferAuthority) Type() string { return TypeMsgTransferOwnership }

// NewMsgDisableMint creates a MsgDisableMint
func NewMsgDisableMint(denom string, authority string) *MsgDisableMint {
	return &MsgDisableMint{
		Denom:     denom,
		Authority: authority,
	}
}

// Route implements Msg
func (msg MsgDisableMint) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgDisableMint) Type() string { return TypeMsgEdit }

// GetSignBytes implements Msg
func (msg MsgDisableMint) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgDisableMint) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgDisableMint) ValidateBasic() error {
	// check owner
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	return ValidateDenom(msg.Denom)
}

// NewMsgMint creates a MsgMint
func NewMsgMint(recipient, denom, authority string, amount sdk.Int) *MsgMint {
	return &MsgMint{
		Recipient: recipient,
		Denom:     denom,
		Authority: authority,
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
	from, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgMint) ValidateBasic() error {
	// check the authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
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
