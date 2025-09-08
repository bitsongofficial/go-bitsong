package types_test

import (
	"testing"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/stretchr/testify/require"
)

func TestSplitNftLengthPrefixedKey(t *testing.T) {
	denom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"
	tokenId := "1"

	keyBz := append(append([]byte(denom), 0), []byte(tokenId)...)

	denomBz, tokenIdBz := types.MustSplitNftLengthPrefixedKey(keyBz)
	require.Equal(t, denom, string(denomBz))
	require.Equal(t, tokenId, string(tokenIdBz))
}
