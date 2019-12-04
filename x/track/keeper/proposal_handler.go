package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

func NewTrackVerifyProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Error {
		switch c := content.(type) {
		case types.TrackVerifyProposal:
			return HandleVerifyTrackProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized artist proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

// HandleVerifyTrackProposal is a handler for executing a passed artist verify proposal
func HandleVerifyTrackProposal(ctx sdk.Context, k Keeper, avp types.TrackVerifyProposal) sdk.Error {
	logger := k.Logger(ctx)
	// TODO:
	// complete...
	logger.Info(fmt.Sprintf("Track ID: %d", avp.TrackID))

	// Get track
	track, ok := k.GetTrack(ctx, avp.TrackID)
	if !ok {
		return types.ErrUnknownTrack(k.codespace, "unknown track")
	}

	// Check if status is nil
	if track.Status != types.StatusNil {
		return types.ErrInvalidTrackStatus(k.codespace, "track status must be nil")
	}

	// Set status verified
	k.SetTrackStatus(ctx, avp.TrackID, types.StatusVerified)

	//logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("TODO: HandleVerifyTrackProposal"))

	return nil
}
