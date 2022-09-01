package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestLaunchPadGetSetDelete() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	// get not available launchpad
	_, err := suite.app.LaunchPadKeeper.GetLaunchPadByCollId(suite.ctx, 0)
	suite.Require().Error(err)

	// get all launchpads when not available
	allPads := suite.app.LaunchPadKeeper.GetAllLaunchPads(suite.ctx)
	suite.Require().Len(allPads, 0)

	// get ending queue launchpads when not available
	endingPads := suite.app.LaunchPadKeeper.GetLaunchPadsToEndByTime(suite.ctx)
	suite.Require().Len(endingPads, 0)

	// set launchpads
	pads := []types.LaunchPad{
		{
			CollId:               1,
			Price:                100,
			Treasury:             addr.String(),
			Denom:                "ubtsg",
			GoLiveDate:           uint64(now.Unix()),
			EndTimestamp:         uint64(now.Unix()),
			MaxMint:              1000,
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
		{
			CollId:               2,
			Price:                100,
			Treasury:             addr.String(),
			Denom:                "ubtsg",
			GoLiveDate:           uint64(now.Unix()),
			EndTimestamp:         uint64(now.Unix()) + 100,
			MaxMint:              1000,
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
		{
			CollId:               3,
			Price:                100,
			Treasury:             addr.String(),
			Denom:                "ubtsg",
			GoLiveDate:           uint64(now.Unix()),
			EndTimestamp:         uint64(now.Unix()) - 100,
			MaxMint:              1000,
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
		{
			CollId:               4,
			Price:                100,
			Treasury:             addr.String(),
			Denom:                "ubtsg",
			GoLiveDate:           uint64(now.Unix()),
			EndTimestamp:         0,
			MaxMint:              1000,
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
	}

	for _, pad := range pads {
		suite.app.LaunchPadKeeper.SetLaunchPad(suite.ctx, pad)
	}

	for _, pad := range pads {
		m, err := suite.app.LaunchPadKeeper.GetLaunchPadByCollId(suite.ctx, pad.CollId)
		suite.Require().NoError(err)
		suite.Require().Equal(pad, m)
	}

	allPads = suite.app.LaunchPadKeeper.GetAllLaunchPads(suite.ctx)
	suite.Require().Len(allPads, 4)
	suite.Require().Equal(pads, allPads)

	endingPads = suite.app.LaunchPadKeeper.GetLaunchPadsToEndByTime(suite.ctx)
	suite.Require().Len(endingPads, 1)

	for _, pad := range pads {
		suite.app.LaunchPadKeeper.DeleteLaunchPad(suite.ctx, pad)
	}

	// get all launchpads after all deletion
	allPads = suite.app.LaunchPadKeeper.GetAllLaunchPads(suite.ctx)
	suite.Require().Len(allPads, 0)

	// get ending queue launchpads after all deletion
	endingPads = suite.app.LaunchPadKeeper.GetLaunchPadsToEndByTime(suite.ctx)
	suite.Require().Len(endingPads, 0)
}
