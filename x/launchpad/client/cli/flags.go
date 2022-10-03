package cli

import (
	flag "github.com/spf13/pflag"

	nftcli "github.com/bitsongofficial/go-bitsong/x/nft/client/cli"
)

const (
	FlagPrice           = "price"
	FlagTreasury        = "treasury"
	FlagDenom           = "denom"
	FlagGoLiveDate      = "go-live-date"
	FlagEndTimestamp    = "end-timestamp"
	FlagMaxMint         = "max-mint"
	FlagMetadataBaseUrl = "metadata-baseurl"
	FlagShuffle         = "shuffle"
)

func FlagCreateLaunchPad() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagPrice, 0, "price of launchpad nft")
	fs.String(FlagTreasury, "", "treasury to receive payment for nft minting")
	fs.String(FlagDenom, "", "denom to be spent for launchpad nft minting price")
	fs.Uint64(FlagEndTimestamp, 0, "end timestamp")
	fs.Uint64(FlagMaxMint, 0, "max mint")
	fs.String(FlagMetadataBaseUrl, "", "metadata base url")
	fs.Uint64(FlagGoLiveDate, 0, "go live date")
	fs.Bool(FlagShuffle, true, "Flag if shuffle metadata ids or not")

	fs.Uint64(nftcli.FlagCollectionId, 0, "collection id for the nft")
	fs.String(nftcli.FlagCreators, "", "Creators of nft")
	fs.String(nftcli.FlagCreatorShares, "", "Shares between creators for seller fee")
	fs.Bool(nftcli.FlagMutable, false, "mutability of the nft")
	fs.Uint32(nftcli.FlagSellerFeeBasisPoints, 0, "Seller fee basis points of the nft")

	return fs
}
