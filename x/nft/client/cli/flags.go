package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName                 = "name"
	FlagSymbol               = "symbol"
	FlagUri                  = "uri"
	FlagSellerFeeBasisPoints = "seller-fee-basis-points"
	FlagCreators             = "creators"
	FlagCreatorShares        = "creator-shares"
	FlagMutable              = "mutable"
	FlagUpdateAuthority      = "update-authority"
)

func FlagCreateNFT() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagName, "", "Name of the nft")
	fs.String(FlagSymbol, "", "Symbol of the nft")
	fs.String(FlagUri, "", "Uri of the nft")
	fs.Uint64(FlagSellerFeeBasisPoints, 0, "Seller fee basis points of the nft")
	fs.String(FlagCreators, "", "Creators of nft")
	fs.String(FlagCreatorShares, "", "Shares between creators for seller fee")
	fs.Bool(FlagMutable, false, "mutability of the nft")
	fs.String(FlagUpdateAuthority, "", "update authority of the nft")

	return fs
}
