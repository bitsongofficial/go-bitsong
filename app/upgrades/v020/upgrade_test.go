package v020_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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
		validators := s.App.AppKeepers.StakingKeeper.GetAllValidators(s.Ctx)
		for _, val := range validators {
			// update the tokens staked to validator due to slashing event
			// mimic slashing event on staking power, but not update slashing event to distribution module
			val.Tokens = math.LegacyNewDecFromInt(val.Tokens).MulTruncate(math.LegacyOneDec().Sub(math.LegacyNewDecWithPrec(1, 3))).RoundInt() // 1 % slash

			dels := s.App.AppKeepers.StakingKeeper.GetAllDelegations(s.Ctx)
			for _, del := range dels {
				endingPeriod := s.App.AppKeepers.DistrKeeper.IncrementValidatorPeriod(s.Ctx, val)
				// assert v018 bug is present prior to upgrade
				s.Require().Panics(func() {
					s.App.AppKeepers.DistrKeeper.CalculateDelegationRewards(s.Ctx, val, del, endingPeriod)
				})
			}
			// create another delegator
			del2 := s.TestAccs[2]
			s.FundAcc(del2, types.NewCoins(types.NewCoin("ubtsg", math.NewInt(1000000))))
			s.StakingHelper.Delegate(del2, val.GetOperator(), math.NewInt(1000000))
			del, found := s.App.AppKeepers.StakingKeeper.GetDelegation(s.Ctx, del2, val.GetOperator())
			s.Require().True(found)
			endingPeriod := s.App.AppKeepers.DistrKeeper.IncrementValidatorPeriod(s.Ctx, val)
			s.Require().NotPanics(func() {
				s.App.AppKeepers.DistrKeeper.CalculateDelegationRewards(s.Ctx, val, del, endingPeriod)
			})
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
				plan := upgradetypes.Plan{Name: "v020", Height: dummyUpgradeHeight}
				err := s.App.AppKeepers.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
				s.Require().NoError(err)
				_, exists := s.App.AppKeepers.UpgradeKeeper.GetUpgradePlan(s.Ctx)
				s.Require().True(exists)

				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
				s.Require().NotPanics(func() {
					beginBlockRequest := abcitypes.RequestBeginBlock{}
					s.App.BeginBlocker(s.Ctx, beginBlockRequest)
				})
			},
			func() {
				// assert rewards can be calculated
				validators := s.App.AppKeepers.StakingKeeper.GetAllValidators(s.Ctx)
				for _, val := range validators {
					dels := s.App.AppKeepers.StakingKeeper.GetAllDelegations(s.Ctx)
					for _, del := range dels {
						// confirm delegators can query, withdraw and stake
						// require all rewards to have been claimed for this delegator
						// confirm delegators claimed tokens was accurate
						s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 10)
						endingPeriod := s.App.AppKeepers.DistrKeeper.IncrementValidatorPeriod(s.Ctx, val)
						s.App.AppKeepers.DistrKeeper.CalculateDelegationRewards(s.Ctx, val, del, endingPeriod)
						s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 10)
						s.FundAcc(del.GetDelegatorAddr(), types.NewCoins(types.NewCoin("ubtsg", math.NewInt(1000000))))
						s.StakingHelper.Delegate(del.GetDelegatorAddr(), del.GetValidatorAddr(), math.NewInt(1000000))
						s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 10)
						_, err := s.App.AppKeepers.DistrKeeper.WithdrawDelegationRewards(s.Ctx, del.GetDelegatorAddr(), del.GetValidatorAddr())
						s.Require().NoError(err)
					}
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
