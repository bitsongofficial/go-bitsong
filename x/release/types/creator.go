package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreatorRelease struct {
	Creator    sdk.AccAddress `json:"creator"`
	ReleaseIDs []string       `json:"release_ids"`
}

func NewCreatorRelease(creator sdk.AccAddress, releaseIDs []string) CreatorRelease {
	return CreatorRelease{
		Creator:    creator,
		ReleaseIDs: releaseIDs,
	}
}

type CreatorReleases []CreatorRelease
