package app

import (
	"encoding/json"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for bitsong.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}

func NewDefaultGenesisStateWithCodec(cdc codec.JSONCodec) GenesisState {
	gen := ModuleBasics.DefaultGenesis(cdc)

	// here we override wasm config to make it permissioned by default
	wasmGen := wasmtypes.GenesisState{
		Params: wasmtypes.Params{
			CodeUploadAccess:             wasmtypes.AllowNobody,
			InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
		},
	}
	// other default genesis
	// mintGenesis := mintGenesisState()
	// stakingGenesis := stakingGenesisState()
	// govGenesis := govGenesisState()

	gen[wasmtypes.ModuleName] = cdc.MustMarshalJSON(&wasmGen)
	// gen["mint"] = cdc.MustMarshalJSON(mintGenesis)
	// gen["staking"] = cdc.MustMarshalJSON(stakingGenesis)
	// gen["gov"] = cdc.MustMarshalJSON(govGenesis)

	return gen
}

// stakingGenesisState returns the default genesis state for the staking module, replacing the
// bond denom from stake to ubtsg
// func stakingGenesisState() *stakingtypes.GenesisState {
// 	historicalEntries := uint32(0)
// 	return &stakingtypes.GenesisState{
// 		Params: stakingtypes.NewParams(
// 			stakingtypes.DefaultUnbondingTime,
// 			stakingtypes.DefaultMaxValidators,
// 			stakingtypes.DefaultMaxEntries,
// 			historicalEntries,
// 			"ubtsg", sdk.ZeroDec(),
// 		),
// 	}
// }

// func govGenesisState() *govv1beta1.GenesisState {
// 	startingPropId := uint64(1)
// 	return govv1beta1.NewGenesisState(
// 		startingPropId,
// 		govv1beta1.NewDepositParams(
// 			sdk.NewCoins(sdk.NewCoin("ubtsg", govv1beta1.DefaultMinDepositTokens)),
// 			govv1beta1.DefaultPeriod,
// 		),
// 		govv1beta1.NewVotingParams(govv1beta1.DefaultPeriod),
// 		govv1beta1.NewTallyParams(govv1beta1.DefaultQuorum, govv1beta1.DefaultThreshold, govv1beta1.DefaultVetoThreshold),
// 	)
// }

// func mintGenesisState() *minttypes.GenesisState {
// 	return &minttypes.GenesisState{
// 		Params: minttypes.NewParams(
// 			"ubtsg",
// 			sdk.NewDecWithPrec(13, 2),
// 			sdk.NewDecWithPrec(20, 2),
// 			sdk.NewDecWithPrec(7, 2),
// 			sdk.NewDecWithPrec(67, 2),
// 			uint64(60*60*8766/6), // assuming 6 second block times
// 		),
// 	}
// }
