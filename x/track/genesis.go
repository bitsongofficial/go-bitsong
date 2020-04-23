package track

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	// set initial track id
	k.SetLastTrackID(ctx, data.LastTrackID)

	for _, track := range data.Tracks {
		k.SetTrack(ctx, track)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	lastTrackID, _ := k.GetLastTrackID(ctx)

	return GenesisState{
		LastTrackID: lastTrackID,
		Tracks:      k.GetTracks(ctx),
	}
}
