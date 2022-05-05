package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryFanToken  = "fantoken"
	QueryFanTokens = "fantokens"
	QueryParams    = "params"
	QueryTotalBurn = "total_burn"
)

// QueryFanTokenParams is the query parameters for 'custom/fantoken/token'
type QueryFanTokenParams struct {
	Denom string
}

// QueryFanTokensParams is the query parameters for 'custom/fantoken/tokens'
type QueryFanTokensParams struct {
	Owner sdk.AccAddress
}
