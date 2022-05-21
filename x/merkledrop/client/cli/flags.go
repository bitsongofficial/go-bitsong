package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMerkleRoot = "merkle-root"
	FlagCoin       = "coin"
	FlagProofs     = "proofs"
	FlagIndex      = "index"
)

func FlagCreateMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagMerkleRoot, "", "Merkle root of the merkledrop")
	fs.String(FlagCoin, "", "Coin of the merkledrop")

	return fs
}

func FlagClaimMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagProofs, "", "Merkle proofs of the merkledrop")
	fs.String(FlagCoin, "", "Coin to claim of the merkledrop")
	fs.Uint64(FlagIndex, 0, "Index of the merkledrop")

	return fs
}
