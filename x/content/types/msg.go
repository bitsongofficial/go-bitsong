package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

// Content messages types and routes
const (
	TypeMsgAddContent = "add_content"
)

var _ sdk.Msg = MsgAddContent{}

type MsgAddContent struct {
	Name       string         `json:"name" yaml:"name"`
	Uri        string         `json:"uri" yaml:"uri"`
	MetaUri    string         `json:"meta_uri" yaml:"meta_uri"`
	ContentUri string         `json:"content_uri" yaml:"content_uri"`
	Denom      string         `json:"denom" yaml:"denom"`
	Creator    sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgAddContent(name, uri, metaUri, contentUri, denom string, creator sdk.AccAddress) MsgAddContent {
	return MsgAddContent{
		Name:       name,
		Uri:        uri,
		MetaUri:    metaUri,
		ContentUri: contentUri,
		Denom:      denom,
		Creator:    creator,
	}
}

func (msg MsgAddContent) Route() string { return RouterKey }
func (msg MsgAddContent) Type() string  { return TypeMsgAddContent }

func (msg MsgAddContent) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Name)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("name cannot be empty"))
	}

	if len(msg.Name) > MaxNameLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("name cannot be longer than %d characters", MaxUriLength))
	}

	if len(strings.TrimSpace(msg.Uri)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("uri cannot be empty"))
	}

	if len(msg.Uri) > MaxUriLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("uri cannot be longer than %d characters", MaxUriLength))
	}

	if len(strings.TrimSpace(msg.MetaUri)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("meta-uri cannot be empty"))
	}

	if len(msg.MetaUri) > MaxUriLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("meta-uri cannot be longer than %d characters", MaxUriLength))
	}

	if len(strings.TrimSpace(msg.ContentUri)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("content-uri cannot be empty"))
	}

	if len(msg.ContentUri) > MaxUriLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("content-uri cannot be longer than %d characters", MaxUriLength))
	}

	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("%s", err.Error()))
	}

	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("Invalid creator: %s", msg.Creator.String()))
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddContent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddContent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg MsgAddContent) String() string {
	return fmt.Sprintf(`Msg Add Content
Name: %s
Uri: %s
MetaUri: %s
ContentUri: %s
Denom: %s
Creator: %s`,
		msg.Name, msg.Uri, msg.MetaUri, msg.ContentUri, msg.Denom, msg.Creator,
	)
}
