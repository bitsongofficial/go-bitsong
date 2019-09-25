package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Track is a struct that contains all the metadata of a title
type Track struct {
	TrackID                 uint64         `json:"track_id"`
	Owner                   sdk.AccAddress `json:"owner"`
	Title                   string         `json:"title"`
	Content                 string         `json:"content"`
	TotalReward             sdk.Int        `json:"total_reward"`
	RedistributionSplitRate sdk.Dec        `json:"redistribution_split_rate"`
	CreateTime              time.Time      `json:"create_time"`
}

func (s Track) String() string {
	return fmt.Sprintf(`Track %d:
		  Owner:					%s
		  Title:					%s
		  Content:					%s
		  TotalReward:				%s
		  RedistributionSplitRate:	%d
		  Create Time:		%s`, s.TrackID, s.Owner, s.Title, s.Content, s.TotalReward, s.RedistributionSplitRate, s.CreateTime)
}

// Tracks is an array of track
// To FIX with new fields
type Tracks []*Track

func (tracks Tracks) String() string {
	out := fmt.Sprintf("%10s - (%15s) - (%40s) - [%10s] - Create Time\n", "ID", "Title", "Owner", "CreateTime")
	for _, track := range tracks {
		out += fmt.Sprintf("%10d - (%15s) - (%40s) - [%10s]\n",
			track.TrackID, track.Title, track.Owner, track.CreateTime)
	}

	return strings.TrimSpace(out)
}
