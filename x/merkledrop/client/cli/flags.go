package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMerkleRoot  = "merkle-root"
	FlagTotalAmount = "total-amount"
	FlagProofs      = "proofs"
	FlagAmount      = "amount"
	FlagIndex       = "index"
)

func FlagCreateMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagMerkleRoot, "", "Merkle root of the merkledrop")
	fs.String(FlagTotalAmount, "", "Total amount of the merkledrop")

	return fs
}

func FlagClaimMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagProofs, "", "Merkle proofs of the merkledrop")
	fs.String(FlagAmount, "", "Amount of the merkledrop")
	fs.Uint64(FlagIndex, 0, "Index of the merkledrop")

	return fs
}
