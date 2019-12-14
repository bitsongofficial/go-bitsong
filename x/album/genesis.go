package album

import (
	"bytes"
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

const (
	// Default period for deposits
	DefaultPeriod time.Duration = 86400 * 2 * time.Second // 2 days
)

// GenesisState - all album state that must be provided at genesis
type GenesisState struct {
	StartingAlbumID uint64        `json:"starting_album_id"`
	Albums          Albums        `json:"albums"`
	Deposits        Deposits      `json:"deposits"`
	DepositParams   DepositParams `json:"deposit_params"`
}

// NewGenesisState creates a new genesis state for the album module
func NewGenesisState(startingAlbumID uint64, dp DepositParams) GenesisState {
	return GenesisState{
		StartingAlbumID: startingAlbumID,
		DepositParams:   dp,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(10)

	return GenesisState{
		StartingAlbumID: 1,
		DepositParams: DepositParams{
			MinDeposit:       sdk.Coins{sdk.NewCoin(btsg.BondDenom, minDepositTokens)},
			MaxDepositPeriod: DefaultPeriod,
		},
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

	if !data.DepositParams.MinDeposit.IsValid() {
		return fmt.Errorf("Album deposit amount must be a valid sdk.Coins amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetAlbumID(ctx, data.StartingAlbumID)
	k.SetDepositParams(ctx, data.DepositParams)

	// check if the deposits pool account exists
	moduleAcc := k.Sk.GetModuleAccount(ctx, ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	var totalDeposits sdk.Coins
	for _, deposit := range data.Deposits {
		k.SetDeposit(ctx, deposit.AlbumID, deposit.Depositor, deposit)
		totalDeposits = totalDeposits.Add(deposit.Amount)
	}

	for _, album := range data.Albums {
		k.SetAlbum(ctx, album)
	}

	// add coins if not provided on genesis
	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(totalDeposits); err != nil {
			panic(err)
		}
		k.Sk.SetModuleAccount(ctx, moduleAcc)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingAlbumID, _ := k.GetAlbumID(ctx)
	depositParams := k.GetDepositParams(ctx)
	// TODO: export only verified albums?
	albums := k.GetAlbumsFiltered(ctx, sdk.AccAddress{}, types.StatusVerified, 0)

	var albumsDeposits Deposits
	for _, album := range albums {
		deposits := k.GetDeposits(ctx, album.AlbumID)
		albumsDeposits = append(albumsDeposits, deposits...)
	}

	return GenesisState{
		StartingAlbumID: startingAlbumID,
		Deposits:        albumsDeposits,
		DepositParams:   depositParams,
		Albums:          albums,
	}
}
