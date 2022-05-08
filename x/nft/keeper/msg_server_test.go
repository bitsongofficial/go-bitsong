package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/keeper"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// TODO: test
// SignMetadata
// UpdateMetadata
// UpdateMetadataAuthority
// CreateCollection
// VerifyCollection
// UnverifyCollection
// UpdateCollectionAuthority

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
