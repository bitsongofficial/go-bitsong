package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryRelease              = "release"
	QueryAllReleaseForCreator = "allReleaseForCreator"
)

type QueryReleaseParams struct {
	ReleaseID string
}

func NewQueryReleaseParams(releaseID string) QueryReleaseParams {
	return QueryReleaseParams{ReleaseID: releaseID}
}

type QueryAllReleaseForCreatorParams struct {
	Creator sdk.AccAddress
}

func NewQueryAllReleaseForCreatorParams(creator sdk.AccAddress) QueryAllReleaseForCreatorParams {
	return QueryAllReleaseForCreatorParams{Creator: creator}
}
