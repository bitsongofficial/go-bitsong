package v021_test

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
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

const dummyUpgradeHeight = 5

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
	upgradeSetup := func(shares, slash math.LegacyDec, jailed bool) {
		vals, _ := s.App.StakingKeeper.GetAllValidators(s.Ctx)
		fmt.Printf("# OF VALS: %d\n", len(vals))
		for _, val := range vals {
			// delAddrStr, err :=s.App.AccountKeeper.AddressCodec().BytesToString(s.TestAccs[0])

			s.FundAcc(s.TestAccs[0], sdktypes.NewCoins(sdktypes.NewCoin("stake", math.NewInt(1000000))))
			_, err := s.App.StakingKeeper.Delegate(s.Ctx, s.TestAccs[0], math.NewInt(1000000), stakingtypes.Unbonded, val, true)
			s.Require().NoError(err)

			if !slash.IsZero() {
				val.Tokens = math.LegacyNewDecFromInt(val.Tokens).MulTruncate(math.LegacyOneDec().Sub(slash)).RoundInt() // 1 % slash
			}
			if jailed {
				val.Jailed = jailed
			}
			err = s.App.StakingKeeper.SetValidator(s.Ctx, val)
			s.Require().NoError(err)

			// get delegations
			dels, err := s.App.StakingKeeper.GetValidatorDelegations(s.Ctx, sdktypes.ValAddress(val.OperatorAddress))
			s.Require().NoError(err)
			fmt.Printf("# OF DELS: %d\n", len(dels))

			// todo: fix staking helper to propoerly manage staking in simulation tests
		}

	}

	testCases := []struct {
		name         string
		pre_upgrade  func()
		upgrade      func()
		post_upgrade func()
	}{
		{
			"test app module params",
			func() {
				upgradeSetup(math.LegacyOneDec(), math.LegacyZeroDec(), false)
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func() {
				//  cheeck icq params
				params := s.App.ICQKeeper.GetParams(s.Ctx)
				fmt.Println(params.HostEnabled)
				s.Require().Equal(len(params.AllowQueries), 1)
				s.Require().True(params.HostEnabled)

				// check smart account params
				smartAccParams := s.App.SmartAccountKeeper.GetParams(s.Ctx)
				authManagers := s.App.AuthenticatorManager.GetRegisteredAuthenticators()
				s.Require().Greater(len(authManagers), 1)
				s.Require().True(smartAccParams.IsSmartAccountActive)

				//check ibchook params
				ibcwasmparams := s.App.IBCKeeper.ClientKeeper.GetParams(s.Ctx)
				s.Require().Equal(len(ibcwasmparams.AllowedClients), 2)

				// check cadance params
				cadanceParams := s.App.CadanceKeeper.GetParams(s.Ctx)
				s.Require().Equal(cadanceParams.ContractGasLimit, uint64(1000000))

				// expidited proposal
				govparams, _ := s.App.GovKeeper.Params.Get(s.Ctx)
				newExpeditedVotingPeriod := time.Minute * 60 * 24
				s.Require().Equal(govparams.ExpeditedVotingPeriod.Seconds(), newExpeditedVotingPeriod.Seconds())
				s.Require().Equal(govparams.ExpeditedThreshold, "0.75")
			},
		},
		{
			"0 shares, jailed",
			func() {
				upgradeSetup(math.LegacyZeroDec(), math.LegacyZeroDec(), true)
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func() {
				dels, _ := s.App.StakingKeeper.GetAllDelegations(s.Ctx)
				for _, del := range dels {
					s.Require().NotZero(del.Shares)
				}

			},
		},
		{
			"test distribution reward invariants patched",
			func() {
				upgradeSetup(math.LegacyOneDec(), math.LegacyNewDecWithPrec(1, 3), true) // 1% slash
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func() {
				dels, _ := s.App.StakingKeeper.GetAllDelegations(s.Ctx)
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

func dummyUpgrade(s *UpgradeTestSuite) {
	s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
	plan := upgradetypes.Plan{Name: "v021", Height: dummyUpgradeHeight}
	err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
	s.Require().NoError(err)
	_, err = s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithHeaderInfo(header.Info{Height: dummyUpgradeHeight, Time: s.Ctx.BlockTime().Add(time.Second)}).WithBlockHeight(dummyUpgradeHeight)
}
