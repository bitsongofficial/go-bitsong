package types

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Player struct {
	ID      PlayerID       `json:"id" yaml:"id"`
	Moniker string         `json:"moniker" yaml:"moniker"`
	Deposit sdk.Coin       `json:"deposit" yaml:"deposit"`
	Owner   sdk.AccAddress `json:"owner" yaml:"owner"`
}

func (p Player) String() string {
	return fmt.Sprintf(`Player
  ID:      %s
  Owner:   %s
  Deposit: %s
  Moniker: %s`, p.ID, p.Owner, p.Deposit, p.Moniker)
}

func (p Player) Validate() error {
	if p.Owner == nil || p.Owner.Empty() {
		return fmt.Errorf("invalid owner")
	}

	if p.Deposit.Denom != btsg.BondDenom {
		return fmt.Errorf("invelid deposit")
	}

	if p.Moniker == "" || len(p.Moniker) < MinMonikerLength || len(p.Moniker) > MaxMonikerLength {
		return fmt.Errorf("invalid moniker")
	}

	return nil
}
