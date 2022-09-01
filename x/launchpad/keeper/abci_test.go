package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

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
			MaxMint:              10,
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
	}

	for _, pad := range pads {
		moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
		suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
			Id:              pad.CollId,
			Symbol:          "PUNK",
			UpdateAuthority: moduleAddr.String(),
		})
		suite.app.LaunchPadKeeper.SetLaunchPad(suite.ctx, pad)
	}

	endingPads := suite.app.LaunchPadKeeper.GetLaunchPadsToEndByTime(suite.ctx)
	suite.Require().Len(endingPads, 1)

	allPads := suite.app.LaunchPadKeeper.GetAllLaunchPads(suite.ctx)
	suite.Require().Len(allPads, 4)
	suite.Require().Equal(pads, allPads)

	suite.app.LaunchPadKeeper.EndBlocker(suite.ctx)

	// check one pad automatically closed
	allPads = suite.app.LaunchPadKeeper.GetAllLaunchPads(suite.ctx)
	suite.Require().Len(allPads, 3)
}
