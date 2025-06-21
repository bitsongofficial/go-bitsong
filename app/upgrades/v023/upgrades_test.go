package v023_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/header"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"

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

	}

	postUpgrade := func() {
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 1)

	}

	testCases := []struct {
		name         string
		pre_upgrade  func()
		upgrade      func()
		post_upgrade func()
	}{
		{
			"test: protocol-pool patch",
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
	err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
	s.Require().NoError(err)
	_, err = s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithHeaderInfo(header.Info{Height: dummyUpgradeHeight, Time: s.Ctx.BlockTime().Add(time.Second)}).WithBlockHeight(dummyUpgradeHeight)
}
