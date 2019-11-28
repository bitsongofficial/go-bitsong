package artist

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const custom = "custom"

func getQueriedArtist(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, artistID uint64) types.Artist {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryArtist}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryArtistParams(artistID)),
	}

	bz, err := querier(ctx, []string{types.QueryArtist}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var artist types.Artist
	err2 := cdc.UnmarshalJSON(bz, artist)
	require.Nil(t, err2)
	return artist
}

func getQueriedArtists(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, owner sdk.AccAddress, status types.ArtistStatus, limit uint64) []types.Artist {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryArtists}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryArtistsParams(owner, status, limit)),
	}

	bz, err := querier(ctx, []string{types.QueryArtists}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var artists types.Artists
	err2 := cdc.UnmarshalJSON(bz, &artists)
	require.Nil(t, err2)
	return artists
}

func TestQueries(t *testing.T) {
	cdc := codec.New()
	input := getMockApp(t, 1000, GenesisState{}, nil)
	querier := NewQuerier(input.keeper)
	handler := NewHandler(input.keeper)

	types.RegisterCodec(cdc)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.NewContext(false, abci.Header{})

	// input.addrs[0] create artist #1 and #2
	res := handler(ctx, types.NewMsgCreateArtist(testArtist(), input.addrs[0]))
	var artistID1 uint64
	require.True(t, res.IsOK())
	cdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &artistID1)

	res = handler(ctx, types.NewMsgCreateArtist(testArtist(), input.addrs[0]))
	var artistID2 uint64
	require.True(t, res.IsOK())
	cdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &artistID2)

	// input.addrs[1] create artist #3
	res = handler(ctx, types.NewMsgCreateArtist(testArtist(), input.addrs[1]))
	var artistID3 uint64
	require.True(t, res.IsOK())
	cdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &artistID3)

	// Test proposals queries with filters

	// Test query all proposals and status nil
	artists := getQueriedArtists(t, ctx, cdc, querier, nil, types.StatusNil, 0)
	require.Equal(t, artistID1, (artists[0]).ArtistID)
	require.Equal(t, artistID2, (artists[1]).ArtistID)
	require.Equal(t, artistID3, (artists[2]).ArtistID)

	// Test query artists added by input.addrs[0] and status nil
	artists = getQueriedArtists(t, ctx, cdc, querier, input.addrs[0], types.StatusNil, 0)
	require.Equal(t, artistID1, (artists[0]).ArtistID)
	require.Equal(t, artistID2, (artists[1]).ArtistID)
	require.NotEqual(t, artistID1, (artists[2]).ArtistID)
	require.NotEqual(t, artistID2, (artists[2]).ArtistID)

	// Test query artists added by input.addrs[1] and status nil
	artists = getQueriedArtists(t, ctx, cdc, querier, input.addrs[1], types.StatusNil, 0)
	require.NotEqual(t, artistID3, (artists[0]).ArtistID)
	require.NotEqual(t, artistID3, (artists[1]).ArtistID)
	require.Equal(t, artistID3, (artists[2]).ArtistID)
}
