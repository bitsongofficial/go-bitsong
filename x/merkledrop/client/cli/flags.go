package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagProofs      = "proofs"
	FlagIndex       = "index"
	FlagStartHeight = "start-height"
	FlagEndHeight   = "end-height"
	FlagAmount      = "amount"
	FlagDenom       = "denom"
	FlagIPFSNode    = "ipfs-node"
)

func FlagsCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Int64(FlagStartHeight, 0, "Start height of the merkledrop")
	fs.Int64(FlagEndHeight, 0, "End height of the merkledrop")
	fs.String(FlagDenom, "", "Denom of the merkledrop")
	fs.String(FlagIPFSNode, "localhost:5001", "IPFS node to use")

	return fs
}

func FlagClaimMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagProofs, "", "Merkle proofs of the merkledrop")
	fs.Int64(FlagAmount, 0, "Amount of the merkledrop")
	fs.Uint64(FlagIndex, 0, "Index of the merkledrop")

	return fs
}

type accountInput struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
	Index   uint64 `json:"index"`
}

type accountsInput struct {
	Accounts []accountInput `json:"accounts"`
}
