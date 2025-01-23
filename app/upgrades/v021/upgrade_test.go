package v021_test

import (
	"fmt"
	"testing"

	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/stretchr/testify/suite"
)

const dummyUpgradeHeight = 5

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
}

func (s *UpgradeTestSuite) SetupTest() {
	s.Setup()
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) TestUpgrade() {
	upgradeSetup := func() {
		// create delegations with 0 power
		validators, _ := s.App.AppKeepers.StakingKeeper.GetAllValidators(s.Ctx)
		// del2 := s.TestAccs[1]
		for _, val := range validators {
			dels, _ := s.App.AppKeepers.StakingKeeper.GetValidatorDelegations(s.Ctx, sdktypes.ValAddress(val.OperatorAddress))

			// log the current delegations
			log := fmt.Sprintf("Current delegations for validator %s: %#v", val.OperatorAddress, dels)
			fmt.Println(log)

			// update delegation with 0 shares
			for _, del := range dels {
				del.Shares = math.LegacyZeroDec()
			}

			// s.Ctx.Logger().Info(fmt.Sprintf("Current delegations for validator %s: %#v", val.OperatorAddress, dels))
			// // create another delegator
			// s.FundAcc(del2, types.NewCoins(types.NewCoin("ubtsg", math.NewInt(1000000))))
			// s.StakingHelper.Delegate(del2, valAddr, math.NewInt(1000000))
		}
	}

	testCases := []struct {
		name         string
		pre_upgrade  func()
		upgrade      func()
		post_upgrade func()
	}{
		{
			"Test that the upgrade succeeds",
			func() {
				upgradeSetup()
			},
			func() {
				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
				plan := upgradetypes.Plan{Name: "v021", Height: dummyUpgradeHeight}
				err := s.App.AppKeepers.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
				s.Require().NoError(err)
				_, err = s.App.AppKeepers.UpgradeKeeper.GetUpgradePlan(s.Ctx)
				s.Require().NoError(err)
				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
				s.Require().NotPanics(func() {
					s.App.BeginBlocker(s.Ctx)
				})

			},
			func() {
				// confirm 0 delegation does not exists
				dels, _ := s.App.AppKeepers.StakingKeeper.GetAllDelegations(s.Ctx)
				for _, del := range dels {
					s.Require().NotZero(del.Shares)
				}
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
