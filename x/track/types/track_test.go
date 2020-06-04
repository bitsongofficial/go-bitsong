package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTrack(t *testing.T) {
	track, err := NewTrack("My best title", []Artist{}, 1, 200, false, nil, nil, "", nil)
	require.NoError(t, err)

	fmt.Println(track.Title)
	fmt.Println(track.Cid)
}
