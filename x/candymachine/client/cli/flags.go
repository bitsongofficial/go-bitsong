package cli

import (
	flag "github.com/spf13/pflag"
)

const ()

func FlagCreateCandyMachine() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	return fs
}
