package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryChannel        = "channel"
	QueryChannelByOwner = "channelByAddress"
)

type QueryChannelParams struct {
	Handle string
}

func NewQueryChannelParams(handle string) QueryChannelParams {
	return QueryChannelParams{Handle: handle}
}

type QueryByOwnerParams struct {
	Owner sdk.AccAddress
}

func NewQueryByOwnerParams(addr sdk.AccAddress) QueryByOwnerParams {
	return QueryByOwnerParams{Owner: addr}
}
