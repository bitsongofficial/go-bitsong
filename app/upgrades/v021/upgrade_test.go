package v021_test

import (
	"fmt"
	"testing"

	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"

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
	upgradeSetup := func() {}

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
				_, err = s.App.AppKeepers.UpgradeKeeper.GetUpgradePlan(s.Ctx)
				s.Require().NoError(err)

				s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight)
				s.Require().NotPanics(func() {
					s.App.BeginBlocker(s.Ctx)
				})
			},
			func() {},
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
