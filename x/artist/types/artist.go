package types

import (
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

type Artist struct {
	ID          btsg.ID        `json:"id"`
	Name        string         `json:"name"`
	URLs        btsg.URLs      `json:"urls"`
	Genres      []string       `json:"genres"`
	MetadataURI string         `json:"metadata"`
	Creator     sdk.AccAddress `json:"creator"`
}

func NewArtist(id btsg.ID, name string, urls btsg.URLs, genres []string, metadataURI string, creator sdk.AccAddress) Artist {
	return Artist{
		ID:          id,
		Name:        name,
		URLs:        urls,
		Genres:      genres,
		MetadataURI: metadataURI,
		Creator:     creator,
	}
}

func (a Artist) Validate() error {
	if a.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "creator cannot be empty")
	}

	if strings.TrimSpace(a.ID.String()) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "id cannot be empty")
	}

	if strings.TrimSpace(a.Name) == "" && len(strings.TrimSpace(a.Name)) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "id cannot be empty or more than 256 characters")
	}

	if len(a.MetadataURI) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel metadataURI cannot be more than 256 characters")
	}

	return nil
}
