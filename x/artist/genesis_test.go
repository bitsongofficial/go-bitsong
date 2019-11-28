package artist

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestEqualArtistID(t *testing.T) {
	state1 := GenesisState{}
	state2 := GenesisState{}
	require.Equal(t, state1, state2)

	// Artists
	state1.StartingArtistID = 1
	require.NotEqual(t, state1, state2)
	require.False(t, state1.Equal(state2))

	state2.StartingArtistID = 1
	require.Equal(t, state1, state2)
	require.True(t, state1.Equal(state2))
}

func TestEqualArtists(t *testing.T) {
	// Generate mock app and keepers
	input := getMockApp(t, 1, GenesisState{}, nil)
	SortAddresses(input.addrs)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	// Create two artists
	artist := testArtist()
	images := testImages()
	artist1, err := input.keeper.CreateArtist(ctx, artist, images, input.addrs[0])
	require.NoError(t, err)
	artist2, err := input.keeper.CreateArtist(ctx, artist, images, input.addrs[0])
	require.NoError(t, err)

	// They are similar but their IDs should be different
	require.NotEqual(t, artist1, artist2)
	require.False(t, ArtistEqual(artist1, artist2))

	// Now create two genesis blocks
	state1 := GenesisState{Artists: []types.Artist{artist1}}
	state2 := GenesisState{Artists: []types.Artist{artist2}}
	require.NotEqual(t, state1, state2)
	require.False(t, state1.Equal(state2))

	// Now make artists identical by setting both IDs to 9
	artist1.ArtistID = 9
	artist2.ArtistID = 9
	require.Equal(t, artist1, artist2)
	require.True(t, ArtistEqual(artist1, artist2))

	// Reassign artists into state
	state1.Artists[0] = artist1
	state2.Artists[0] = artist2

	// State should be identical now..
	require.Equal(t, state1, state2)
	require.True(t, state1.Equal(state2))
}

func TestImportExportArtists(t *testing.T) {
	// Generate mock app and keepers
	input := getMockApp(t, 1, GenesisState{}, nil)
	SortAddresses(input.addrs)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	// Create two artists
	// TODO: put the second into the status Verified
	artist := testArtist()
	images := testImages()
	artist1, err := input.keeper.CreateArtist(ctx, artist, images, input.addrs[0])
	require.NoError(t, err)
	artistID1 := artist1.ArtistID

	artist2, err := input.keeper.CreateArtist(ctx, artist, images, input.addrs[0])
	require.NoError(t, err)
	artistID2 := artist2.ArtistID

	artist1, ok := input.keeper.GetArtist(ctx, artistID1)
	require.True(t, ok)
	artist2, ok = input.keeper.GetArtist(ctx, artistID2)
	require.True(t, ok)
	require.True(t, artist1.Status == types.StatusNil)
	require.True(t, artist2.Status == types.StatusNil)

	genAccs := input.mApp.AccountKeeper.GetAllAccounts(ctx)

	// Export the state and import it into a new Mock App
	genState := ExportGenesis(ctx, input.keeper)
	input2 := getMockApp(t, 1, genState, genAccs)

	header = abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input2.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx2 := input2.mApp.BaseApp.NewContext(false, abci.Header{})

	artist1, ok = input2.keeper.GetArtist(ctx2, artistID1)
	require.True(t, ok)
	artist2, ok = input2.keeper.GetArtist(ctx2, artistID2)
	require.True(t, ok)
	require.True(t, artist1.Status == types.StatusNil)
	require.True(t, artist2.Status == types.StatusNil)
}
