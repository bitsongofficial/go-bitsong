package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/keeper"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) CreateNFT(creator sdk.AccAddress, collectionId uint64) *types.MsgCreateNFTResponse {
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
		creator, collectionId, creator.String(), "Punk", "punk.com", 0, false, false, []types.Creator{}, 1,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) CreateMutableNFT(creator sdk.AccAddress, collectionId uint64) *types.MsgCreateNFTResponse {
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
		creator, collectionId, creator.String(), "Punk", "punk.com", 0, false, true, []types.Creator{}, 1,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) CreateNFTWithCreators(creator sdk.AccAddress, collectionId uint64, creatorAccs []sdk.AccAddress) *types.MsgCreateNFTResponse {
	creators := []types.Creator{}
	for _, creatorAcc := range creatorAccs {
		creators = append(creators, types.Creator{
			Address:  creatorAcc.String(),
			Verified: false,
			Share:    100,
		})
	}
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
		creator, collectionId, creator.String(), "Punk", "punk.com", 0, false, false, creators, 1,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) CreateCollection(creator sdk.AccAddress) *types.MsgCreateCollectionResponse {
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateCollection(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateCollection(
		creator, "PUNK", "Punk Collection", "punk.com", creator.String(), false,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) TestMsgServerCreateNFT() {
	tests := []struct {
		testCase           string
		nftId              uint64
		expectPass         bool
		expectedNFTId      string
		expectedMetadataId uint64
	}{
		{
			"create an nft",
			0,
			true,
			"1:1:0",
			1,
		},
	}

	for _, tc := range tests {
		creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

		// set params for issue fee
		issuePrice := sdk.NewInt64Coin("stake", 1000000)
		suite.app.NFTKeeper.SetParamSet(suite.ctx, types.Params{
			IssuePrice: issuePrice,
		})

		// mint coins for issue fee
		suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{issuePrice})
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, creator, sdk.Coins{issuePrice})

		collInfo := suite.CreateCollection(creator)

		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
			creator, collInfo.Id, creator.String(), "Punk", "punk.com", 0, false, false, []types.Creator{
				{
					Address:  creator.String(),
					Verified: true,
					Share:    1,
				},
			}, 1,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// test response is correct
			suite.Require().Equal(resp.MetadataId, tc.expectedMetadataId)
			suite.Require().Equal(resp.Id, tc.expectedNFTId)

			// test lastmetadataId and lastNftId are updated correctly
			lastMetadataId := suite.app.NFTKeeper.GetLastMetadataId(suite.ctx, collInfo.Id)
			suite.Require().Equal(lastMetadataId, tc.expectedMetadataId)

			// test Verified field false
			metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, collInfo.Id, resp.MetadataId)
			suite.Require().NoError(err)
			suite.Require().Equal(len(metadata.Creators), 1)
			suite.Require().Equal(metadata.Creators[0].Verified, false)

			// test metadataId and nftId to set correctly
			nft, err := suite.app.NFTKeeper.GetNFTById(suite.ctx, resp.Id)
			suite.Require().NoError(err)
			suite.Require().Equal(nft.Id(), tc.expectedNFTId)
			suite.Require().Equal(nft.MetadataId, tc.expectedMetadataId)

			// test fees are paid correctly
			balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, creator)
			suite.Require().Equal(balances, sdk.Coins{})
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerPrintEdition() {
	metadataAuthority := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	editionOwner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase      string
		collId        uint64
		metadataId    uint64
		sender        sdk.AccAddress
		editionOwner  string
		masterEdition *types.MasterEdition
		expectPass    bool
	}{
		{
			"metadata does not exist",
			1,
			0,
			metadataAuthority,
			editionOwner.String(),
			&types.MasterEdition{
				Supply:    1,
				MaxSupply: 2,
			},
			false,
		},
		{
			"not metadata authority",
			1,
			1,
			editionOwner,
			editionOwner.String(),
			&types.MasterEdition{
				Supply:    1,
				MaxSupply: 2,
			},
			false,
		},
		{
			"empty master edition",
			1,
			1,
			metadataAuthority,
			editionOwner.String(),
			nil,
			false,
		},
		{
			"exceed max supply",
			1,
			1,
			metadataAuthority,
			editionOwner.String(),
			&types.MasterEdition{
				Supply:    2,
				MaxSupply: 2,
			},
			false,
		},
		{
			"master edition nft check",
			0,
			1,
			metadataAuthority,
			editionOwner.String(),
			&types.MasterEdition{
				Supply:    1,
				MaxSupply: 2,
			},
			true,
		},
		{
			"successful printing",
			1,
			1,
			metadataAuthority,
			editionOwner.String(),
			&types.MasterEdition{
				Supply:    1,
				MaxSupply: 2,
			},
			true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.testCase, func() {
			suite.app.NFTKeeper.SetMetadata(suite.ctx, types.Metadata{
				CollId:               tc.collId,
				Id:                   1,
				MetadataAuthority:    metadataAuthority.String(),
				MintAuthority:        metadataAuthority.String(),
				Name:                 "meta1",
				Uri:                  "uri1",
				SellerFeeBasisPoints: 10,
				Creators: []types.Creator{
					{
						Address:  metadataAuthority.String(),
						Verified: false,
						Share:    1,
					},
				},
				PrimarySaleHappened: false,
				IsMutable:           true,
				MasterEdition:       tc.masterEdition,
			})
			suite.app.NFTKeeper.SetNFT(suite.ctx, types.NFT{
				CollId:     tc.collId,
				MetadataId: 1,
				Seq:        0,
				Owner:      metadataAuthority.String(),
			})

			// set params for issue fee
			issuePrice := sdk.NewInt64Coin("stake", 1000000)
			suite.app.NFTKeeper.SetParamSet(suite.ctx, types.Params{
				IssuePrice: issuePrice,
			})

			// mint coins for issue fee
			suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{issuePrice})
			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, tc.sender, sdk.Coins{issuePrice})

			// get old balance for future check
			oldBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.sender, "stake")

			msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
			resp, err := msgServer.PrintEdition(sdk.WrapSDKContext(suite.ctx), types.NewMsgPrintEdition(
				tc.sender, tc.collId, tc.metadataId, tc.editionOwner,
			))
			if tc.expectPass {
				suite.Require().NoError(err)

				// metadata supply change check
				meta, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, 1, tc.metadataId)
				suite.Require().NoError(err)
				suite.Require().Equal(meta.MasterEdition.Supply, tc.masterEdition.Supply+1)

				// nft data check (edition, id)
				nft, err := suite.app.NFTKeeper.GetNFTById(suite.ctx, resp.Id)
				suite.Require().NoError(err)
				suite.Require().Equal(nft.Id(), resp.Id)
				suite.Require().Equal(nft.Seq, tc.masterEdition.Supply)

				// nft issue fee check
				newBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.sender, "stake")
				suite.Require().Equal(newBalance.Amount.Int64()+1000000, oldBalance.Amount.Int64())
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerTransferNFT() {

	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collInfo1 := suite.CreateCollection(creator1)
	collInfo2 := suite.CreateCollection(creator2)
	nftInfo1 := suite.CreateNFT(creator1, collInfo1.Id)
	nftInfo2 := suite.CreateNFT(creator1, collInfo1.Id)
	nftInfo3 := suite.CreateNFT(creator2, collInfo2.Id)

	tests := []struct {
		testCase   string
		nftId      string
		sender     sdk.AccAddress
		target     string
		expectPass bool
	}{
		{
			"transfer not existing nft",
			"",
			creator3,
			creator1.String(),
			false,
		},
		{
			"transfer my nft to other",
			nftInfo1.Id,
			creator1,
			creator3.String(),
			true,
		},
		{
			"transfer other's nft",
			nftInfo2.Id,
			creator3,
			creator1.String(),
			false,
		},
		{
			"transfer nft to original address",
			nftInfo2.Id,
			creator1,
			creator1.String(),
			true,
		},
		{
			"transfer nft to empty address",
			nftInfo3.Id,
			creator2,
			creator2.String(),
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		_, err := msgServer.TransferNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgTransferNFT(
			tc.sender, tc.nftId, tc.target,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			nft, err := suite.app.NFTKeeper.GetNFTById(suite.ctx, tc.nftId)
			suite.Require().NoError(err)
			suite.Require().Equal(nft.Id(), tc.nftId)
			suite.Require().Equal(nft.Owner, tc.target)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerSignMetadata() {
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collInfo1 := suite.CreateCollection(creator1)
	nftInfo := suite.CreateNFTWithCreators(creator1, collInfo1.Id, []sdk.AccAddress{creator1, creator2})

	tests := []struct {
		testCase   string
		metadataId uint64
		sender     sdk.AccAddress
		expectPass bool
	}{
		{
			"not existing metadata",
			0,
			creator3,
			false,
		},
		{
			"sign correct metadata",
			nftInfo.MetadataId,
			creator1,
			true,
		},
		{
			"try sign metadata - not mine",
			nftInfo.MetadataId,
			creator3,
			false,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		_, err := msgServer.SignMetadata(sdk.WrapSDKContext(suite.ctx), types.NewMsgSignMetadata(
			tc.sender, collInfo1.Id, tc.metadataId,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, collInfo1.Id, tc.metadataId)
			suite.Require().NoError(err)
			suite.Require().Equal(metadata.Id, tc.metadataId)

			for _, creator := range metadata.Creators {
				if creator.Address == tc.sender.String() {
					suite.Require().Equal(creator.Verified, true)
				}
			}
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerUpdateMetadata() {
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collInfo1 := suite.CreateCollection(creator1)
	immutableNft := suite.CreateNFT(creator1, collInfo1.Id)
	collInfo2 := suite.CreateCollection(creator2)
	mutableNft := suite.CreateMutableNFT(creator2, collInfo2.Id)

	tests := []struct {
		testCase   string
		collId     uint64
		metadataId uint64
		sender     sdk.AccAddress
		expectPass bool
	}{
		{
			"not existing metadata",
			collInfo1.Id,
			0,
			creator3,
			false,
		},
		{
			"try updating not mutable metadata",
			immutableNft.CollId,
			immutableNft.MetadataId,
			creator1,
			false,
		},
		{
			"try updating not owned metadata",
			mutableNft.CollId,
			mutableNft.MetadataId,
			creator3,
			false,
		},
		{
			"update with correct values",
			mutableNft.CollId,
			mutableNft.MetadataId,
			creator2,
			true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.testCase, func() {
			msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
			msg := types.NewMsgUpdateMetadata(
				tc.sender, tc.collId, tc.metadataId, "NewPUNK", "NewURI", 10, []types.Creator{
					{Address: creator1.String(), Verified: true, Share: 100},
				},
			)

			_, err := msgServer.UpdateMetadata(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expectPass {
				suite.Require().NoError(err)

				metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, tc.collId, tc.metadataId)
				suite.Require().NoError(err)
				suite.Require().Equal(metadata.Id, tc.metadataId)

				suite.Require().Equal(len(metadata.Creators), 1)
				suite.Require().Equal(metadata.Creators[0].Verified, false)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerUpdateMetadataAuthority() {
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collInfo1 := suite.CreateCollection(creator1)
	immutableNft := suite.CreateNFT(creator1, collInfo1.Id)
	collInfo2 := suite.CreateCollection(creator2)
	mutableNft := suite.CreateMutableNFT(creator2, collInfo2.Id)

	tests := []struct {
		testCase   string
		collId     uint64
		metadataId uint64
		sender     sdk.AccAddress
		newOwner   string
		expectPass bool
	}{
		{
			"not existing metadata",
			collInfo1.Id,
			0,
			creator3,
			creator3.String(),
			false,
		},
		{
			"try updating not owned metadata",
			collInfo2.Id,
			mutableNft.MetadataId,
			creator3,
			creator3.String(),
			false,
		},
		{
			"update with correct value",
			collInfo2.Id,
			mutableNft.MetadataId,
			creator2,
			creator3.String(),
			true,
		},
		{
			"update with original value",
			collInfo1.Id,
			immutableNft.MetadataId,
			creator1,
			creator1.String(),
			true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.testCase, func() {
			msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
			_, err := msgServer.UpdateMetadataAuthority(sdk.WrapSDKContext(suite.ctx), types.NewMsgUpdateMetadataAuthority(
				tc.sender, tc.collId, tc.metadataId, tc.newOwner),
			)
			if tc.expectPass {
				suite.Require().NoError(err)

				metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, tc.collId, tc.metadataId)
				suite.Require().NoError(err)
				suite.Require().Equal(metadata.Id, tc.metadataId)
				suite.Require().Equal(metadata.MetadataAuthority, tc.newOwner)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerCreateCollection() {
	tests := []struct {
		testCase             string
		expectPass           bool
		expectedCollectionId uint64
	}{
		{
			"create a collection",
			true,
			1,
		},
	}

	for _, tc := range tests {
		creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		resp, err := msgServer.CreateCollection(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateCollection(
			creator, "PUNK", "Punk Collection", "punk.com", creator.String(), false,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// test response is correct
			suite.Require().Equal(resp.Id, tc.expectedCollectionId)

			// test last collectionId id updated correctly
			lastCollectionId := suite.app.NFTKeeper.GetLastCollectionId(suite.ctx)
			suite.Require().Equal(lastCollectionId, tc.expectedCollectionId)

			// test collection is set correctly
			collection, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, resp.Id)
			suite.Require().NoError(err)
			suite.Require().Equal(collection.Id, tc.expectedCollectionId)
			suite.Require().Equal(collection.UpdateAuthority, creator.String())
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerUpdateCollectionAuthority() {
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collectionInfo1 := suite.CreateCollection(creator1)
	collectionInfo2 := suite.CreateCollection(creator2)

	tests := []struct {
		testCase     string
		sender       sdk.AccAddress
		targetOwner  string
		collectionId uint64
		expectPass   bool
	}{
		{
			"update collection authority with owner",
			creator1,
			creator2.String(),
			collectionInfo1.Id,
			true,
		},
		{
			"try updating collection authority with non-owner",
			creator1,
			creator2.String(),
			collectionInfo2.Id,
			false,
		},
	}

	for _, tc := range tests {

		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		_, err := msgServer.UpdateCollectionAuthority(sdk.WrapSDKContext(suite.ctx), types.NewMsgUpdateCollectionAuthority(
			tc.sender, tc.collectionId, tc.targetOwner,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// test authority is updated correctly
			collection, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, tc.collectionId)
			suite.Require().NoError(err)
			suite.Require().Equal(collection.Id, tc.collectionId)
			suite.Require().Equal(collection.UpdateAuthority, tc.targetOwner)
		} else {
			suite.Require().Error(err)
		}
	}
}
