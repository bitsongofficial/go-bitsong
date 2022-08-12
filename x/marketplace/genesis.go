package marketplace

import (
	"github.com/bitsongofficial/go-bitsong/x/marketplace/keeper"
	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params: types.DefaultParams(),
	}
}

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	// initialize params
	k.SetParamSet(ctx, data.Params)

	for _, auction := range data.Auctions {
		k.SetAuction(ctx, auction)
	}

	for _, bid := range data.Bids {
		k.SetBid(ctx, bid)
	}

	for _, bidder := range data.BidderMetadata {
		k.SetBidderMetadata(ctx, bidder)
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:         k.GetParamSet(ctx),
		Auctions:       k.GetAllAuctions(ctx),
		Bids:           k.GetAllBids(ctx),
		BidderMetadata: k.GetAllBidderMetadata(ctx),
	}
}
