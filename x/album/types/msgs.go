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
	TypeMsgAddTrack    = "add_track"
	TypeMsgDeposit     = "deposit"
)

/****************************************
 * MsgCreateAlbum
 ****************************************/

var _ sdk.Msg = MsgCreateAlbum{}

// MsgCreateAlbum defines CreateAlbum message
type MsgCreateAlbum struct {
	AlbumType   AlbumType      `json:"album_type"` // The type of the album: one of 'album', 'single', or 'compilation'.
	Title       string         `json:"title"`      // Album name
	MetadataURI string         `json:"metadata_uri"`
	Owner       sdk.AccAddress `json:"owner"` // Album owner
}

func NewMsgCreateAlbum(albumType AlbumType, title string, metadataUri string, owner sdk.AccAddress) MsgCreateAlbum {
	return MsgCreateAlbum{
		AlbumType:   albumType,
		Title:       title,
		MetadataURI: metadataUri,
		Owner:       owner,
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
	// - Add more check for CID (Metadata uri ipfs:)
	if len(strings.TrimSpace(msg.MetadataURI)) == 0 {
		return ErrInvalidAlbumMetadataURI(DefaultCodespace, "artist metadata uri cannot be blank")
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
  Metadata URI: %s
  Owner: %s
`, msg.AlbumType.String(), msg.Title, msg.MetadataURI, msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateAlbum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateAlbum) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

/****************************************
 * MsgAddTrackAlbum
 ****************************************/

var _ sdk.Msg = MsgAddTrackAlbum{}

// MsgAddTrackAlbum defines AddTrackAlbum message
type MsgAddTrackAlbum struct {
	AlbumID uint64         "json:`album_id` yaml:`album_id`"
	TrackID uint64         "json:`track_id` yaml:`track_id`"
	Owner   sdk.AccAddress `json:"owner"` // Artist owner
}

func NewMsgAddTrackAlbum(albumID uint64, trackID uint64, owner sdk.AccAddress) MsgAddTrackAlbum {
	return MsgAddTrackAlbum{
		AlbumID: albumID,
		TrackID: trackID,
		Owner:   owner,
	}
}

//nolint
func (msg MsgAddTrackAlbum) Route() string { return RouterKey }
func (msg MsgAddTrackAlbum) Type() string  { return TypeMsgAddTrack }

// ValidateBasic
func (msg MsgAddTrackAlbum) ValidateBasic() sdk.Error {
	// TODO:
	// - improve check

	if msg.AlbumID == 0 {
		return ErrUnknownAlbum(DefaultCodespace, "album-id cannot be blank")
	}

	if msg.TrackID == 0 {
		return ErrUnknownTrack(DefaultCodespace, "track-id cannot be blank")
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

// Implements Msg.
func (msg MsgAddTrackAlbum) String() string {
	return fmt.Sprintf(`Add Track Album Message:
  AlbumID: %d
  TrackID: %d
  Address: %s
`, msg.AlbumID, msg.TrackID, msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgAddTrackAlbum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddTrackAlbum) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

/****************************************
 * MsgDeposit
 ****************************************/

var _ sdk.Msg = MsgDeposit{}

type MsgDeposit struct {
	AlbumID   uint64         `json:"album_id" yaml:"album_id"`   // ID of the album
	Depositor sdk.AccAddress `json:"depositor" yaml:"depositor"` // Address of the depositor
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`       // Coins to add to the proposal's deposit
}

func NewMsgDeposit(depositor sdk.AccAddress, albumID uint64, amount sdk.Coins) MsgDeposit {
	return MsgDeposit{albumID, depositor, amount}
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Route() string { return RouterKey }
func (msg MsgDeposit) Type() string  { return TypeMsgDeposit }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

func (msg MsgDeposit) String() string {
	return fmt.Sprintf(`Deposit Message:
  Depositer:   %s
  Album ID: %d
  Amount:      %s
`, msg.Depositor, msg.AlbumID, msg.Amount)
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}
