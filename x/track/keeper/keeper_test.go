package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/types/util"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"testing"
	"time"
)

func initConfig() (sdk.AccAddress, time.Time, types.Track) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(util.Bech32PrefixAccAddr, util.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(util.Bech32PrefixValAddr, util.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(util.Bech32PrefixConsAddr, util.Bech32PrefixConsPub)

	owner := sdk.AccAddress(crypto.AddressHash([]byte(`owner`)))
	timeZone, _ := time.LoadLocation("UTC")
	date := time.Date(2020, 1, 1, 12, 00, 00, 000, timeZone)
	trackAttrs := map[string]string{
		"title":   `The Show Must Go On`,
		"artists": `Queen`,
	}
	trackRewards := types.TrackRewards{
		Users:     10,
		Playlists: 10,
	}

	trackRightsHolders := types.RightsHolders{
		types.RightHolder{
			Address: sdk.AccAddress(crypto.AddressHash([]byte(`test`))),
			Quota:   100,
		},
	}

	trackMedia := types.TrackMedia{
		Audio: types.Content{
			Path:        "/ipfs/QM....",
			ContentType: "audio/x-mpeg",
			Duration:    5,
			Attributes:  nil,
		},
		Image: types.Content{
			Path:        "/ipfs/QM....",
			ContentType: "image/jpeg",
			Duration:    0,
			Attributes:  nil,
		},
		Video: types.Content{},
	}

	track := types.NewTrack(
		"The Show Must Go On",
		trackMedia,
		trackAttrs,
		trackRewards,
		trackRightsHolders,
		date,
		owner,
	)

	return owner, date, track
}

func TestKeeper_GetTrack(t *testing.T) {
	ctx, k := SetupTestInput()
	_, _, track := initConfig()

	tests := []struct {
		name     string
		address  crypto.Address
		track    types.Track
		expected types.Track
	}{
		{
			name:     "Get track id expected track nil",
			address:  generateTrackAddress(uint64(1)),
			track:    types.Track{},
			expected: types.Track{},
		},
		{
			name:     "Set track and get track id 1 expected track equal",
			address:  generateTrackAddress(uint64(1)),
			track:    track,
			expected: track,
		},
		{
			name:     "Set track and get track id 2 expected track equal",
			address:  generateTrackAddress(uint64(2)),
			track:    track,
			expected: track,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if !test.track.Equals(types.Track{}) {
				k.Create(ctx, test.track)
			}

			actual, ok := k.GetTrack(ctx, test.address)
			fmt.Println(actual.String())

			if test.track.Equals(types.Track{}) {
				require.False(t, ok)
			} else {
				require.Equal(t, test.address, actual.Address)
			}
		})
	}
}

func TestKeeper_CreateTrack(t *testing.T) {
	ctx, k := SetupTestInput()
	_, _, track := initConfig()

	tests := []struct {
		name           string
		existingTracks types.Tracks
		lastTrackID    uint64
		trackAddr      string
		newTrack       types.Track
	}{
		{
			name:           "Create track with lastTrackID 1",
			existingTracks: types.Tracks{},
			lastTrackID:    uint64(1),
			trackAddr:      "B0FA2953B126722264F67828AF7443144C85D867",
			newTrack:       track,
		},
		{
			name:           "Create track with lastTrackID 2",
			existingTracks: types.Tracks{},
			lastTrackID:    uint64(2),
			trackAddr:      "F1CAEDF8C538569A1884892B3144E9D566AD2607",
			newTrack:       track,
		},
		{
			name:           "Create track with lastTrackID 3",
			existingTracks: types.Tracks{},
			lastTrackID:    uint64(3),
			trackAddr:      "F9640309DE484B43A463EC66ED616E312501EF64",
			newTrack:       track,
		},
		{
			name:           "Create track with lastTrackID 4",
			existingTracks: types.Tracks{},
			lastTrackID:    uint64(4),
			trackAddr:      "007995DDDC258A86BB59A3F3C7CAEB94D559C0F7",
			newTrack:       track,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			track := test.newTrack
			addr := k.Create(ctx, track)
			generatedTrackAddr := generateTrackAddress(test.lastTrackID)
			tAddr, _ := hex.DecodeString(test.trackAddr)
			trackAddr := crypto.Address(tAddr)
			require.Equal(t, generatedTrackAddr, addr)
			require.Equal(t, generatedTrackAddr, trackAddr)

			// TODO: implement existingTracks

			expected, _ := k.GetTrack(ctx, generatedTrackAddr)
			require.Equal(t, expected.Address, generatedTrackAddr)
			require.Equal(t, expected.Title, test.newTrack.Title)
			require.True(t, expected.Media.Equals(test.newTrack.Media))
		})
	}
}

func Test_generateTrackAddress(t *testing.T) {
	initConfig()

	acc := generateTrackAddress(uint64(1))
	require.Equal(t, "B0FA2953B126722264F67828AF7443144C85D867", acc.String())

	acc = generateTrackAddress(uint64(2))
	require.Equal(t, "F1CAEDF8C538569A1884892B3144E9D566AD2607", acc.String())

	acc = generateTrackAddress(uint64(3))
	require.Equal(t, "F9640309DE484B43A463EC66ED616E312501EF64", acc.String())

	acc = generateTrackAddress(uint64(4))
	require.Equal(t, "007995DDDC258A86BB59A3F3C7CAEB94D559C0F7", acc.String())

	acc = generateTrackAddress(uint64(400000000000))
	require.Equal(t, "CCE1469BDEFFA9CC460B3ED46007A2DCAE5515FD", acc.String())
}

func TestKeeper_autoIncrementID(t *testing.T) {
	ctx, k := SetupTestInput()

	for lastId := uint64(1); lastId <= uint64(10); lastId++ {
		k.autoIncrementID(ctx)
		store := ctx.KVStore(k.storeKey)
		bz := store.Get(types.KeyLastTrackID)
		bzLastId := sdk.Uint64ToBigEndian(lastId + 1)
		require.Equal(t, bz, bzLastId)
	}
}
