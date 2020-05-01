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
	Name          string         `json:"name" yaml:"name"`
	Uri           string         `json:"uri" yaml:"uri"`
	MetaUri       string         `json:"meta_uri" yaml:"meta_uri"`
	ContentUri    string         `json:"content_uri" yaml:"content_uri"`
	Denom         string         `json:"denom" yaml:"denom"`
	StreamPrice   string         `json:"stream_price" yaml:"stream_price"`
	DownloadPrice string         `json:"download_price" yaml:"download_price"`
	Creator       sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgAddContent(name, uri, metaUri, contentUri, denom, streamPrice, downloadPrice string, creator sdk.AccAddress) MsgAddContent {
	return MsgAddContent{
		Name:          name,
		Uri:           uri,
		MetaUri:       metaUri,
		ContentUri:    contentUri,
		Denom:         denom,
		StreamPrice:   streamPrice,
		DownloadPrice: downloadPrice,
		Creator:       creator,
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

	if _, err := sdk.ParseCoin(msg.StreamPrice); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid stream-price %s", msg.StreamPrice))
	}

	if _, err := sdk.ParseCoin(msg.DownloadPrice); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid download-price %s", msg.DownloadPrice))
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

var _ sdk.Msg = MsgStream{}

type MsgStream struct {
	Uri  string         `json:"uri" yaml:"uri"`
	From sdk.AccAddress `json:"from" yaml:"from"`
}

func NewMsgStream(uri string, from sdk.AccAddress) MsgStream {
	return MsgStream{
		Uri:  uri,
		From: from,
	}
}

func (msg MsgStream) Route() string { return RouterKey }
func (msg MsgStream) Type() string  { return TypeMsgAddContent }

func (msg MsgStream) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Uri)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("uri cannot be empty"))
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("Invalid from: %s", msg.From.String()))
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgStream) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgStream) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgStream) String() string {
	return fmt.Sprintf(`Msg Stream
Uri: %s
From: %s`,
		msg.Uri, msg.From,
	)
}

var _ sdk.Msg = MsgDownload{}

type MsgDownload struct {
	Uri  string         `json:"uri" yaml:"uri"`
	From sdk.AccAddress `json:"from" yaml:"from"`
}

func NewMsgDownload(uri string, from sdk.AccAddress) MsgDownload {
	return MsgDownload{
		Uri:  uri,
		From: from,
	}
}

func (msg MsgDownload) Route() string { return RouterKey }
func (msg MsgDownload) Type() string  { return TypeMsgAddContent }

func (msg MsgDownload) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Uri)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("uri cannot be empty"))
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("Invalid from: %s", msg.From.String()))
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDownload) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDownload) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgDownload) String() string {
	return fmt.Sprintf(`Msg Download
Uri: %s
From: %s`,
		msg.Uri, msg.From,
	)
}
