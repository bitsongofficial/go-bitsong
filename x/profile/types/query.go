package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryProfile          = "profile"
	QueryProfileByAddress = "profileByAddress"
)

type QueryProfileParams struct {
	Handle string
}

func NewQueryProfileParams(handle string) QueryProfileParams {
	return QueryProfileParams{Handle: handle}
}

type QueryByAddressParams struct {
	Address sdk.AccAddress
}

func NewQueryByAddressParams(addr sdk.AccAddress) QueryByAddressParams {
	return QueryByAddressParams{Address: addr}
}
