package simulation

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// RandomizedGenState generates a random GenesisState for distribution
func RandomizedGenState(simState *module.SimulationState) {
	// 1. Randomize module parameters
	bondDenom := simState.BondDenom

	issueFee := sdk.Coin{
		Denom:  bondDenom,
		Amount: math.NewInt(int64(simulation.RandIntBetween(simState.Rand, 100, 1000000))),
	}
	mintFee := sdk.Coin{
		Denom:  bondDenom,
		Amount: math.NewInt(int64(simulation.RandIntBetween(simState.Rand, 10, 100000))),
	}
	burnFee := sdk.Coin{
		Denom:  bondDenom,
		Amount: math.NewInt(int64(simulation.RandIntBetween(simState.Rand, 10, 100000))),
	}

	params := types.Params{
		IssueFee: issueFee,
		MintFee:  mintFee,
		BurnFee:  burnFee,
	}

	// 2. Generate random fantokens
	numFanTokens := simState.Rand.Intn(5) + 1 // 1 to 5 tokens
	fanTokens := make([]types.FanToken, numFanTokens)

	for i := 0; i < numFanTokens; i++ {
		denom := "ft" + simtypes.RandStringOfLength(simState.Rand, 6)
		maxSupply := int64(simulation.RandIntBetween(simState.Rand, 1000, 10000000))

		// Select a random account for minter
		minterAcc := simState.Accounts[simState.Rand.Intn(len(simState.Accounts))]
		minter := minterAcc.Address.String()

		// Random metadata
		name := fmt.Sprintf("Random FanToken %d", i)
		symbol := strings.ToUpper(simtypes.RandStringOfLength(simState.Rand, 4))

		var uri string
		if simState.Rand.Intn(2) == 0 {
			uri = fmt.Sprintf("https://example.com/token/%s", simtypes.RandStringOfLength(simState.Rand, 8))
		}

		authorityAcc := simState.Accounts[simState.Rand.Intn(len(simState.Accounts))]
		authority := authorityAcc.Address.String()

		metaData := types.Metadata{
			Name:      name,
			Symbol:    symbol,
			URI:       uri,
			Authority: authority,
		}

		fanTokens[i] = types.FanToken{
			Denom:     denom,
			MaxSupply: math.NewInt(maxSupply),
			Minter:    minter,
			MetaData:  metaData,
		}
	}

	// 3. Construct and marshal the randomized genesis state
	genesisState := &types.GenesisState{
		Params:    params,
		FanTokens: fanTokens,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)
}
