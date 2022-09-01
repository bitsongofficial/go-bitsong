package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/launchpad/keeper"
	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestMsgServerCreateLaunchPad() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase        string
		sender          sdk.AccAddress
		collectionOwner sdk.AccAddress
		collectionFee   sdk.Coin
		collectionId    uint64
		maxMint         uint64
		shuffle         bool
		expectPass      bool
	}{
		{
			"when collection is not available",
			addr1,
			addr1,
			sdk.NewInt64Coin("ubtsg", 0),
			0,
			1000,
			false,
			false,
		},
		{
			"when collection is not owned by the sender",
			addr1,
			addr2,
			sdk.NewInt64Coin("ubtsg", 0),
			1,
			1000,
			false,
			false,
		},
		{
			"when max mint is bigger than max mint params",
			addr1,
			addr2,
			sdk.NewInt64Coin("ubtsg", 0),
			1,
			10000000000,
			false,
			false,
		},
		{
			"successful launchpad creation when fee is positive",
			addr1,
			addr1,
			sdk.NewInt64Coin("ubtsg", 1000),
			1,
			1000,
			false,
			true,
		},
		{
			"successful launchpad creation when fee is zero",
			addr1,
			addr1,
			sdk.NewInt64Coin("ubtsg", 0),
			1,
			1000,
			false,
			true,
		},
		{
			"successful launchpad creation with shuffle",
			addr1,
			addr1,
			sdk.NewInt64Coin("ubtsg", 0),
			1,
			1000,
			false,
			true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.testCase, func() {
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(tc.collectionFee))
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, tc.sender, sdk.NewCoins(tc.collectionFee))
			suite.Require().NoError(err)

			if tc.collectionId > 0 {
				suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
					Id:              tc.collectionId,
					Symbol:          "PUNK",
					UpdateAuthority: tc.collectionOwner.String(),
				})
			}

			params := suite.app.LaunchPadKeeper.GetParamSet(suite.ctx)
			params.LaunchpadCreationPrice = tc.collectionFee
			suite.app.LaunchPadKeeper.SetParamSet(suite.ctx, params)

			msgServer := keeper.NewMsgServerImpl(suite.app.LaunchPadKeeper)
			pad := types.LaunchPad{
				CollId:               tc.collectionId,
				Price:                0,
				Treasury:             addr1.String(),
				Denom:                "ubtsg",
				GoLiveDate:           1659870342,
				EndTimestamp:         0,
				MaxMint:              1000,
				Minted:               0,
				Authority:            tc.sender.String(),
				MetadataBaseUrl:      "https://punk.com/metadata",
				Mutable:              true,
				SellerFeeBasisPoints: 100,
				Creators:             []nfttypes.Creator(nil),
				Shuffle:              tc.shuffle,
			}

			oldSenderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.sender, "ubtsg")
			_, err = msgServer.CreateLaunchPad(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateLaunchPad(
				tc.sender, pad,
			))
			if tc.expectPass {
				suite.Require().NoError(err)

				// check collection authority is upgraded
				moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
				collection, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, tc.collectionId)
				suite.Require().NoError(err)
				suite.Require().Equal(collection.UpdateAuthority, moduleAddr.String())

				// check launchpad object is created from params
				savedPad, err := suite.app.LaunchPadKeeper.GetLaunchPadByCollId(suite.ctx, tc.collectionId)
				suite.Require().NoError(err)
				suite.Require().Equal(pad, savedPad)

				// check fee spent when it is positive
				newSenderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.sender, "ubtsg")
				suite.Require().Equal(newSenderBalance.Add(tc.collectionFee), oldSenderBalance)

				metadataIds := suite.app.LaunchPadKeeper.GetMintableMetadataIds(suite.ctx, tc.collectionId)
				suite.Require().Len(metadataIds, int(pad.MaxMint))

				// check if not ordered if not shuffle
				ordered := true
				for i, metadataId := range metadataIds {
					if i > 0 && metadataId < metadataIds[i-1] {
						ordered = false
					}
				}
				if !pad.Shuffle {
					suite.Require().True(ordered)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerUpdateLaunchPad() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase     string
		sender       sdk.AccAddress
		padAuthority sdk.AccAddress
		collectionId uint64
		expectPass   bool
	}{
		{
			"when launchpad is not available case",
			addr1,
			addr1,
			0,
			false,
		},
		{
			"when sender is not the authority",
			addr1,
			addr2,
			1,
			false,
		},
		{
			"successful launchpad upgrade",
			addr1,
			addr1,
			1,
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.LaunchPadKeeper)
		pad := types.LaunchPad{
			CollId:               1,
			Price:                0,
			Treasury:             tc.padAuthority.String(),
			Denom:                "ubtsg",
			GoLiveDate:           1659870342,
			EndTimestamp:         0,
			MaxMint:              1000,
			Minted:               0,
			Authority:            tc.padAuthority.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		}

		if tc.collectionId > 0 {
			suite.app.LaunchPadKeeper.SetLaunchPad(suite.ctx, pad)
		}

		pad.MetadataBaseUrl = "https://punk.com/newmeatadata"
		_, err := msgServer.UpdateLaunchPad(sdk.WrapSDKContext(suite.ctx), types.NewMsgUpdateLaunchPad(
			tc.sender, pad,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// check launchpad upgraded
			savedPad, err := suite.app.LaunchPadKeeper.GetLaunchPadByCollId(suite.ctx, pad.CollId)
			suite.Require().NoError(err)
			suite.Require().Equal(pad, savedPad)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerCloseLaunchPad() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase     string
		sender       sdk.AccAddress
		padAuthority sdk.AccAddress
		collectionId uint64
		expectPass   bool
	}{
		{
			"when launchpad is not available case",
			addr1,
			addr1,
			0,
			false,
		},
		{
			"when sender is not the authority",
			addr1,
			addr2,
			1,
			false,
		},
		{
			"successful launchpad close",
			addr1,
			addr1,
			1,
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.LaunchPadKeeper)
		pad := types.LaunchPad{
			CollId:               1,
			Price:                0,
			Treasury:             tc.padAuthority.String(),
			Denom:                "ubtsg",
			GoLiveDate:           1659870342,
			EndTimestamp:         0,
			MaxMint:              1000,
			Minted:               0,
			Authority:            tc.padAuthority.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		}

		if tc.collectionId > 0 {
			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
			suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
				Id:              tc.collectionId,
				Symbol:          "PUNK",
				UpdateAuthority: moduleAddr.String(),
			})
			suite.app.LaunchPadKeeper.SetLaunchPad(suite.ctx, pad)
		}

		pad.MetadataBaseUrl = "https://punk.com/newmeatadata"
		_, err := msgServer.CloseLaunchPad(sdk.WrapSDKContext(suite.ctx), types.NewMsgCloseLaunchPad(
			tc.sender, tc.collectionId,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// check launchpad deleted
			_, err := suite.app.LaunchPadKeeper.GetLaunchPadByCollId(suite.ctx, pad.CollId)
			suite.Require().Error(err)

			// check collection ownership transferred
			collection, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, tc.collectionId)
			suite.Require().NoError(err)
			suite.Require().Equal(collection.UpdateAuthority, tc.sender.String())
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerMintNFT() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase     string
		sender       sdk.AccAddress
		collectionId uint64
		expectPass   bool
	}{
		{
			"when launchpad is not available case",
			addr1,
			0,
			false,
		},
		{
			"successful launchpad mint nft",
			addr1,
			1,
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.LaunchPadKeeper)
		pad := types.LaunchPad{
			CollId:               1,
			Price:                0,
			Treasury:             addr1.String(),
			Denom:                "ubtsg",
			GoLiveDate:           1659870342,
			EndTimestamp:         0,
			MaxMint:              1000,
			Minted:               0,
			Authority:            addr1.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		}

		if tc.collectionId > 0 {
			suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
				Id:              tc.collectionId,
				Symbol:          "PUNK",
				UpdateAuthority: addr1.String(),
			})
			err := suite.app.LaunchPadKeeper.CreateLaunchPad(suite.ctx, &types.MsgCreateLaunchPad{
				Sender: addr1.String(),
				Pad:    pad,
			})
			suite.Require().NoError(err)
		}

		_, err := msgServer.MintNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgMintNFT(
			tc.sender, tc.collectionId, "punk1",
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// check launchpad minted count increased
			pad, err := suite.app.LaunchPadKeeper.GetLaunchPadByCollId(suite.ctx, pad.CollId)
			suite.Require().NoError(err)
			suite.Require().Equal(pad.Minted, uint64(1))

			// check nft minted
			nfts := suite.app.NFTKeeper.GetCollectionNfts(suite.ctx, tc.collectionId)
			suite.Require().NoError(err)
			suite.Require().Len(nfts, 1)

			// check minted nft ownership
			suite.Require().Equal(nfts[0].Owner, tc.sender.String())
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerMintNFTs() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase   string
		sender     sdk.AccAddress
		maxMint    uint64
		numMint    uint64
		expectPass bool
	}{
		{
			"when exceeding total number on launchpad",
			addr1,
			10,
			100,
			false,
		},
		{
			"successful launchpad mint nft",
			addr1,
			10,
			9,
			true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.testCase, func() {
			suite.SetupTest()
			msgServer := keeper.NewMsgServerImpl(suite.app.LaunchPadKeeper)
			suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
				Id:              1,
				Symbol:          "PUNK",
				UpdateAuthority: addr1.String(),
			})

			pad := types.LaunchPad{
				CollId:               1,
				Price:                0,
				Treasury:             addr1.String(),
				Denom:                "ubtsg",
				GoLiveDate:           1659870342,
				EndTimestamp:         0,
				MaxMint:              tc.maxMint,
				Minted:               0,
				Authority:            addr1.String(),
				MetadataBaseUrl:      "https://punk.com/metadata",
				Mutable:              true,
				SellerFeeBasisPoints: 100,
				Creators:             []nfttypes.Creator(nil),
			}

			err := suite.app.LaunchPadKeeper.CreateLaunchPad(suite.ctx, &types.MsgCreateLaunchPad{
				Sender: addr1.String(),
				Pad:    pad,
			})
			suite.Require().NoError(err)

			_, err = msgServer.MintNFTs(sdk.WrapSDKContext(suite.ctx), types.NewMsgMintNFTs(
				tc.sender, 1, tc.numMint,
			))
			if tc.expectPass {
				suite.Require().NoError(err)

				// check launchpad minted count increased
				pad, err := suite.app.LaunchPadKeeper.GetLaunchPadByCollId(suite.ctx, pad.CollId)
				suite.Require().NoError(err)
				suite.Require().Equal(pad.Minted, uint64(tc.numMint))

				// check nft minted
				nfts := suite.app.NFTKeeper.GetCollectionNfts(suite.ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Len(nfts, int(tc.numMint))

				// check minted nft ownership
				suite.Require().Equal(nfts[0].Owner, tc.sender.String())
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
