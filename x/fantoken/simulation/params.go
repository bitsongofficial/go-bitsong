package simulation

import (
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

const (
	keyIssueFee = "IssueFee"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyIssueFee,
			func(r *rand.Rand) string {
				return RandomInt(r).String()
			},
		),
	}
}
