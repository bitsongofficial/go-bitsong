package artist

import (
	"github.com/stretchr/testify/require"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestGetSetArtist(t *testing.T) {
	input := getMockApp(t, 1, GenesisState{}, nil)
	SortAddresses(input.addrs)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	artist := testArtist()
	artist1, err := input.keeper.CreateArtist(ctx, artist, input.addrs[0])
	require.NoError(t, err)
	artistID := artist1.ArtistID
	input.keeper.SetArtist(ctx, artist1)

	gotArtist, ok := input.keeper.GetArtist(ctx, artistID)
	require.True(t, ok)
	require.True(t, ArtistEqual(artist1, gotArtist))
}

func TestIncrementArtistNumber(t *testing.T) {
	input := getMockApp(t, 1, GenesisState{}, nil)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	artist := testArtist()
	input.keeper.CreateArtist(ctx, artist, input.addrs[0])
	input.keeper.CreateArtist(ctx, artist, input.addrs[0])
	input.keeper.CreateArtist(ctx, artist, input.addrs[0])
	input.keeper.CreateArtist(ctx, artist, input.addrs[0])
	input.keeper.CreateArtist(ctx, artist, input.addrs[0])
	artist6, err := input.keeper.CreateArtist(ctx, artist, input.addrs[0])
	require.NoError(t, err)

	require.Equal(t, uint64(6), artist6.ArtistID)
}

func TestCreateArtist(t *testing.T) {
	input := getMockApp(t, 1, GenesisState{}, nil)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})
	input.mApp.InitChainer(ctx, abci.RequestInitChain{})

	// TODO: implement new tests

	/*testCases := []struct {
		meta        types.Meta
		expectedErr sdk.Error
	}{
		{validArtist{}, nil},
		// Keeper does not check the validity of name, no error
		{invalidArtistName1{}, nil},
		{invalidArtistName2{}, nil},
		// TODO: error only when invalid route
		// {invalidArtistRoute{}, types.ErrNoArtistHandlerExists(types.DefaultCodespace, invalidArtistRoute{})},
		// Keeper does not call ValidateBasic, msg.ValidateBasic does
		{invalidArtistValidation{}, nil},
	}

	for _, tc := range testCases {
		_, err := input.keeper.CreateArtist(ctx, tc.meta, input.addrs[0])
		require.Equal(t, tc.expectedErr, err, "unexpected type of error: %s", err)
	}*/
}
