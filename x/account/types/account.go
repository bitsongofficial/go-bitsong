package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"time"
)

type Account struct {
	auth.BaseAccount

	Address     sdk.AccAddress `json:"address"`
	Handle      string         `json:"handle"`
	MetadataURI string         `json:"metadata_uri"`
	CreatedAt   time.Time      `json:"created_at"`
}

func NewAccount(accAddr sdk.AccAddress, handle, metadataURI string, createdAt time.Time) Account {
	return Account{
		Address:     accAddr,
		Handle:      handle,
		MetadataURI: metadataURI,
		CreatedAt:   createdAt,
	}
}

func (acc Account) String() string {
	return fmt.Sprintf(`
  Address:       %s
  Handle:        %s
  MetadataURI:   %s
  CreatedAt:     %s`,
		acc.Address.String(), acc.Handle, acc.MetadataURI, acc.CreatedAt.String())
}
