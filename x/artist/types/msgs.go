package types

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

const (
	TypeMsgArtistCreate = "artist_create"
)

var _ sdk.Msg = MsgArtistCreate{}

type MsgArtistCreate struct {
	ID          string         `json:"id" yaml:"id"`
	Name        string         `json:"name" yaml:"name"`
	URLs        btsg.URLs      `json:"urls" yaml:"urls"`
	Genres      []string       `json:"genres" yaml:"genres"`
	MetadataURI string         `json:"metadata_uri" yaml:"metadata_uri"`
	Creator     sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgArtistCreate(id, name string, urls btsg.URLs, genres []string, metadataURI string, creator sdk.AccAddress) MsgArtistCreate {
	return MsgArtistCreate{
		ID:          id,
		Name:        name,
		URLs:        urls,
		Genres:      genres,
		MetadataURI: metadataURI,
		Creator:     creator,
	}
}

func (msg MsgArtistCreate) Route() string { return RouterKey }
func (msg MsgArtistCreate) Type() string  { return TypeMsgArtistCreate }

func (msg MsgArtistCreate) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "creator cannot be empty")
	}

	if strings.TrimSpace(msg.ID) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "id cannot be empty")
	}

	if strings.TrimSpace(msg.Name) == "" && len(strings.TrimSpace(msg.Name)) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "id cannot be empty or more than 256 characters")
	}

	if len(msg.MetadataURI) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel metadataURI cannot be more than 256 characters")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgArtistCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgArtistCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg MsgArtistCreate) String() string {
	return fmt.Sprintf(`Msg Artist Create
Creator: %s,
Id: %s
Name: %s`,
		msg.Creator.String(), msg.ID, msg.Name,
	)
}
