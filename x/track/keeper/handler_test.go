package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

var mockMsgCreate = types.NewMsgCreate(
	mockTitle,
	mockEmptyAttribute,
	mockTrackMediaNoVideo,
	mockRewards,
	mockRightsHoldersSingle,
	mockOwner,
)

// TODO: implement more tests
func Test_handleMsgCreate(t *testing.T) {
	tests := []struct {
		name         string
		storedTracks types.Tracks
		msg          types.MsgCreate
		expTrack     types.Track
		expTrackAddr crypto.Address
		expError     error
	}{
		{
			name:         "Track is stored properly",
			msg:          mockMsgCreate,
			expTrack:     mockTrack,
			expTrackAddr: mockTrackAddr1,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			for _, t := range test.storedTracks {
				k.Create(ctx, t)
			}

			handler := NewHandler(k)
			res, err := handler(ctx, test.msg)

			// Valid response
			if res != nil {
				// Check the post
				stored, ok := k.GetTrack(ctx, test.expTrackAddr)
				require.True(t, ok)
				require.Equal(t, stored.Address, test.expTrackAddr)

				// Check the data
				require.Equal(t, k.cdc.MustMarshalBinaryLengthPrefixed(test.expTrackAddr), res.Data)

				// Check the events
				event := sdk.NewEvent(
					types.EventTypeTrackCreated,
					sdk.NewAttribute(types.AttributeKeyTrackAddr, test.expTrackAddr.String()),
				)
				require.Len(t, ctx.EventManager().Events(), 1)
				require.Contains(t, ctx.EventManager().Events(), event)
			}

			// Invalid response
			if res == nil {
				require.NotNil(t, err)
				require.Equal(t, test.expError.Error(), err.Error())
			}
		})
	}
}
