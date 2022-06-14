package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagAuctionId        = "auction-id"
	FlagNewAuthority     = "new-authority"
	FlagNftId            = "nft-id"
	FlagPrizeType        = "prize-type"
	FlagBidDenom         = "bid-denom"
	FlagDuration         = "duration"
	FlagPriceFloor       = "price-floor"
	FlagInstantSalePrice = "instant-sale-price"
	FlagTickSize         = "tick-size"
	FlagAmount           = "amount"

	FlagAuctionState = "state"
	FlagAuthority    = "authority"
	FlagEditionLimit = "edition-limit"
)

func FlagCreateAuction() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagNftId, 0, "Id of the nft to put on auction")
	fs.String(FlagPrizeType, "", "Type of prize, NFT_ONLY_TRANSFER | FULL_RIGHTS_TRANSFER | LIMITED_EDITION_PRINTS | OPEN_EDITION_PRINTS")
	fs.String(FlagBidDenom, "", "Denom to be used for bidding on the auction")
	fs.Duration(FlagDuration, 0, "Duration of the auction")
	fs.Uint64(FlagPriceFloor, 0, "Minimum price of the nft.")
	fs.Uint64(FlagInstantSalePrice, 0, "Instant sale price of the nft.")
	fs.Uint64(FlagTickSize, 0, "Tick size to be increased at least for new bid.")
	fs.Uint64(FlagEditionLimit, 0, "Edition limit for printing: valid for LIMITED_EDITION_PRINTS auction")

	return fs
}

func FlagSetAuctionAuthority() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagAuctionId, 0, "Id of the auction")
	fs.String(FlagNewAuthority, "", "new authority of the auction")

	return fs
}

func FlagStartAuction() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagAuctionId, 0, "Id of the auction")

	return fs
}

func FlagEndAuction() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagAuctionId, 0, "Id of the auction")

	return fs
}

func FlagPlaceBid() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagAuctionId, 0, "Id of the auction")
	fs.String(FlagAmount, "", "Amount to bid")

	return fs
}

func FlagCancelBid() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagAuctionId, 0, "Id of the auction")

	return fs
}

func FlagClaimBid() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagAuctionId, 0, "Id of the auction")

	return fs
}

func FlagQueryAuctions() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAuctionState, "", "State of the auction, EMPTY | CREATED | STARTED | ENDED")
	fs.String(FlagAuthority, "", "Authority of the auction")

	return fs
}
