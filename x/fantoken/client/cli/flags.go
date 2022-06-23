package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSymbol       = "symbol"
	FlagName         = "name"
	FlagMaxSupply    = "max-supply"
	FlagRecipient    = "recipient"
	FlagNewAuthority = "new-authority"
	FlagNewMinter    = "new-minter"
	FlagAmount       = "amount"
	FlagURI          = "uri"
)

var (
	FsIssue        = flag.NewFlagSet("", flag.ContinueOnError)
	FsMint         = flag.NewFlagSet("", flag.ContinueOnError)
	FsDisableMint  = flag.NewFlagSet("", flag.ContinueOnError)
	FsSetAuthority = flag.NewFlagSet("", flag.ContinueOnError)
	FsSetMinter    = flag.NewFlagSet("", flag.ContinueOnError)
	FsSetUri       = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssue.String(FlagSymbol, "", "The fantoken symbol. Once created, it cannot be modified")
	FsIssue.String(FlagName, "", "The fantoken name, e.g. Bitsong Network")
	FsIssue.String(FlagMaxSupply, "", "The maximum supply of the fantoken")
	FsIssue.String(FlagURI, "", "The fantoken uri")

	FsMint.String(FlagRecipient, "", "Address to which the fantoken is to be minted")

	FsDisableMint.String(FlagName, "[do-not-modify]", "The fantoken name, e.g. IRIS Network")
	FsDisableMint.String(FlagMaxSupply, "", "The maximum supply of the fantoken")

	FsSetAuthority.String(FlagNewAuthority, "", "The new authority")

	FsSetMinter.String(FlagNewMinter, "", "The new minter")

	FsSetUri.String(FlagURI, "", "The uri of the fantoken")
}
