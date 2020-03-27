package track

import (
	"bytes"
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"time"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

const (
	// Default period for deposits
	DefaultPeriod time.Duration = 86400 * 2 * time.Second // 2 days
)

// GenesisState - all track state that must be provided at genesis
type GenesisState struct {
	StartingTrackID uint64        `json:"starting_track_id"`
	Tracks          Tracks        `json:"tracks"`
	Deposits        Deposits      `json:"deposits"`
	DepositParams   DepositParams `json:"deposit_params"`
}

// NewGenesisState creates a new genesis state for the track module
func NewGenesisState(startingTrackID uint64, dp DepositParams) GenesisState {
	return GenesisState{
		StartingTrackID: startingTrackID,
		DepositParams:   dp,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(10)

	return GenesisState{
		StartingTrackID: 1,
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
		return fmt.Errorf("Track deposit amount must be a valid sdk.Coins amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, bankKeeper bank.Keeper, data GenesisState) {
	k.SetTrackID(ctx, data.StartingTrackID)
	k.SetDepositParams(ctx, data.DepositParams)

	// check if the deposits pool account exists
	moduleAcc := k.Sk.GetModuleAccount(ctx, ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	totalDeposits := sdk.NewCoins()
	for _, deposit := range data.Deposits {
		k.SetDeposit(ctx, deposit.TrackID, deposit.Depositor, deposit)
		totalDeposits = totalDeposits.Add(deposit.Amount...)
	}

	for _, track := range data.Tracks {
		k.SetTrack(ctx, track)
	}

	// add coins if not provided on genesis
	coins := bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if coins.IsZero() {
		if err := bankKeeper.SetBalances(ctx, moduleAcc.GetAddress(), totalDeposits); err != nil {
			panic(err)
		}
		k.Sk.SetModuleAccount(ctx, moduleAcc)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingTrackID, _ := k.GetTrackID(ctx)
	depositParams := k.GetDepositParams(ctx)
	// TODO: export only verified tracks?
	tracks := k.GetTracksFiltered(ctx, sdk.AccAddress{}, types.StatusVerified, 0)

	var tracksDeposits Deposits
	for _, track := range tracks {
		deposits := k.GetDeposits(ctx, track.TrackID)
		tracksDeposits = append(tracksDeposits, deposits...)
	}

	return GenesisState{
		StartingTrackID: startingTrackID,
		Deposits:        tracksDeposits,
		DepositParams:   depositParams,
		Tracks:          tracks,
	}
}
