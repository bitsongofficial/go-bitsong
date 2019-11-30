package keeper

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleVerifyArtistProposal is a handler for executing a passed artist verify proposal
func HandleVerifyArtistProposal(ctx sdk.Context, k Keeper, p types.ArtistVerifyProposal) sdk.Error {
	// TODO:
	// Handle Verify Artist

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("TODO: HandleVerifyArtistProposal"))

	return nil
}

func NewArtistVerifyProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Error {
		switch c := content.(type) {
		case types.ArtistVerifyProposal:
			return HandleVerifyArtistProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized artist proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}
