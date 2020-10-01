package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"
)

type Release struct {
	ReleaseID   string         `json:"release_id"`
	MetadataURI string         `json:"metadata_uri"`
	Creator     sdk.AccAddress `json:"creator"`
	CreatedAt   time.Time      `json:"created_at"`
}

func NewRelease(releaseID, metadataURI string, creator sdk.AccAddress, createdAt time.Time) Release {
	return Release{
		ReleaseID:   releaseID,
		MetadataURI: metadataURI,
		Creator:     creator,
		CreatedAt:   createdAt,
	}
}

func (r Release) String() string {
	return fmt.Sprintf(`
  ReleaseID:   %s
  MetadataURI: %s
  Creator:     %s
  CreatedAt:   %s`,
		r.ReleaseID, r.MetadataURI, r.Creator.String(), r.CreatedAt.String())
}

func (r Release) Validate() error {
	if r.ReleaseID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "releaseID is required")
	}

	if len(r.MetadataURI) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "relese metadataURI cannot be more than 256 characters")
	}

	if r.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "creator address cannot be empty")
	}

	return nil
}
