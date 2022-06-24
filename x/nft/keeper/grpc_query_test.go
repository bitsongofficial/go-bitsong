package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestGRPCNFTInfo() {
	// create nfts
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	nftInfo1 := suite.CreateNFT(creator, 1)
	nftInfo2 := suite.CreateNFT(creator, 1)

	tests := []struct {
		testCase           string
		nftId              string
		expectPass         bool
		expectedNFTId      string
		expectedMetadataId uint64
	}{
		{
			"not existing nft id query",
			"",
			false,
			"",
			1,
		},
		{
			"query for existing nft 1",
			nftInfo1.Id,
			true,
			nftInfo1.Id,
			nftInfo1.MetadataId,
		},
		{
			"query for existing nft 2",
			nftInfo2.Id,
			true,
			nftInfo2.Id,
			nftInfo2.MetadataId,
		},
	}

	for _, tc := range tests {
		resp, err := suite.app.NFTKeeper.NFTInfo(sdk.WrapSDKContext(suite.ctx), &types.QueryNFTInfoRequest{
			Id: tc.nftId,
		})
		if tc.expectPass {
			suite.Require().NoError(err)
			suite.Require().Equal(resp.Nft.Id(), tc.expectedNFTId)
			suite.Require().Equal(resp.Metadata.Id, tc.expectedMetadataId)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestGRPCNFTsByOwner() {
	// create nfts
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	suite.CreateNFT(creator1, 1)
	suite.CreateNFT(creator1, 1)
	suite.CreateNFT(creator2, 1)

	tests := []struct {
		testCase        string
		owner           string
		expectPass      bool
		expectedNFTsLen int
	}{
		{
			"empty address",
			"",
			false,
			0,
		},
		{
			"invalid address",
			"0xAddressWrong",
			false,
			0,
		},
		{
			"creator1 address",
			creator1.String(),
			true,
			2,
		},
		{
			"creator2 address",
			creator2.String(),
			true,
			1,
		},
		{
			"creator3 address",
			creator3.String(),
			true,
			0,
		},
	}

	for _, tc := range tests {
		resp, err := suite.app.NFTKeeper.NFTsByOwner(sdk.WrapSDKContext(suite.ctx), &types.QueryNFTsByOwnerRequest{
			Owner: tc.owner,
		})
		if tc.expectPass {
			suite.Require().NoError(err)
			suite.Require().Equal(len(resp.Nfts), tc.expectedNFTsLen)
			suite.Require().Equal(len(resp.Metadata), tc.expectedNFTsLen)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestGRPCMetadata() {
	// create nfts
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	nftInfo1 := suite.CreateNFT(creator1, 1)
	nftInfo2 := suite.CreateNFT(creator2, 1)

	tests := []struct {
		testCase          string
		id                uint64
		expectPass        bool
		expectedAuthority string
	}{
		{
			"not existing id",
			0,
			false,
			"",
		},
		{
			"metadata1",
			nftInfo1.MetadataId,
			true,
			creator1.String(),
		},
		{
			"metadata2",
			nftInfo2.MetadataId,
			true,
			creator2.String(),
		},
	}

	for _, tc := range tests {
		resp, err := suite.app.NFTKeeper.Metadata(sdk.WrapSDKContext(suite.ctx), &types.QueryMetadataRequest{
			Id: tc.id,
		})
		if tc.expectPass {
			suite.Require().NoError(err)
			suite.Require().Equal(resp.Metadata.UpdateAuthority, tc.expectedAuthority)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestGRPCCollection() {
	// create nfts
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	collectionInfo1 := suite.CreateCollection(creator)
	collectionInfo2 := suite.CreateCollection(creator)
	suite.CreateNFT(creator, collectionInfo1.Id)
	suite.CreateNFT(creator, collectionInfo1.Id)

	tests := []struct {
		testCase          string
		id                uint64
		expectPass        bool
		expectedAuthority string
		expectedNftsCount int
	}{
		{
			"not existing id",
			0,
			false,
			"",
			0,
		},
		{
			"collection1",
			collectionInfo1.Id,
			true,
			creator.String(),
			2,
		},
		{
			"collection2",
			collectionInfo2.Id,
			true,
			creator.String(),
			0,
		},
	}

	for _, tc := range tests {
		resp, err := suite.app.NFTKeeper.Collection(sdk.WrapSDKContext(suite.ctx), &types.QueryCollectionRequest{
			Id: tc.id,
		})
		if tc.expectPass {
			suite.Require().NoError(err)
			suite.Require().Equal(resp.Collection.UpdateAuthority, tc.expectedAuthority)
			suite.Require().Equal(len(resp.NftIds), tc.expectedNftsCount)
		} else {
			suite.Require().Error(err)
		}
	}
}
