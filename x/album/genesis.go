package album

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

// GenesisState - all album state that must be provided at genesis
type GenesisState struct {
	StartingAlbumID uint64 `json:"starting_album_id"`
	Albums          Albums `json:"albums"`
}

// NewGenesisState creates a new genesis state for the album module
func NewGenesisState(startingAlbumID uint64) GenesisState {
	return GenesisState{
		StartingAlbumID: startingAlbumID,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingAlbumID: 1,
	}
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// ValidateGenesis validates the given genesis state and returns an error if something is invalid
func ValidateGenesis(data GenesisState) error {
	// TODO: add validation
	/*for _, record := range data.Artists {
		if err := record.Validate(); err != nil {
			return err
		}
	}*/

	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetAlbumID(ctx, data.StartingAlbumID)

	for _, album := range data.Albums {
		k.SetAlbum(ctx, album)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingAlbumID, _ := k.GetAlbumID(ctx)
	// TODO: export only verified albums?
	albums := k.GetAlbumsFiltered(ctx, sdk.AccAddress{}, types.StatusVerified, 0)

	return GenesisState{
		StartingAlbumID: startingAlbumID,
		Albums:          albums,
	}
}
