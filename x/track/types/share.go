package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Share struct {
	TrackID    uint64  `json:"track_id" yaml:"track_id"`
	TotalShare sdk.Dec `json:"total_share" yaml:"total_share"`
}

func NewShare(trackID uint64) Share {
	return Share{
		TrackID:    trackID,
		TotalShare: sdk.NewDec(0),
	}
}

type Shares []Share

func (s Shares) String() string {
	if len(s) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Shares on Track %d:", s[0].TrackID)
	for _, share := range s {
		out += fmt.Sprintf("\n  %d: %s", share.TrackID, share.TotalShare.String())
	}
	return out
}
