package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"kkn.fi/base62"
	"testing"
	"time"
)

func TestNewTrack(t *testing.T) {
	track, err := NewTrack("My best title", []string{}, 1, 200, false, nil, nil, "", nil)
	require.NoError(t, err)

	fmt.Println(track.Title)
	fmt.Println(track.Cid)
}

func TestBase62(t *testing.T) {
	b62 := base62.Encode(1000000)
	fmt.Println(b62)

	b62 = base62.Encode(10000000)
	fmt.Println(b62)

	b62 = base62.Encode(100000000)
	fmt.Println(b62)

	b62 = base62.Encode(1000000000)
	fmt.Println(b62)

	b62 = base62.Encode(10000000000)
	fmt.Println(b62)

	b62 = base62.Encode(time.Now().UnixNano())
	fmt.Println(b62)
}
