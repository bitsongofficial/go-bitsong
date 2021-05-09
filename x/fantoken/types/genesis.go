package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var nativeToken = FanToken{
	Denom:     sdk.DefaultBondDenom,
	Name:      "Network staking token",
	MaxSupply: sdk.NewInt(10000000000),
	Mintable:  true,
	Owner:     sdk.AccAddress(crypto.AddressHash([]byte(ModuleName))).String(),
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens []FanToken) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

// SetNativeToken resets the system's default native token
func SetNativeToken(
	denom string,
	name string,
	maxSupply sdk.Int,
	mintable bool,
	owner sdk.AccAddress,
) {
	nativeToken = NewFanToken(denom, name, maxSupply, mintable, owner)
}

//GetNativeToken returns the system's default native token
func GetNativeToken() FanToken {
	return nativeToken
}

// ValidateGenesis validates the provided token genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	// validate token
	for _, token := range data.Tokens {
		if err := ValidateToken(token); err != nil {
			return err
		}
	}

	// validate token
	for _, coin := range data.BurnedCoins {
		if err := coin.Validate(); err != nil {
			return err
		}
	}
	return nil
}
