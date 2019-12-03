package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

func NewAlbumVerifyProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Error {
		switch c := content.(type) {
		case types.AlbumVerifyProposal:
			return HandleAlbumVerifyProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized album proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

// HandleAlbumVerifyProposal is a handler for executing a passed album verify proposal
func HandleAlbumVerifyProposal(ctx sdk.Context, k Keeper, avp types.AlbumVerifyProposal) sdk.Error {
	logger := k.Logger(ctx)
	// TODO:
	// complete...
	logger.Info(fmt.Sprintf("Album ID: %d", avp.AlbumID))

	// Get album
	album, ok := k.GetAlbum(ctx, avp.AlbumID)
	if !ok {
		return types.ErrUnknownAlbum(k.codespace, "unknown album")
	}

	// Check if status is nil
	if album.Status != types.StatusNil {
		return types.ErrInvalidAlbumStatus(k.codespace, "album status must be nil")
	}

	// Set status verified
	k.SetAlbumStatus(ctx, avp.AlbumID, types.StatusVerified)

	//logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("TODO: HandleVerifyArtistProposal"))

	return nil
}
