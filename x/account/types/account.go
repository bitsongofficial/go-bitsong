package types

import (
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
)

var _ authexported.Account = (*BitSongAccount)(nil)

type BitSongAccount struct {
	authexported.Account

	Handle string `json:"handle" yaml:"handle"`
}

func NewBitSongAccount(acc authexported.Account, handle string) *BitSongAccount {
	return &BitSongAccount{
		Account: acc,
		Handle:  handle,
	}
}
