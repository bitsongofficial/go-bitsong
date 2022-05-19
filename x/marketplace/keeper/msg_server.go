package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateAuction(goCtx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	auctionId, err := m.Keeper.CreateAuction(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateAuctionResponse{Id: auctionId}, nil
}

func (m msgServer) StartAuction(goCtx context.Context, msg *types.MsgStartAuction) (*types.MsgStartAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.StartAuction(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgStartAuctionResponse{}, nil
}

func (m msgServer) SetAuctionAuthority(goCtx context.Context, msg *types.MsgSetAuctionAuthority) (*types.MsgSetAuctionAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.SetAuctionAuthority(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSetAuctionAuthorityResponse{}, nil
}

func (m msgServer) EndAuction(goCtx context.Context, msg *types.MsgEndAuction) (*types.MsgEndAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// 1.  Check executor is a correct authority of the auction
	// 2.  Check auction is not already ended
	// 3.  Set auction end time
	// 4.  Set auction status as ended
	// 5.  Check auction has winning bid
	// 6.  If winning bid does not exists, send nft and metadata ownership to auction authority

	// 7.  Emit event for auction end
	// ctx.EventManager().EmitTypedEvent(&types.EventEndAuction{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	return &types.MsgEndAuctionResponse{}, nil
}

func (m msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// 1. Verify auction is `Started` status
	// 2. Verify bid is valid for the auction (check `bid_denom`, `tick_size` and `last_bid_amount`)
	// 3. Add new bid for the auction on the storage
	// 4. Confirm payer does have enough token to pay the bid
	// 5. Transfer amount of token to bid account
	// 6. Serialize new auction state with new bid
	// 7. Update or create bidder metadata
	// 8. If the amount exceeds `instant_sale_price`, end the auction

	// 9. Emit event for placing bid
	// ctx.EventManager().EmitTypedEvent(&types.EventPlaceBid{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	return &types.MsgPlaceBidResponse{}, nil
}

func (m msgServer) CancelBid(goCtx context.Context, msg *types.MsgCancelBid) (*types.MsgCancelBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// 1. Load the auction and verify this bid is valid.
	// 2. Refuse to cancel if the auction ended and this person is a winning account.
	// 3. Remove bid from the storage
	// 4. Transfer tokens back to the bidder
	// 5. Update bidder Metadata
	// 6. Update auction with remaining bids

	// 7. Emit event for placing bid
	// ctx.EventManager().EmitTypedEvent(&types.EventCancelBid{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	return &types.MsgCancelBidResponse{}, nil
}

func (m msgServer) ClaimBid(goCtx context.Context, msg *types.MsgClaimBid) (*types.MsgClaimBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// 1. Load the auction and verify this bid is valid.
	// 2. Ensure the sender is winner bidder
	// 3. Send bid amount to auction authority
	// 4. If `primary_sale_happened` is true, process royalties from NFT's `seller_fee_basis_points` field to creators
	// 5. Set `primary_sale_happened` as true if it was not set already
	// 6. Transfer ownership of NFT to bidder
	// 7. If auction type is for transferring metadata ownership as well, transfer metadata ownership as well

	// 8. Emit event for claiming bid
	// ctx.EventManager().EmitTypedEvent(&types.EventClaimBid{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	return &types.MsgClaimBidResponse{}, nil
}
