package keeper

import (
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestKeeper_CreateArtist(t *testing.T) {
	ctx, _, keeper := SetupTestInput()
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	id := btsg.ID("dasdas-dasdsa")
	name := "Freedy Mercury"
	mockArtist := types.NewArtist(id, name, nil, nil, "", addr)

	artist, err := keeper.CreateArtist(ctx, mockArtist)
	require.NoError(t, err)
	require.Equal(t, artist.ID, mockArtist.ID)
	require.Equal(t, artist.Name, mockArtist.Name)
	require.Equal(t, artist.URLs, mockArtist.URLs)
	require.Equal(t, artist.MetadataURI, mockArtist.MetadataURI)
	require.Equal(t, artist.Creator, mockArtist.Creator)

	artist, err = keeper.CreateArtist(ctx, mockArtist)
	require.Error(t, err)
}
