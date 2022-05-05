package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "fantoken"

	TypeMsgIssueFanToken         = "issue_fan_token"
	TypeMsgEditFanToken          = "edit_fan_token_mintable"
	TypeMsgMintFanToken          = "mint_fan_token"
	TypeMsgBurnFanToken          = "burn_fan_token"
	TypeMsgTransferFanTokenOwner = "transfer_fan_token_owner"

	// DoNotModify used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"
)

var (
	_ sdk.Msg = &MsgIssueFanToken{}
	_ sdk.Msg = &MsgEditFanToken{}
	_ sdk.Msg = &MsgMintFanToken{}
	_ sdk.Msg = &MsgBurnFanToken{}
	_ sdk.Msg = &MsgTransferFanTokenOwner{}
)

// NewMsgIssueToken - construct token issue msg.
func NewMsgIssueFanToken(
	symbol string, name string,
	maxSupply sdk.Int,
	descriptioin string, owner, uri string, issueFee sdk.Coin,
) *MsgIssueFanToken {
	return &MsgIssueFanToken{
		Symbol:      symbol,
		Name:        name,
		MaxSupply:   maxSupply,
		Description: descriptioin,
		Owner:       owner,
		URI:         uri,
		IssueFee:    issueFee,
	}
}

// Route Implements Msg.
func (msg MsgIssueFanToken) Route() string { return MsgRoute }

// Type Implements Msg.
func (msg MsgIssueFanToken) Type() string { return TypeMsgIssueFanToken }

// ValidateBasic Implements Msg.
func (msg MsgIssueFanToken) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	denom := GetFantokenDenom(owner, msg.Symbol, msg.Name)
	denomMetaData := banktypes.Metadata{
		Description: msg.Description,
		Base:        denom,
		Display:     msg.Symbol,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: denom, Exponent: 0},
			{Denom: msg.Symbol, Exponent: FanTokenDecimal},
		},
	}

	return ValidateToken(
		NewFanToken(
			msg.Name,
			msg.MaxSupply,
			owner,
			msg.URI,
			denomMetaData,
		),
	)
}

// GetSignBytes Implements Msg.
func (msg MsgIssueFanToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgIssueFanToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferTokenOwner return a instance of MsgTransferTokenOwner
func NewMsgTransferFanTokenOwner(denom, srcOwner, dstOwner string) *MsgTransferFanTokenOwner {
	return &MsgTransferFanTokenOwner{
		Denom:    denom,
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
	}
}

// GetSignBytes implements Msg
func (msg MsgTransferFanTokenOwner) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferFanTokenOwner) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.SrcOwner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgTransferFanTokenOwner) ValidateBasic() error {
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
func (msg MsgTransferFanTokenOwner) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferFanTokenOwner) Type() string { return TypeMsgTransferFanTokenOwner }

// NewMsgEditToken creates a MsgEditToken
func NewMsgEditFanToken(denom string, mintable bool, owner string) *MsgEditFanToken {
	return &MsgEditFanToken{
		Denom:    denom,
		Mintable: mintable,
		Owner:    owner,
	}
}

// Route implements Msg
func (msg MsgEditFanToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEditFanToken) Type() string { return TypeMsgEditFanToken }

// GetSignBytes implements Msg
func (msg MsgEditFanToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEditFanToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgEditFanToken) ValidateBasic() error {
	// check owner
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// check symbol
	return ValidateDenom(msg.Denom)
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgMintFanToken(recipient, denom, owner string, amount sdk.Int) *MsgMintFanToken {
	return &MsgMintFanToken{
		Recipient: recipient,
		Denom:     denom,
		Owner:     owner,
		Amount:    amount,
	}
}

// Route implements Msg
func (msg MsgMintFanToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgMintFanToken) Type() string { return TypeMsgMintFanToken }

// GetSignBytes implements Msg
func (msg MsgMintFanToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgMintFanToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgMintFanToken) ValidateBasic() error {
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

// NewMsgBurnToken creates a MsgMintToken
func NewMsgBurnFanToken(denom string, owner string, amount sdk.Int) *MsgBurnFanToken {
	return &MsgBurnFanToken{
		Denom:  denom,
		Amount: amount,
		Sender: owner,
	}
}

// Route implements Msg
func (msg MsgBurnFanToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgBurnFanToken) Type() string { return TypeMsgBurnFanToken }

// GetSignBytes implements Msg
func (msg MsgBurnFanToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgBurnFanToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgBurnFanToken) ValidateBasic() error {
	// check the owner
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err := ValidateAmount(msg.Amount); err != nil {
		return err
	}

	return ValidateDenom(msg.Denom)
}
