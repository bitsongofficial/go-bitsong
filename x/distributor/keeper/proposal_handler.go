package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
)

func NewDistributorVerifyProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Error {
		switch c := content.(type) {
		case types.DistributorVerifyProposal:
			return HandleVerifyDistributorProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized distributor proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

func HandleVerifyDistributorProposal(ctx sdk.Context, k Keeper, avp types.DistributorVerifyProposal) sdk.Error {
	logger := k.Logger(ctx)
	// TODO:
	// complete...
	logger.Info(fmt.Sprintf("Distributor address: %s", avp.Address.String()))

	// Get distributor
	distributor, ok := k.GetDistributor(ctx, avp.Address)
	if !ok {
		return types.ErrInvalidDistributor(k.codespace, "unknown distributor")
	}

	// Check if status is nil
	if distributor.Status != types.StatusNil {
		return types.ErrInvalidDistributor(k.codespace, "distributor status must be nil")
	}

	// Set status verified
	k.SetStatus(ctx, avp.Address, types.StatusVerified)

	//logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("TODO: HandleVerifyDistributorProposal"))

	return nil
}
