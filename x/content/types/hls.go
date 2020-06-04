package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type HlsInfo struct {
	HlsCode []byte         `json:"hls_code"`
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewHlsInfo(hlsCode []byte, creator sdk.AccAddress) HlsInfo {
	return HlsInfo{
		HlsCode: hlsCode,
		Creator: creator,
	}
}
