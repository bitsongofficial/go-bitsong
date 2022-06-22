package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "fantoken"

	TypeMsgIssue        = "issue"
	TypeMsgEdit         = "edit"
	TypeMsgMint         = "mint"
	TypeMsgBurn         = "burn"
	TypeMsgSetAuthority = "set_authority"
	TypeMsgSetMinter    = "set_minter"
)

var (
	_ sdk.Msg = &MsgIssue{}
	_ sdk.Msg = &MsgDisableMint{}
	_ sdk.Msg = &MsgMint{}
	_ sdk.Msg = &MsgBurn{}
	_ sdk.Msg = &MsgSetAuthority{}
	_ sdk.Msg = &MsgSetMinter{}
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

	minter, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
	}

	fantoken := &FanToken{
		MaxSupply: msg.MaxSupply,
		Minter:    minter.String(),
		MetaData: Metadata{
			Name:      msg.Name,
			Symbol:    msg.Symbol,
			URI:       msg.URI,
			Authority: authority.String(),
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

// NewMsgSetAuthority return a instance of MsgSetAuthority
func NewMsgSetAuthority(denom, oldAuthority, newAuthority string) *MsgSetAuthority {
	return &MsgSetAuthority{
		Denom:        denom,
		OldAuthority: oldAuthority,
		NewAuthority: newAuthority,
	}
}

// GetSignBytes implements Msg
func (msg MsgSetAuthority) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgSetAuthority) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.OldAuthority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgSetAuthority) ValidateBasic() error {
	oldAuthority, err := sdk.AccAddressFromBech32(msg.OldAuthority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid old authority address (%s)", err)
	}

	newAuthority, err := sdk.AccAddressFromBech32(msg.NewAuthority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new authority address (%s)", err)
	}

	// check if the `newAuthority` is same as the original authority
	if oldAuthority.Equals(newAuthority) {
		return ErrInvalidToAddress
	}

	// check the symbol
	if err := ValidateDenom(msg.Denom); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgSetAuthority) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgSetAuthority) Type() string { return TypeMsgSetAuthority }

// NewMsgSetMinter return a instance of MsgSetMinter
func NewMsgSetMinter(denom, oldMinter, newMinter string) *MsgSetMinter {
	return &MsgSetMinter{
		Denom:     denom,
		OldMinter: oldMinter,
		NewMinter: newMinter,
	}
}

// GetSignBytes implements Msg
func (msg MsgSetMinter) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgSetMinter) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.OldMinter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgSetMinter) ValidateBasic() error {
	oldMinter, err := sdk.AccAddressFromBech32(msg.OldMinter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid old minter address (%s)", err)
	}

	newMinter, err := sdk.AccAddressFromBech32(msg.NewMinter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new minter address (%s)", err)
	}

	// check if the `newMinter` is same as the original minter
	if oldMinter.Equals(newMinter) {
		return ErrInvalidToAddress
	}

	// check the symbol
	if err := ValidateDenom(msg.Denom); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgSetMinter) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgSetMinter) Type() string { return TypeMsgSetMinter }

// NewMsgDisableMint creates a MsgDisableMint
func NewMsgDisableMint(denom string, minter string) *MsgDisableMint {
	return &MsgDisableMint{
		Denom:  denom,
		Minter: minter,
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
	from, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgDisableMint) ValidateBasic() error {
	// check minter
	if _, err := sdk.AccAddressFromBech32(msg.Minter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
	}

	return ValidateDenom(msg.Denom)
}

// NewMsgMint creates a MsgMint
func NewMsgMint(recipient, denom, minter string, amount sdk.Int) *MsgMint {
	return &MsgMint{
		Recipient: recipient,
		Denom:     denom,
		Minter:    minter,
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
	from, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgMint) ValidateBasic() error {
	// check the minter
	if _, err := sdk.AccAddressFromBech32(msg.Minter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
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
