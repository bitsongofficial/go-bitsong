package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

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

// HandleVerifyArtistProposal is a handler for executing a passed artist verify proposal
func HandleVerifyArtistProposal(ctx sdk.Context, k Keeper, avp types.ArtistVerifyProposal) sdk.Error {
	logger := k.Logger(ctx)
	// TODO:
	// complete...
	logger.Info(fmt.Sprintf("Artist ID: %d", avp.ArtistID))

	fmt.Println()
	fmt.Println()
	fmt.Printf("HandlerArtist ID: %d", avp.ArtistID)
	fmt.Println()
	fmt.Println()

	// Get artist
	artist, ok := k.GetArtist(ctx, avp.ArtistID)
	if !ok {
		return types.ErrUnknownArtist(k.codespace, "unknown artist")
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("Artist: %s", artist.String())
	fmt.Println()
	fmt.Println()

	// Check if status is nil
	if artist.Status != types.StatusNil {
		return types.ErrInvalidArtistStatus(k.codespace, "artist status must be nil")
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("Artist Status: %s", artist.Status.String())
	fmt.Println()
	fmt.Println()

	// Set status verified
	k.SetArtistStatus(ctx, avp.ArtistID, types.StatusVerified)

	//logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("TODO: HandleVerifyArtistProposal"))

	return nil
}
