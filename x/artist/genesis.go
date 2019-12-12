package artist

import (
	"bytes"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"time"

	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Default period for deposits
	DefaultPeriod time.Duration = 86400 * 2 * time.Second // 2 days
)

// GenesisState - all artist state that must be provided at genesis
type GenesisState struct {
	StartingArtistID uint64        `json:"starting_artist_id"`
	Artists          Artists       `json:"artists"`
	Deposits         Deposits      `json:"deposits"`
	DepositParams    DepositParams `json:"deposit_params"`
}

// NewGenesisState creates a new genesis state for the artist module
func NewGenesisState(startingArtistID uint64, dp DepositParams) GenesisState {
	return GenesisState{
		StartingArtistID: startingArtistID,
		DepositParams:    dp,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(10)

	return GenesisState{
		StartingArtistID: 1,
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
		return fmt.Errorf("Artist deposit amount must be a valid sdk.Coins amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetArtistID(ctx, data.StartingArtistID)
	k.SetDepositParams(ctx, data.DepositParams)

	// check if the deposits pool account exists
	moduleAcc := k.Sk.GetModuleAccount(ctx, ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	var totalDeposits sdk.Coins
	for _, deposit := range data.Deposits {
		k.SetDeposit(ctx, deposit.ArtistID, deposit.Depositor, deposit)
		totalDeposits = totalDeposits.Add(deposit.Amount)
	}

	for _, artist := range data.Artists {
		k.SetArtist(ctx, artist)
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
	startingArtistID, _ := k.GetArtistID(ctx)
	depositParams := k.GetDepositParams(ctx)
	// TODO: export only verified artists?
	artists := k.GetArtistsFiltered(ctx, sdk.AccAddress{}, types.StatusVerified, 0)

	var artistsDeposits Deposits
	for _, artist := range artists {
		deposits := k.GetDeposits(ctx, artist.ArtistID)
		artistsDeposits = append(artistsDeposits, deposits...)
	}

	return GenesisState{
		StartingArtistID: startingArtistID,
		Deposits:         artistsDeposits,
		DepositParams:    depositParams,
		Artists:          artists,
	}
}
