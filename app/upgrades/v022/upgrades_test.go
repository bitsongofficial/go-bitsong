package v022_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/header"
	"cosmossdk.io/math"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

const dummyUpgradeHeight = 5

var testDescription = stakingtypes.NewDescription("test_moniker", "test_identity", "test_website", "test_security_contact", "test_details")

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
	preModule appmodule.HasPreBlocker
}

func (s *UpgradeTestSuite) SetupTest() {
	s.Setup()
	s.preModule = upgrade.NewAppModule(s.App.UpgradeKeeper, addresscodec.NewBech32Codec("bitsong"))
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) TestUpgrade() {
	upgradeSetup := func(zeroDel bool) {
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 2)

		vals, _ := s.App.StakingKeeper.GetAllValidators(s.Ctx)
		for i, val := range vals {

			valAddr, _ := sdk.ValAddressFromBech32(val.OperatorAddress)
			val.Tokens = math.LegacyNewDecFromInt(val.Tokens).MulTruncate(math.LegacyOneDec().Sub(math.LegacyNewDecWithPrec(1, 3))).RoundInt() // 1 % slash
			//manually set validator historic rewards
			s.App.DistrKeeper.SetValidatorHistoricalRewards(
				s.Ctx,
				valAddr,
				1,
				types.NewValidatorHistoricalRewards(sdk.NewDecCoins(sdk.NewDecCoin("ubtsg", math.OneInt())), 2), // set reference to 2
			)

			err := s.App.StakingKeeper.SetValidator(s.Ctx, val)
			s.Require().NoError(err)
			// store 0 share delegation to fist validator if test requires
			if i == 0 && zeroDel {
				s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.NewDelegation(s.TestAccs[0].String(), val.OperatorAddress, math.LegacyZeroDec()))
			}
		}

		// create another validator with normal state for control
		s.controlVal()

	}

	postUpgrade := func(sdk.ValAddress) {
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 1)

		// Verify control validator's rewards remain unchanged

		// withdraw all delegator rewards
		dels, err := s.App.StakingKeeper.GetAllDelegations(s.Ctx)
		s.Require().NoError(err)
		for _, delegation := range dels {
			valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
			s.Require().NoError(err)
			delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
			s.Require().NoError(err)
			fmt.Println("~~~~~~~~~ POST UPGRADE DEBUG ~~~~~~~~~~~~")
			fmt.Printf("delAddr: %v\n", delAddr)
			fmt.Printf("valAddr: %v\n", valAddr)
			_, err = s.App.DistrKeeper.WithdrawDelegationRewards(s.Ctx, delAddr, valAddr)
			s.Require().NoError(err)
		}
	}

	testCases := []struct {
		name         string
		pre_upgrade  func() (sdk.ValAddress, types.ValidatorHistoricalRewards)
		upgrade      func()
		post_upgrade func(sdk.ValAddress, types.ValidatorHistoricalRewards)
	}{
		{
			"test: missing slash event patched",
			func() (sdk.ValAddress, types.ValidatorHistoricalRewards) {
				upgradeSetup(false)
				controlValAddr, controlRewardsBefore := s.getControlValidatorState()
				return controlValAddr, controlRewardsBefore
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func(controlValAddr sdk.ValAddress, controlRewardsBefore types.ValidatorHistoricalRewards) {
				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 1)

				// Verify control validator's rewards remain unchanged
				controlRewardsAfter, err := s.App.DistrKeeper.GetValidatorHistoricalRewards(s.Ctx, controlValAddr, 1)
				s.Require().NoError(err, "Control validator rewards should still exist after upgrade")
				// Compare the control validator's rewards before and after upgrade
				s.Require().Equal(controlRewardsBefore.CumulativeRewardRatio, controlRewardsAfter.CumulativeRewardRatio,
					"Control validator rewards should be untouched")
			},
		},
		{
			"test: zero delegation removed from store",
			func() (sdk.ValAddress, types.ValidatorHistoricalRewards) {
				upgradeSetup(true)
				controlValAddr, controlRewardsBefore := s.getControlValidatorState()
				return controlValAddr, controlRewardsBefore
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func(controlValAddr sdk.ValAddress, controlRewardsBefore types.ValidatorHistoricalRewards) {
				postUpgrade(controlValAddr)

			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset
			controlValAddr, controlRewardsBefore := tc.pre_upgrade()
			tc.upgrade()
			tc.post_upgrade(controlValAddr, controlRewardsBefore)
		})
	}
}

func dummyUpgrade(s *UpgradeTestSuite) {
	s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
	plan := upgradetypes.Plan{Name: "v022", Height: dummyUpgradeHeight}
	err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
	s.Require().NoError(err)
	_, err = s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithHeaderInfo(header.Info{Height: dummyUpgradeHeight, Time: s.Ctx.BlockTime().Add(time.Second)}).WithBlockHeight(dummyUpgradeHeight)
}

func (s *UpgradeTestSuite) controlVal() {
	valPub := secp256k1.GenPrivKey().PubKey()
	valAddr := sdk.ValAddress(valPub.Address())
	ZeroCommission := stakingtypes.NewCommissionRates(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec())
	selfBond := sdk.NewCoins(sdk.Coin{Amount: math.OneInt(), Denom: "stake"})
	stakingCoin := sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: selfBond[0].Amount}
	s.FundAcc(sdk.AccAddress(valAddr), selfBond)
	valCreateMsg, err := stakingtypes.NewMsgCreateValidator(valAddr.String(), valPub, stakingCoin, testDescription, ZeroCommission, math.OneInt())
	s.Require().NoError(err)
	stakingMsgSvr := stakingkeeper.NewMsgServerImpl(s.App.StakingKeeper)
	res, err := stakingMsgSvr.CreateValidator(s.Ctx, valCreateMsg)
	fmt.Printf("err: %v\n", err)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	// set normal rewards
	s.App.DistrKeeper.SetValidatorHistoricalRewards(
		s.Ctx,
		valAddr,
		1,
		types.NewValidatorHistoricalRewards(sdk.NewDecCoins(sdk.NewDecCoin("ubtsg", math.NewInt(100))), 2), // set reference to 2
	)
	val, err := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.Require().NoError(err)
	err = s.App.StakingKeeper.SetValidator(s.Ctx, val)
}

// Helper function to get control validator state
func (s *UpgradeTestSuite) getControlValidatorState() (sdk.ValAddress, types.ValidatorHistoricalRewards) {
	// Find the control validator (the one created in controlVal function)
	vals, err := s.App.StakingKeeper.GetAllValidators(s.Ctx)
	s.Require().NoError(err)

	// The control validator should be the last one we added
	var controlValAddr sdk.ValAddress
	for _, val := range vals {
		// Find the validator that was created in controlVal() - you could add a distinctive
		// attribute in controlVal() to make this easier to identify
		if val.GetMoniker() == "test_moniker" && val.GetTokens().Equal(math.OneInt()) {
			fmt.Printf("val.OperatorAddress: %v\n", val.OperatorAddress)
			valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
			s.Require().NoError(err)
			controlValAddr = valAddr
			break
		}
	}

	s.Require().NotNil(controlValAddr, "Control validator not found")
	fmt.Printf("controlValAddr: %v\n", controlValAddr)

	// Get its historical rewards
	rewards, found := s.App.DistrKeeper.GetValidatorHistoricalRewards(s.Ctx, controlValAddr, 1)
	s.Require().NoError(found, "Control validator rewards should exist")

	return controlValAddr, rewards
}
