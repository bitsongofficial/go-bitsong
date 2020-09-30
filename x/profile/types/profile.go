package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"
)

type Profile struct {
	Address     sdk.AccAddress `json:"address"`
	Handle      string         `json:"handle"`
	MetadataURI string         `json:"metadata_uri"`
	CreatedAt   time.Time      `json:"created_at"`
}

func NewProfile(accAddr sdk.AccAddress, handle, metadataURI string, createdAt time.Time) Profile {
	return Profile{
		Address:     accAddr,
		Handle:      handle,
		MetadataURI: metadataURI,
		CreatedAt:   createdAt,
	}
}

func (p Profile) String() string {
	return fmt.Sprintf(`
  Address:       %s
  Handle:        %s
  MetadataURI:   %s
  CreatedAt:     %s`,
		p.Address.String(), p.Handle, p.MetadataURI, p.CreatedAt.String())
}

func (p Profile) Validate() error {
	if p.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "profile address cannot be empty")
	}

	if p.Handle == "" || len(p.Handle) < 3 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "profile handle must have a length > 3")
	}

	if len(p.Handle) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "profile metadataURI cannot be more than 256 characters")
	}

	return nil
}
