package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// module sentinel errors
var (
	ErrAuctionDoesNotExist         = sdkerrors.Register(ModuleName, 2, "auction does not exist")
	ErrNotAuctionAuthority         = sdkerrors.Register(ModuleName, 3, "not an auction authority")
	ErrAuctionAlreadyStarted       = sdkerrors.Register(ModuleName, 4, "auction already started")
	ErrAuctionAlreadyEnded         = sdkerrors.Register(ModuleName, 5, "auction already ended")
	ErrAuctionNotStarted           = sdkerrors.Register(ModuleName, 6, "auction not started")
	ErrInvalidBidDenom             = sdkerrors.Register(ModuleName, 7, "bid denom is invalid")
	ErrInvalidBidAmount            = sdkerrors.Register(ModuleName, 8, "bid amount is invalid")
	ErrBidDoesNotExists            = sdkerrors.Register(ModuleName, 9, "bid does not exists")
	ErrBidAlreadyExists            = sdkerrors.Register(ModuleName, 10, "bid already exists")
	ErrBidderMetadataDoesNotExists = sdkerrors.Register(ModuleName, 11, "bidder metadata does not exists")
	ErrCanNotCancelWinningBid      = sdkerrors.Register(ModuleName, 12, "cannot cancel winning bid")
	ErrNotWinningBid               = sdkerrors.Register(ModuleName, 13, "not a winning winning bid")
	ErrInvalidNftId                = sdkerrors.Register(ModuleName, 14, "invalid nft id")
	ErrInvalidInstantSalePrice     = sdkerrors.Register(ModuleName, 15, "invalid instant sale price")
	ErrInvalidDuration             = sdkerrors.Register(ModuleName, 16, "invalid duration")
	ErrInvalidPrizeType            = sdkerrors.Register(ModuleName, 17, "invalid prize type")
	ErrInvalidAuctionId            = sdkerrors.Register(ModuleName, 18, "invalid auction id")
)
