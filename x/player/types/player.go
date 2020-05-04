package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Player struct {
	Moniker    string         `json:"moniker" yaml:"moniker"`
	PlayerAddr sdk.AccAddress `json:"player_addr" yaml:"player_addr"`
	Validator  sdk.ValAddress `json:"validator" yaml:"validator"`
}

func (p Player) String() string {
	return fmt.Sprintf(`Player
  Player: %s
  Validator: %s
  Moniker:   %s`, p.PlayerAddr, p.Validator, p.Moniker)
}

func (p Player) Validate() error {
	if p.Validator == nil || p.Validator.Empty() {
		return fmt.Errorf("invalid validator")
	}

	if p.Moniker == "" || len(p.Moniker) < MinMonikerLength || len(p.Moniker) > MaxMonikerLength {
		return fmt.Errorf("invalid moniker")
	}

	if p.PlayerAddr == nil || p.PlayerAddr.Empty() {
		return fmt.Errorf("invalid player address")
	}

	return nil
}
