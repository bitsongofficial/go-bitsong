package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateGenesis(t *testing.T) {

	p := InitialPool()
	require.Nil(t, p.ValidateGenesis())

	p2 := Pool{Rewards: sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(-1)}}}
	require.NotNil(t, p2.ValidateGenesis())

}
