package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Song is a struct that contains all the metadata of a title
type Song struct {
	SongID                  uint64         `json:"song_id"`
	Owner                   sdk.AccAddress `json:"owner"`
	Title                   string         `json:"title"`
	Content                 string         `json:"content"`
	TotalReward             sdk.Int        `json:"total_reward"`
	RedistributionSplitRate string         `json:"redistribution_split_rate"`
	CreateTime              time.Time      `json:"create_time"`
}

func (s Song) String() string {
	return fmt.Sprintf(`Song %d:
		  Owner:					%s
		  Title:					%s
		  Content:					%s
		  TotalReward:				%s
		  RedistributionSplitRate:	%s
		  Create Time:		%s`, s.SongID, s.Owner, s.Title, s.Content, s.TotalReward, s.RedistributionSplitRate, s.CreateTime)
}

// Songs is an array of song
// To FIX with new fields
type Songs []*Song

func (songs Songs) String() string {
	out := fmt.Sprintf("%10s - (%15s) - (%40s) - [%10s] - Create Time\n", "ID", "Title", "Owner", "CreateTime")
	for _, song := range songs {
		out += fmt.Sprintf("%10d - (%15s) - (%40s) - [%10s]\n",
			song.SongID, song.Title, song.Owner, song.CreateTime)
	}

	return strings.TrimSpace(out)
}
