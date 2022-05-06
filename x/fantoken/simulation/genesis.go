package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// Simulation parameter constants
const (
	IssueFee = "issue_fee"
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

	var issueFee sdk.Int
	var fantokens []tokentypes.FanToken

	simState.AppParams.GetOrGenerate(
		simState.Cdc, IssueFee, &issueFee, simState.Rand,
		func(r *rand.Rand) {
			issueFee = sdk.NewInt(int64(10))

			for i := 0; i < 5; i++ {
				fantokens = append(fantokens, randFanToken(r, simState.Accounts))
			}
		},
	)

	fantokenGenesis := tokentypes.NewGenesisState(
		tokentypes.NewParams(sdk.NewCoin(sdk.DefaultBondDenom, issueFee)),
		fantokens,
	)

	bz, err := json.MarshalIndent(&fantokenGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", tokentypes.ModuleName, bz)

	simState.GenState[tokentypes.ModuleName] = simState.Cdc.MustMarshalJSON(&fantokenGenesis)
}
