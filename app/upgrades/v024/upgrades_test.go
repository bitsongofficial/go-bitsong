package v024_test

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
	v023 "github.com/bitsongofficial/go-bitsong/app/upgrades/v023"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

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
	coin := types.NewCoin("ft12345", math.NewInt(12345))
	communityPoolFunds := types.NewCoins(types.NewCoin("ubtsg", math.NewInt(69420)), types.NewCoin("ftfantoken", math.NewInt(1234567890)))
	protocolpoolEscrow := protocolpooltypes.ProtocolPoolEscrowAccount
	protocolpool := protocolpooltypes.ModuleName
	distributionModule := distrtypes.ModuleName

	upgradeSetup := func() {
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 2)
		s.FundAcc(s.TestAccs[0], types.NewCoins(coin))

		// fund protocolpool_escrow (cannot use protocolpool fund-community-pool since we have disbabled external community pool, so we just seed with balance to module account)
		s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, communityPoolFunds)
		s.App.BankKeeper.SendCoinsFromModuleToModule(s.Ctx, minttypes.ModuleName, protocolpoolEscrow, communityPoolFunds)

	}

	postUpgrade := func() {
		s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight + 1)

		// assert protocolpool & protocolpool_escrow are empty
		protocolPoolBalance := s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.AccountKeeper.GetModuleAddress(protocolpool))
		protocolPoolEscrowBalance := s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.AccountKeeper.GetModuleAddress(protocolpoolEscrow))
		distrModuleBalance := s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.AccountKeeper.GetModuleAddress(distributionModule))

		fmt.Printf("distrModuleBalance: %v\n", distrModuleBalance)
		fmt.Printf("protocolPoolBalance: %v\n", protocolPoolBalance)
		fmt.Printf("protocolPoolEscrowBalance: %v\n", protocolPoolEscrowBalance)

		s.Require().Equal(protocolPoolBalance.Empty(), true)
		s.Require().Equal(protocolPoolEscrowBalance.Empty(), true)

		// assert fund community pool does not error
		err := s.App.DistrKeeper.FundCommunityPool(s.Ctx, types.NewCoins(coin), s.TestAccs[0])
		s.Require().NoError(err)

		// check for funded token in distribution module balance
		distrModuleBalance = s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.AccountKeeper.GetModuleAddress(distributionModule))
		amountOf := distrModuleBalance.AmountOf(coin.Denom)
		s.Require().Equal(amountOf, coin.Amount)

		// assert we are unblocking fund communityu pool requests
		s.Require().Equal(s.App.DistrKeeper.HasExternalCommunityPool(), false)

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
				upgradeSetup()
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
	plan := upgradetypes.Plan{Name: v023.UpgradeName, Height: dummyUpgradeHeight}
	err := s.App.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
	s.Require().NoError(err)
	_, err = s.App.UpgradeKeeper.GetUpgradePlan(s.Ctx)
	s.Require().NoError(err)

	s.Ctx = s.Ctx.WithHeaderInfo(header.Info{Height: dummyUpgradeHeight, Time: s.Ctx.BlockTime().Add(time.Second)}).WithBlockHeight(dummyUpgradeHeight)
}
