package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/****
 * Album Msg
 ***/

// Album messages types and routes
const (
	TypeMsgCreateAlbum = "create_album"
)

/****************************************
 * MsgCreateAlbum
 ****************************************/

var _ sdk.Msg = MsgCreateAlbum{}

// MsgCreateAlbum defines CreateAlbum message
type MsgCreateAlbum struct {
	AlbumType            AlbumType      `json:"album_type"`             // The type of the album: one of 'album', 'single', or 'compilation'.
	Title                string         `json:"title"`                  // Artist name
	ReleaseDate          string         `json:"release_date"`           // The date the album was first released, for example '1981-12-15'. Depending on the precision, it might be shown as '1981' or '1981-12'.
	ReleaseDatePrecision string         `json:"release_date_precision"` // The precision with which release_date value is known: 'year', 'month', or 'day'.
	Owner                sdk.AccAddress `json:"owner"`                  // Artist owner
}

func NewMsgCreateAlbum(albumType AlbumType, title string, releaseDate string, releasePrecision string, owner sdk.AccAddress) MsgCreateAlbum {
	return MsgCreateAlbum{
		AlbumType:            albumType,
		Title:                title,
		ReleaseDate:          releaseDate,
		ReleaseDatePrecision: releasePrecision,
		Owner:                owner,
	}
}

//nolint
func (msg MsgCreateAlbum) Route() string { return RouterKey }
func (msg MsgCreateAlbum) Type() string  { return TypeMsgCreateAlbum }

// ValidateBasic
func (msg MsgCreateAlbum) ValidateBasic() sdk.Error {
	if !msg.AlbumType.Valid() {
		return ErrInvalidAlbumType(DefaultCodespace, "album type is not valid")
	}

	if len(strings.TrimSpace(msg.Title)) == 0 {
		return ErrInvalidAlbumName(DefaultCodespace, "album name cannot be blank")
	}

	if len(msg.Title) > MaxTitleLength {
		return ErrInvalidAlbumName(DefaultCodespace, fmt.Sprintf("album name is longer than max length of %d", MaxTitleLength))
	}

	// TODO:
	// - improve check on release_date
	// - improve check on release_date_precision
	if len(strings.TrimSpace(msg.ReleaseDate)) == 0 {
		return ErrInvalidAlbumReleaseDate(DefaultCodespace, "album release date cannot be blank")
	}

	if len(strings.TrimSpace(msg.ReleaseDatePrecision)) == 0 {
		return ErrInvalidAlbumReleaseDatePrecision(DefaultCodespace, "album release date precision cannot be blank")
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

// Implements Msg.
func (msg MsgCreateAlbum) String() string {
	return fmt.Sprintf(`Create Album Message:
  Album Type:         %s
  Title: %s
  Release Date: %s
  Release Date Precision: %s
  Owner: %s
`, msg.AlbumType.String(), msg.Title, msg.ReleaseDate, msg.ReleaseDatePrecision, msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateAlbum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateAlbum) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
