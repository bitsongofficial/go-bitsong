package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSymbol      = "symbol"
	FlagName        = "name"
	FlagMaxSupply   = "max-supply"
	FlagMintable    = "mintable"
	FlagDescription = "description"
	FlagRecipient   = "recipient"
	FlagAmount      = "amount"
	FlagIssueFee    = "issue-fee"
	FlagURI         = "uri"
)

var (
	FsIssueFanToken         = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditFanToken          = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferFanTokenOwner = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintFanToken          = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueFanToken.String(FlagSymbol, "", "The fantoken symbol. Once created, it cannot be modified")
	FsIssueFanToken.String(FlagName, "", "The fantoken name, e.g. Bitsong Network")
	FsIssueFanToken.String(FlagMaxSupply, "", "The maximum supply of the fantoken")
	FsIssueFanToken.Bool(FlagMintable, false, "Whether the fantoken can be minted, default to false")
	FsIssueFanToken.String(FlagDescription, "", "The fantoken description")
	FsIssueFanToken.String(FlagURI, "", "The fantoken uri")
	FsIssueFanToken.String(FlagIssueFee, "", "The fan fantoken issue fee")

	FsEditFanToken.String(FlagName, "[do-not-modify]", "The fantoken name, e.g. IRIS Network")
	FsEditFanToken.String(FlagMaxSupply, "", "The maximum supply of the fantoken")
	FsEditFanToken.String(FlagMintable, "", "Whether the fantoken can be minted, default to false")

	FsTransferFanTokenOwner.String(FlagRecipient, "", "The new owner")

	FsMintFanToken.String(FlagRecipient, "", "Address to which the fantoken is to be minted")
	FsMintFanToken.String(FlagAmount, "", "Amount of the fantoken to be minted")
}
