package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMerkleRoot = "merkle-root"
	FlagCoin       = "coin"
	FlagProofs     = "proofs"
	FlagIndex      = "index"
	FlagStartTime  = "start-time"
	FlagEndTime    = "end-time"
)

func FlagCreateMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagMerkleRoot, "", "Merkle root of the merkledrop")
	fs.String(FlagCoin, "", "Coin of the merkledrop")
	fs.String(FlagStartTime, "", "Start time of the merkledrop")
	fs.String(FlagEndTime, "", "End time of the merkledrop")

	return fs
}

func FlagClaimMerkledrop() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagProofs, "", "Merkle proofs of the merkledrop")
	fs.String(FlagCoin, "", "Coin to claim of the merkledrop")
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
