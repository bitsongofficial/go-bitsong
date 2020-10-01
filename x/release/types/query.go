package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryRelease          = "release"
	QueryReleaseByCreator = "releaseByCreator"
)

type QueryReleaseParams struct {
	ReleaseID string
}

func NewQueryReleaseParams(releaseID string) QueryReleaseParams {
	return QueryReleaseParams{ReleaseID: releaseID}
}

type QueryByCreatorParams struct {
	Creator sdk.AccAddress
}

func NewQueryByCreatorParams(creator sdk.AccAddress) QueryByCreatorParams {
	return QueryByCreatorParams{Creator: creator}
}
