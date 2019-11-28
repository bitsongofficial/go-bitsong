package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

var addr = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

func TestArtistKeys(t *testing.T) {
	// key artist
	key := ArtistKey(1)
	artistID := SplitArtistKey(key)
	require.Equal(t, int(artistID), 1)
}
