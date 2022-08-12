package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Auctions(c context.Context, req *types.QueryAuctionsRequest) (*types.QueryAuctionsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	allAuctions := k.GetAllAuctions(ctx)
	auctions := []types.Auction{}
	for _, auction := range allAuctions {
		if req.State != types.AuctionState_Empty && auction.State != req.State {
			continue
		}
		if req.Authority != "" && auction.Authority != req.Authority {
			continue
		}
		auctions = append(auctions, auction)
	}

	return &types.QueryAuctionsResponse{
		Auctions: auctions,
	}, nil
}

func (k Keeper) Auction(c context.Context, req *types.QueryAuctionRequest) (*types.QueryAuctionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	auction, err := k.GetAuctionById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryAuctionResponse{
		Auction: auction,
	}, nil
}

func (k Keeper) BidsByAuction(c context.Context, req *types.QueryBidsByAuctionRequest) (*types.QueryBidsByAuctionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bids := k.GetBidsByAuction(ctx, req.Id)
	return &types.QueryBidsByAuctionResponse{
		Bids: bids,
	}, nil
}

func (k Keeper) BidsByBidder(c context.Context, req *types.QueryBidsByBidderRequest) (*types.QueryBidsByBidderResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bidder, err := sdk.AccAddressFromBech32(req.Bidder)
	if err != nil {
		return nil, err
	}
	bids := k.GetBidsByBidder(ctx, bidder)
	return &types.QueryBidsByBidderResponse{
		Bids: bids,
	}, nil
}

func (k Keeper) BidderMetadata(c context.Context, req *types.QueryBidderMetadataRequest) (*types.QueryBidderMetadataResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bidder, err := sdk.AccAddressFromBech32(req.Bidder)
	if err != nil {
		return nil, err
	}
	bidderdata, err := k.GetBidderMetadata(ctx, bidder)
	if err != nil {
		return nil, err
	}
	return &types.QueryBidderMetadataResponse{
		BidderMetadata: bidderdata,
	}, nil
}
