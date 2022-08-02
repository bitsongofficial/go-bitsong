package cli

import (
	flag "github.com/spf13/pflag"

	nftcli "github.com/bitsongofficial/go-bitsong/x/nft/client/cli"
)

const (
	FlagPrice            = "price"
	FlagTreasury         = "treasury"
	FlagDenom            = "denom"
	FlagGoLiveDate       = "go-live-date"
	FlagEndSettingsType  = "end-settings-type"
	FlagEndSettingsValue = "end-settings-value"
	FlagMetadataBaseUrl  = "metadata-baseurl"
)

func FlagCreateCandyMachine() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(FlagPrice, 0, "price of candymachine nft")
	fs.String(FlagTreasury, "", "treasury to receive payment for nft minting")
	fs.String(FlagDenom, "", "denom to be spent for candymachine nft minting price")
	fs.String(FlagEndSettingsType, "", "end settings type")
	fs.Uint64(FlagEndSettingsValue, 0, "end settings value")
	fs.String(FlagMetadataBaseUrl, "", "metadata base url")
	fs.String(FlagGoLiveDate, "", "go live date")

	fs.Uint64(nftcli.FlagCollectionId, 0, "collection id for the nft")
	fs.String(nftcli.FlagCreators, "", "Creators of nft")
	fs.String(nftcli.FlagCreatorShares, "", "Shares between creators for seller fee")
	fs.Bool(nftcli.FlagMutable, false, "mutability of the nft")
	fs.Uint32(nftcli.FlagSellerFeeBasisPoints, 0, "Seller fee basis points of the nft")

	return fs
}
