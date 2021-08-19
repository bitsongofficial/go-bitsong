package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Channel struct {
	Owner       sdk.AccAddress `json:"owner"`
	Handle      string         `json:"handle"`
	MetadataURI string         `json:"metadata_uri"`
	CreatedAt   time.Time      `json:"created_at"`
}

func NewChannel(owner sdk.AccAddress, handle, metadataURI string, createdAt time.Time) Channel {
	return Channel{
		Owner:       owner,
		Handle:      handle,
		MetadataURI: metadataURI,
		CreatedAt:   createdAt,
	}
}

func (p Channel) String() string {
	return fmt.Sprintf(`
  Owner:       %s
  Handle:        %s
  MetadataURI:   %s
  CreatedAt:     %s`,
		p.Owner.String(), p.Handle, p.MetadataURI, p.CreatedAt.String())
}

func (p Channel) Validate() error {
	if p.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel owner cannot be empty")
	}

	if p.Handle == "" || len(p.Handle) < 3 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel handle must have a length > 3")
	}

	if len(p.MetadataURI) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel metadataURI cannot be more than 256 characters")
	}

	return nil
}
