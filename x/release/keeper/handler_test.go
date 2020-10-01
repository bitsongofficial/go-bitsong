package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/release/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func Test_handleMsgReleaseCreate(t *testing.T) {
	ctx, _, keeper := SetupTestInput()
	creator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	releaseID := "ReleaseID"
	releaseID2 := "ReleaseID2"
	metadataURI := "metadata"

	h := NewHandler(keeper)
	require.NotNil(t, h)

	msg := types.NewMsgReleseCreate(releaseID, metadataURI, creator)
	require.NotNil(t, msg)

	result, err := h(ctx, msg)
	require.NoError(t, err)

	release := types.Release{}
	err = keeper.codec.UnmarshalBinaryLengthPrefixed(result.Data, &release)
	require.NoError(t, err)

	require.Equal(t, release.ReleaseID, releaseID)
	require.Equal(t, release.Creator, creator)
	require.Equal(t, release.MetadataURI, metadataURI)

	msg = types.NewMsgReleseCreate(releaseID, metadataURI, creator)
	require.NotNil(t, msg)

	result, err = h(ctx, msg)
	require.Error(t, err)

	msg = types.NewMsgReleseCreate(releaseID2, metadataURI, creator)
	require.NotNil(t, msg)

	result, err = h(ctx, msg)
	require.NoError(t, err)
}
