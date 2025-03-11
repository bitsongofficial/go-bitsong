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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

const dummyUpgradeHeight = 5

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
	preModule appmodule.HasPreBlocker
}

func (s *UpgradeTestSuite) SetupTest() {
	s.Setup()
	s.preModule = upgrade.NewAppModule(s.App.AppKeepers.UpgradeKeeper, addresscodec.NewBech32Codec("bitsong"))
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) TestUpgrade() {
	upgradeSetup := func(zeroDel bool) {
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 2)

		vals, _ := s.App.AppKeepers.StakingKeeper.GetAllValidators(s.Ctx)
		for i, val := range vals {

			valAddr, _ := sdk.ValAddressFromBech32(val.OperatorAddress)
			val.Tokens = math.LegacyNewDecFromInt(val.Tokens).MulTruncate(math.LegacyOneDec().Sub(math.LegacyNewDecWithPrec(1, 3))).RoundInt() // 1 % slash
			//manually set validator historic rewards
			s.App.AppKeepers.DistrKeeper.SetValidatorHistoricalRewards(
				s.Ctx,
				valAddr,
				1,
				types.NewValidatorHistoricalRewards(sdk.NewDecCoins(sdk.NewDecCoin("ubtsg", math.OneInt())), 2),
			)

			err := s.App.AppKeepers.StakingKeeper.SetValidator(s.Ctx, val)
			if err != nil {
				panic(err)
			}
			// store 0 share delegation to fist validator if test requires
			if i == 0 && zeroDel {
				s.App.AppKeepers.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.NewDelegation(s.TestAccs[0].String(), val.OperatorAddress, math.LegacyZeroDec()))
			}
		}
	}

	postUpgrade := func() {
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 1)
		// vals, _ := s.App.AppKeepers.StakingKeeper.GetAllValidators(s.Ctx)

		// for i, val := range vals {

		// }
		// withdraw all delegator rewards
		dels, err := s.App.AppKeepers.StakingKeeper.GetAllDelegations(s.Ctx)
		if err != nil {
			panic(err)
		}
		for _, delegation := range dels {
			valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
			if err != nil {
				panic(err)
			}

			delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
			if err != nil {
				panic(err)
			}
			fmt.Println("~~~~~~~~~ POST UPGRADE DEBUG ~~~~~~~~~~~~")
			fmt.Printf("delAddr: %v\n", delAddr)
			fmt.Printf("valAddr: %v\n", valAddr)
			_, err = s.App.AppKeepers.DistrKeeper.WithdrawDelegationRewards(s.Ctx, delAddr, valAddr)
			if err != nil {
				panic(err)
			}
		}
	}

	testCases := []struct {
		name         string
		pre_upgrade  func()
		upgrade      func()
		post_upgrade func()
	}{
		{
			"test: missing slash event delegations patched",
			func() {
				upgradeSetup(false)
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func() {
				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 1)

				//  check rewards normall
			},
		},
		{
			"test: zero delegation removed from store",
			func() {
				upgradeSetup(true)
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func() {
				postUpgrade()

			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset

			tc.pre_upgrade()
			tc.upgrade()
			tc.post_upgrade()
		})
	}
}

func dummyUpgrade(s *UpgradeTestSuite) {
	s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
	plan := upgradetypes.Plan{Name: "v022", Height: dummyUpgradeHeight}
	err := s.App.AppKeepers.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
	s.Require().NoError(err)
	_, err = s.App.AppKeepers.UpgradeKeeper.GetUpgradePlan(s.Ctx)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithHeaderInfo(header.Info{Height: dummyUpgradeHeight, Time: s.Ctx.BlockTime().Add(time.Second)}).WithBlockHeight(dummyUpgradeHeight)
}
