package types

import (
	"fmt"
	"time"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Song is a struct that contains all the metadata of a title
type Song struct {
	SongId  	uint64 		    `json:"song_id"`
	Owner   	sdk.AccAddress	`json:"owner"`
	Title 		string         	`json:"title"`
	CreateTime 	time.Time       `json:"create_time"`
}

func (s Song) String() string {
	return fmt.Sprintf(`Song %d:
		  Owner:			%s
		  Title:			%s
		  Create Time:		%s`, s.SongId, s.Owner, s.Title, s.CreateTime)
}

// Songs is an array of song
type Songs []*Song

func (songs Songs) String() string {
	out := fmt.Sprintf("%10s - (%15s) - (%40s) - [%10s] - Create Time\n", "ID", "Title", "Owner", "CreateTime")
	for _, song := range songs {
		out += fmt.Sprintf("%10d - (%15s) - (%40s) - [%10s]\n",
			song.SongId, song.Title, song.Owner, song.CreateTime)
	}

	return strings.TrimSpace(out)
}