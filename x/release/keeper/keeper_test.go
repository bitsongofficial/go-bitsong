package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestKeeper_CreateRelease(t *testing.T) {
	ctx, _, keeper := SetupTestInput()
	creator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	releaseID := "releaseID"
	releaseID2 := "releaseID2"
	metadataURI := "metadata"

	release, err := keeper.CreateRelease(ctx, creator, releaseID, metadataURI)
	require.NoError(t, err)
	require.Equal(t, release.Creator, creator)
	require.Equal(t, release.ReleaseID, releaseID)
	require.Equal(t, release.MetadataURI, metadataURI)

	release, err = keeper.CreateRelease(ctx, creator, releaseID, metadataURI)
	require.Error(t, err)

	release, err = keeper.CreateRelease(ctx, creator, releaseID2, metadataURI)
	require.NoError(t, err)
}
