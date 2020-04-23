package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/types"
	trackTypes "github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"testing"
	"time"
)

func Test_queryTrack(t *testing.T) {
	fakeTrackAddr := crypto.AddressHash([]byte("test"))

	owner := sdk.AccAddress(crypto.AddressHash([]byte(`owner`)))
	timeZone, _ := time.LoadLocation("UTC")
	date := time.Date(2020, 1, 1, 12, 00, 00, 000, timeZone)
	testDate := time.Date(0001, 1, 1, 0, 00, 00, 000, timeZone)
	trackAttrs := map[string]string{
		"title":   `The Show Must Go On`,
		"artists": `Queen`,
	}
	trackRewards := trackTypes.TrackRewards{
		Users:     10,
		Playlists: 10,
	}

	trackRightsHolders := trackTypes.RightsHolders{
		trackTypes.RightHolder{
			Address: sdk.AccAddress(crypto.AddressHash([]byte(`test`))),
			Quota:   100,
		},
	}

	trackMedia := trackTypes.TrackMedia{
		Audio: trackTypes.Content{
			Path:        "/ipfs/QM....",
			ContentType: "audio/x-mpeg",
			Duration:    5,
			Attributes:  nil,
		},
		Image: trackTypes.Content{
			Path:        "/ipfs/QM....",
			ContentType: "image/jpeg",
			Duration:    0,
			Attributes:  nil,
		},
		Video: trackTypes.Content{},
	}

	tests := []struct {
		name         string
		path         []string
		storedTracks trackTypes.Tracks
		expResult    trackTypes.Track
		expError     error
	}{
		{
			name:     "Invalid query endpoint",
			path:     []string{"invalid", ""},
			expError: fmt.Errorf("unknown track query endpoint: unknown request"),
		},
		{
			name:     "Invalid Track Addr returns error",
			path:     []string{trackTypes.QueryTrack, ""},
			expError: sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown track addr "),
		},
		{
			name:     "Track not found returns error",
			path:     []string{trackTypes.QueryTrack, fakeTrackAddr.String()},
			expError: sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("track addr %s not found", fakeTrackAddr.String())),
		},
		{
			name: "Track is returned properly",
			storedTracks: trackTypes.Tracks{
				trackTypes.NewTrack(
					"The Show Must Go On",
					trackMedia,
					trackAttrs,
					trackRewards,
					trackRightsHolders,
					date,
					owner,
				),
			},
			path: []string{trackTypes.QueryTrack, "B0FA2953B126722264F67828AF7443144C85D867"},
			expResult: trackTypes.Track{
				Title:         "The Show Must Go On",
				Address:       generateTrackAddress(uint64(1)),
				Attributes:    trackAttrs,
				Media:         trackMedia,
				Rewards:       trackRewards,
				RightsHolders: trackRightsHolders,
				Totals: trackTypes.TrackTotals{
					Streams:  0,
					Rewards:  sdk.NewCoin(types.BondDenom, sdk.ZeroInt()),
					Accounts: 0,
				},
				SubmitTime: testDate,
				Owner:      owner,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			for _, t := range test.storedTracks {
				k.Create(ctx, t)
			}

			querier := NewQuerier(k)
			result, err := querier(ctx, test.path, abci.RequestQuery{})

			if result != nil {
				require.Nil(t, err)
				expectedIndented, err := codec.MarshalJSONIndent(k.cdc, &test.expResult)
				require.NoError(t, err)
				require.Equal(t, string(expectedIndented), string(result))
			}

			if result == nil {
				require.NotNil(t, err)
				require.Equal(t, test.expError.Error(), err.Error())
				require.Nil(t, result)
			}
		})
	}
}
