package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
)

//var _ authexported.Account = (*BitSongAccount)(nil)

type BitSongAccount struct {
	*authtypes.BaseAccount

	Handle string `json:"handle" yaml:"handle"`
}

func NewBitSongAccount(acc authtypes.BaseAccount, handle string) *BitSongAccount {
	return &BitSongAccount{
		BaseAccount: &acc,
		Handle:      handle,
	}
}

type bAccountJSON struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	Handle        string         `json:"handle" yaml:"handle"`
}

func (bacc BitSongAccount) MarshalJSON() ([]byte, error) {
	alias := bAccountJSON{
		Address:       bacc.Address,
		Coins:         bacc.Coins,
		PubKey:        bacc.PubKey,
		AccountNumber: bacc.AccountNumber,
		Sequence:      bacc.Sequence,
		Handle:        bacc.Handle,
	}

	return codec.Cdc.MarshalJSON(alias)
}

func (bacc *BitSongAccount) UnmarshalJSON(bz []byte) error {
	var alias bAccountJSON
	if err := codec.Cdc.UnmarshalJSON(bz, &alias); err != nil {
		return err
	}

	bacc.BaseAccount = authtypes.NewBaseAccount(alias.Address, alias.Coins, alias.PubKey, alias.AccountNumber, alias.Sequence)
	bacc.Handle = alias.Handle

	return nil
}

func (bacc BitSongAccount) Validate() error {
	return bacc.BaseAccount.Validate()
}
