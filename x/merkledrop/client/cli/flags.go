package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMerkleRoot  = "merkle-root"
	FlagTotalAmount = "total-amount"
)

func FlagCreateMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagMerkleRoot, "", "Merkle root of the merkledrop")
	fs.Uint64(FlagTotalAmount, 0, "Total amount of the merkledrop")

	return fs
}
