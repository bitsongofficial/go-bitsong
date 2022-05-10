package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/keeper"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) CreateNFT(creator sdk.AccAddress) *types.MsgCreateNFTResponse {
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
		creator, creator.String(), types.Data{
			Name:                 "Punk",
			Symbol:               "PUNK",
			Uri:                  "punk.com",
			SellerFeeBasisPoints: 0,
			Creators:             []*types.Creator{},
		}, false, false,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) CreateMutableNFT(creator sdk.AccAddress) *types.MsgCreateNFTResponse {
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
		creator, creator.String(), types.Data{
			Name:                 "Punk",
			Symbol:               "PUNK",
			Uri:                  "punk.com",
			SellerFeeBasisPoints: 0,
			Creators:             []*types.Creator{},
		}, false, true,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) CreateNFTWithCreators(creator sdk.AccAddress, creatorAccs []sdk.AccAddress) *types.MsgCreateNFTResponse {
	creators := []*types.Creator{}
	for _, creatorAcc := range creatorAccs {
		creators = append(creators, &types.Creator{
			Address:  creatorAcc.String(),
			Verified: false,
			Share:    100,
		})
	}
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
		creator, creator.String(), types.Data{
			Name:                 "Punk",
			Symbol:               "PUNK",
			Uri:                  "punk.com",
			SellerFeeBasisPoints: 0,
			Creators:             creators,
		}, false, false,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) CreateCollection(creator sdk.AccAddress) *types.MsgCreateCollectionResponse {
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.CreateCollection(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateCollection(
		creator, "Punk Collection", "punk.com", creator.String(),
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) VerifyCollection(sender sdk.AccAddress, collectionId, nftId uint64) *types.MsgVerifyCollectionResponse {
	msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
	resp, err := msgServer.VerifyCollection(sdk.WrapSDKContext(suite.ctx), types.NewMsgVerifyCollection(
		sender, collectionId, nftId,
	))
	suite.Require().NoError(err)
	return resp
}

func (suite *KeeperTestSuite) TestMsgServerCreateNFT() {
	tests := []struct {
		testCase           string
		nftId              uint64
		expectPass         bool
		expectedNFTId      uint64
		expectedMetadataId uint64
	}{
		{
			"create an nft",
			0,
			true,
			1,
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

		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		resp, err := msgServer.CreateNFT(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateNFT(
			creator, creator.String(), types.Data{
				Name:                 "Punk",
				Symbol:               "PUNK",
				Uri:                  "punk.com",
				SellerFeeBasisPoints: 0,
				Creators: []*types.Creator{
					{
						Address:  creator.String(),
						Verified: true,
						Share:    1,
					},
				},
			}, false, false,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// test response is correct
			suite.Require().Equal(resp.MetadataId, tc.expectedMetadataId)
			suite.Require().Equal(resp.Id, tc.expectedNFTId)

			// test lastmetadataId and lastNftId are updated correctly
			lastNftId := suite.app.NFTKeeper.GetLastNftId(suite.ctx)
			suite.Require().Equal(lastNftId, tc.expectedNFTId)
			lastMetadataId := suite.app.NFTKeeper.GetLastMetadataId(suite.ctx)
			suite.Require().Equal(lastMetadataId, tc.expectedMetadataId)

			// test Verified field false
			metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, resp.MetadataId)
			suite.Require().NoError(err)
			suite.Require().Equal(len(metadata.Data.Creators), 1)
			suite.Require().Equal(metadata.Data.Creators[0].Verified, false)

			// test metadataId and nftId to set correctly
			nft, err := suite.app.NFTKeeper.GetNFTById(suite.ctx, resp.Id)
			suite.Require().NoError(err)
			suite.Require().Equal(nft.Id, tc.expectedNFTId)
			suite.Require().Equal(nft.MetadataId, tc.expectedMetadataId)

			// test fees are paid correctly
			balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, creator)
			suite.Require().Equal(balances, sdk.Coins{})
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerTransferNFT() {

	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	nftInfo1 := suite.CreateNFT(creator1)
	nftInfo2 := suite.CreateNFT(creator1)
	nftInfo3 := suite.CreateNFT(creator2)

	tests := []struct {
		testCase   string
		nftId      uint64
		sender     sdk.AccAddress
		target     string
		expectPass bool
	}{
		{
			"transfer not existing nft",
			0,
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
			suite.Require().Equal(nft.Id, tc.nftId)
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
	nftInfo := suite.CreateNFTWithCreators(creator1, []sdk.AccAddress{creator1, creator2})

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
			tc.sender, tc.metadataId,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, tc.metadataId)
			suite.Require().NoError(err)
			suite.Require().Equal(metadata.Id, tc.metadataId)

			for _, creator := range metadata.Data.Creators {
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
	immutableNft := suite.CreateNFT(creator1)
	mutableNft := suite.CreateMutableNFT(creator2)

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
			"try updating not mutable metadata",
			immutableNft.MetadataId,
			creator1,
			false,
		},
		{
			"try updating not owned metadata",
			mutableNft.MetadataId,
			creator3,
			false,
		},
		{
			"update with correct values",
			mutableNft.MetadataId,
			creator2,
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		_, err := msgServer.UpdateMetadata(sdk.WrapSDKContext(suite.ctx), types.NewMsgUpdateMetadata(
			tc.sender, tc.metadataId, true, &types.Data{Name: "NewPUNK", Creators: []*types.Creator{
				{Address: creator1.String(), Verified: true, Share: 100},
			}},
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, tc.metadataId)
			suite.Require().NoError(err)
			suite.Require().Equal(metadata.Id, tc.metadataId)

			suite.Require().Equal(len(metadata.Data.Creators), 1)
			suite.Require().Equal(metadata.Data.Creators[0].Verified, false)
			suite.Require().Equal(metadata.PrimarySaleHappened, true)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerUpdateMetadataAuthority() {
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	immutableNft := suite.CreateNFT(creator1)
	mutableNft := suite.CreateMutableNFT(creator2)

	tests := []struct {
		testCase   string
		metadataId uint64
		sender     sdk.AccAddress
		newOwner   string
		expectPass bool
	}{
		{
			"not existing metadata",
			0,
			creator3,
			creator3.String(),
			false,
		},
		{
			"try updating not owned metadata",
			mutableNft.MetadataId,
			creator3,
			creator3.String(),
			false,
		},
		{
			"update with correct value",
			mutableNft.MetadataId,
			creator2,
			creator3.String(),
			true,
		},
		{
			"update with original value",
			immutableNft.MetadataId,
			creator1,
			creator1.String(),
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		_, err := msgServer.UpdateMetadataAuthority(sdk.WrapSDKContext(suite.ctx), types.NewMsgUpdateMetadataAuthority(
			tc.sender, tc.metadataId, tc.newOwner),
		)
		if tc.expectPass {
			suite.Require().NoError(err)

			metadata, err := suite.app.NFTKeeper.GetMetadataById(suite.ctx, tc.metadataId)
			suite.Require().NoError(err)
			suite.Require().Equal(metadata.Id, tc.metadataId)
			suite.Require().Equal(metadata.UpdateAuthority, tc.newOwner)
		} else {
			suite.Require().Error(err)
		}
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
			creator, "Punk Collection", "punk.com", creator.String(),
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

func (suite *KeeperTestSuite) TestMsgServerVerifyCollection() {
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collectionInfo := suite.CreateCollection(creator)
	nftInfo1 := suite.CreateNFT(creator)

	tests := []struct {
		testCase     string
		sender       sdk.AccAddress
		collectionId uint64
		nftId        uint64
		expectPass   bool
	}{
		{
			"verify collection with owner",
			creator,
			collectionInfo.Id,
			nftInfo1.Id,
			true,
		},
		{
			"try verifying collection with non-owner",
			creator2,
			collectionInfo.Id,
			nftInfo1.Id,
			false,
		},
	}

	for _, tc := range tests {

		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		_, err := msgServer.VerifyCollection(sdk.WrapSDKContext(suite.ctx), types.NewMsgVerifyCollection(
			tc.sender, tc.collectionId, tc.nftId,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// test number of nfts are correctly put on the collection
			nftIds := suite.app.NFTKeeper.GetCollectionNftRecords(suite.ctx, tc.collectionId)
			suite.Require().NoError(err)
			suite.Require().Equal(len(nftIds), 1)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerUnverifyCollection() {
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collectionInfo := suite.CreateCollection(creator)
	nftInfo1 := suite.CreateNFT(creator)
	nftInfo2 := suite.CreateNFT(creator)
	suite.VerifyCollection(creator, collectionInfo.Id, nftInfo1.Id)
	suite.VerifyCollection(creator, collectionInfo.Id, nftInfo2.Id)

	tests := []struct {
		testCase     string
		sender       sdk.AccAddress
		collectionId uint64
		nftId        uint64
		expectPass   bool
	}{
		{
			"unverify a nft on collection with owner",
			creator,
			collectionInfo.Id,
			nftInfo1.Id,
			true,
		},
		{
			"try unverifying collection with non-owner",
			creator2,
			collectionInfo.Id,
			nftInfo1.Id,
			false,
		},
	}

	for _, tc := range tests {

		msgServer := keeper.NewMsgServerImpl(suite.app.NFTKeeper)
		_, err := msgServer.UnverifyCollection(sdk.WrapSDKContext(suite.ctx), types.NewMsgUnverifyCollection(
			tc.sender, tc.collectionId, tc.nftId,
		))
		if tc.expectPass {
			suite.Require().NoError(err)
		} else {
			suite.Require().Error(err)
		}
	}
}
