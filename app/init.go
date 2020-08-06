package app

import (
	"github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Init initializes the application, overriding the default genesis states that should be changed
func Init() {
	mint.DefaultGenesisState = mintGenesisState
	staking.DefaultGenesisState = stakingGenesisState
	gov.DefaultGenesisState = govGenesisState
}

// stakingGenesisState returns the default genesis state for the staking module, replacing the
// bond denom from stake to desmos
func stakingGenesisState() staking.GenesisState {
	return staking.GenesisState{
		Params: staking.NewParams(
			staking.DefaultUnbondingTime,
			staking.DefaultMaxValidators,
			staking.DefaultMaxEntries,
			0,
			types.BondDenom,
		),
	}
}

func govGenesisState() gov.GenesisState {
	return gov.NewGenesisState(
		1,
		gov.NewDepositParams(
			sdk.NewCoins(sdk.NewCoin(types.BondDenom, govTypes.DefaultMinDepositTokens)),
			gov.DefaultPeriod,
		),
		gov.NewVotingParams(gov.DefaultPeriod),
		gov.NewTallyParams(govTypes.DefaultQuorum, govTypes.DefaultThreshold, govTypes.DefaultVeto),
	)
}

func mintGenesisState() mint.GenesisState {
	return mint.GenesisState{
		Params: mint.NewParams(
			types.BondDenom,
			sdk.NewDecWithPrec(13, 2),
			sdk.NewDecWithPrec(20, 2),
			sdk.NewDecWithPrec(7, 2),
			sdk.NewDecWithPrec(67, 2),
			uint64(60*60*8766/5), // assuming 5 second block times
		),
	}
}
