package v020_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/app/helpers"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	abci "github.com/cometbft/cometbft/abci/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestUpgradeSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))

}

func (s *UpgradeTestSuite) SetupTest() {
	V020Setup(s.T())
}

const dummyUpgradeHeight = 5

func (s *UpgradeTestSuite) dummyUpgrade() {
	s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
	plan := upgradetypes.Plan{Name: "v20", Height: dummyUpgradeHeight}
	err := s.App.AppKeepers.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
	s.Require().NoError(err)
	_, exists := s.App.AppKeepers.UpgradeKeeper.GetUpgradePlan(s.Ctx)
	s.Require().True(exists)

	s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
	s.Require().NotPanics(func() {
		beginBlockRequest := abcitypes.RequestBeginBlock{}
		s.App.BeginBlocker(s.Ctx, beginBlockRequest)
	})
}

// Setup initializes a new BitsongApp
func V020Setup(t *testing.T) *app.BitsongApp {
	t.Helper()

	privVal := helpers.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(appparams.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	app := app.SetupWithGenesisAccounts(t, valSet, []authtypes.GenesisAccount{acc}, balance)
	return app
}

func (s *UpgradeTestSuite) TestUpgrade() {

	// corrupt state to match mainnet
	s.setupCorruptedState()

	// upgrade software
	s.dummyUpgrade()
	s.App.BeginBlocker(s.Ctx, abci.RequestBeginBlock{})
	s.Ctx = s.Ctx.WithBlockTime(s.Ctx.BlockTime().Add(time.Hour * 24))

	// ensure post upgrade reward query works
	s.ensurePostUpgradeDistributionWorks()
}

func (s *UpgradeTestSuite) ensurePostUpgradeDistributionWorks() {

}

// setupCorruptedState aligns the testing environment with the mainnet state.
// By running this method, it will modify the delegations tokens from shares,
// reflecting a slashing event that
func (s *UpgradeTestSuite) setupCorruptedState() {
	// infractionHeight := int64(3)

	power := int64(1000000)
	slashFactor := sdk.NewDecWithPrec(1, 2) //1 %
	staking := s.App.AppKeepers.StakingKeeper

	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := staking.TokensFromConsensusPower(s.Ctx, power)
	slashAmountDec := sdk.NewDecFromInt(amount).Mul(slashFactor)
	slashAmount := slashAmountDec.TruncateInt()

	// mimic slashing event on staking power, but not update slashing event to distribution module
	validator, _ := staking.GetValidatorByConsAddr(s.Ctx, sdk.ConsAddress{})

	operatorAddress := validator.GetOperator()

	// call the before-modification hook
	if err := staking.Hooks().BeforeValidatorModified(s.Ctx, operatorAddress); err != nil {
		staking.Logger(s.Ctx).Error("failed to call before validator modified hook", "error", err)
	}
	remainingSlashAmount := slashAmount

	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(remainingSlashAmount, validator.Tokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, math.ZeroInt()) // defensive.

	// we need to calculate the *effective* slash fraction for distribution
	if validator.Tokens.IsPositive() {
		effectiveFraction := sdk.NewDecFromInt(tokensToBurn).QuoRoundUp(sdk.NewDecFromInt(validator.Tokens))
		// possible if power has changed
		if effectiveFraction.GT(math.LegacyOneDec()) {
			effectiveFraction = math.LegacyOneDec()
		}
		// call the before-slashed hook. Omitted to simulate v0.18.0 error bitsong is experiencing
		// if err := staking.Hooks().BeforeValidatorSlashed(s.Ctx, operatorAddress, effectiveFraction); err != nil {
		// 	staking.Logger(s.Ctx).Error("failed to call before validator slashed hook", "error", err)
		// }
		validator = staking.RemoveValidatorTokens(s.Ctx, validator, tokensToBurn)

		if err := s.burnBondedTokens(s.App.AppKeepers, s.Ctx, tokensToBurn); err != nil {
			panic(err)
		}
	}

	// Deduct from validator's bonded tokens and update the validator.
	// Burn the slashed tokens from the pool account and decrease the total supply.

}

func (s *UpgradeTestSuite) ensurePreUpgradeDistributionPanics() {
	// calculate rewards from distribution

	// expect difference due to unregistered slashing event from validator
}

// burnBondedTokens removes coins from the bonded pool module account
func (s *UpgradeTestSuite) burnBondedTokens(k keepers.AppKeepers, ctx sdk.Context, amt math.Int) error {
	if !amt.IsPositive() {
		// skip as no coins need to be burned
		return nil
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.StakingKeeper.BondDenom(ctx), amt))
	return k.BankKeeper.BurnCoins(ctx, types.BondedPoolName, coins)
}
