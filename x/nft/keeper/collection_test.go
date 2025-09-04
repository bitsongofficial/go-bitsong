package keeper

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_createCollectionDenom(t *testing.T) {
	creator := sdk.AccAddress(tmhash.SumTruncated([]byte("creator")))
	symbol := "MYNFT"

	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"
	k := Keeper{}

	denom := k.createCollectionDenom(creator, symbol)
	if denom != expectedDenom {
		t.Errorf("expected %s, got %s", expectedDenom, denom)
	}
}

func TestSplitNftLengthPrefixedKey(t *testing.T) {
	denom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"
	tokenId := "1"

	keyBz := append(append([]byte(denom), 0), []byte(tokenId)...)

	denomBz, tokenIdBz := MustSplitNftLengthPrefixedKey(keyBz)
	require.Equal(t, denom, string(denomBz))
	require.Equal(t, tokenId, string(tokenIdBz))
}
