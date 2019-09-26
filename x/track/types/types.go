package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Track is a struct that contains all the metadata of a single track
type Track struct {
	TrackID                 uint64         `json:"track_id"`
	Owner                   sdk.AccAddress `json:"owner"`
	Title                   string         `json:"title"`
	Content                 string         `json:"content"`
	TotalReward             sdk.Int        `json:"total_reward"`
	RedistributionSplitRate sdk.Dec        `json:"redistribution_split_rate"`
	CreateTime              time.Time      `json:"create_time"`
}

func (t Track) String() string {
	return fmt.Sprintf(`Track %d:
		  Owner:					%s
		  Title:					%s
		  Content:					%s
		  TotalReward:				%s
		  RedistributionSplitRate:	%d
		  Create Time:		%s`, t.TrackID, t.Owner, t.Title, t.Content, t.TotalReward, t.RedistributionSplitRate, t.CreateTime)
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

// Play is a struct that contains all the metadata of a single play
type Play struct {
	AccAddress sdk.AccAddress `json:"acc_address"`
	TrackId    uint64         `json:"track_id"`
	Shares     sdk.Dec        `json:"shares"`
	Streams    sdk.Int        `json:"streams"`
	CreateTime time.Time      `json:"create_time"`
}

func (p Play) String() string {
	return fmt.Sprintf(`AccAddress %s:
		  TrackId:					%d
		  Shares:					%d
		  Streams:					%s
		  CreateTime:				%s`, p.AccAddress, p.TrackId, p.Shares, p.Streams, p.CreateTime)
}

// current rewards and current period for an account
// kept as a running counter and incremented each block
// as long as the account's tokens remain constant
type AccountCurrentRewards struct {
	Rewards sdk.DecCoins `json:"rewards" yaml:"rewards"` // current rewards
	Period  uint64       `json:"period" yaml:"period"`   // current period
}
