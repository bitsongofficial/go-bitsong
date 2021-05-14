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
)

var (
	FsIssueFanToken         = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditFanToken          = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferFanTokenOwner = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintFanToken          = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueFanToken.String(FlagSymbol, "", "The token symbol. Once created, it cannot be modified")
	FsIssueFanToken.String(FlagName, "", "The token name, e.g. Bitsong Network")
	FsIssueFanToken.String(FlagMaxSupply, "", "The maximum supply of the token")
	FsIssueFanToken.Bool(FlagMintable, false, "Whether the token can be minted, default to false")
	FsIssueFanToken.String(FlagDescription, "", "The token description")

	FsEditFanToken.String(FlagName, "[do-not-modify]", "The token name, e.g. IRIS Network")
	FsEditFanToken.String(FlagMaxSupply, "", "The maximum supply of the token")
	FsEditFanToken.String(FlagMintable, "", "Whether the token can be minted, default to false")

	FsTransferFanTokenOwner.String(FlagRecipient, "", "The new owner")

	FsMintFanToken.String(FlagRecipient, "", "Address to which the token is to be minted")
	FsMintFanToken.String(FlagAmount, "", "Amount of the token to be minted")
}
