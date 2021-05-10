package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/bitsongofficial/bitsong/x/fantoken/types"
)

// Simulation parameter constants
const (
	IssuePrice = "issue_price"
)

// RandomDec randomized sdk.RandomDec
func RandomDec(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(r.Int63())
}

// RandomInt randomized sdk.Int
func RandomInt(r *rand.Rand) sdk.Int {
	return sdk.NewInt(r.Int63())
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {

	var issuePrice sdk.Int
	var tokens []types.FanToken

	simState.AppParams.GetOrGenerate(
		simState.Cdc, IssuePrice, &issuePrice, simState.Rand,
		func(r *rand.Rand) {
			issuePrice = sdk.NewInt(int64(10))

			for i := 0; i < 5; i++ {
				tokens = append(tokens, randFanToken(r, simState.Accounts))
			}
		},
	)

	tokenGenesis := types.NewGenesisState(
		types.NewParams(sdk.NewCoin(sdk.DefaultBondDenom, issuePrice)),
		tokens,
	)

	bz, err := json.MarshalIndent(&tokenGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenGenesis)
}
