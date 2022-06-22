package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestNftGetSet() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// get nft by not available id
	_, err := suite.app.NFTKeeper.GetNFTById(suite.ctx, "")
	suite.Require().Error(err)

	// get nfts by owner when there's nothing
	nftsByOwner := suite.app.NFTKeeper.GetNFTsByOwner(suite.ctx, addr)
	suite.Require().Len(nftsByOwner, 0)

	// get all nfts when not available
	allNfts := suite.app.NFTKeeper.GetAllNFTs(suite.ctx)
	suite.Require().Len(allNfts, 0)

	// create a new nft
	nfts := []types.NFT{
		{
			MetadataId: 1,
			Owner:      addr.String(),
		},
		{
			MetadataId: 2,
			Owner:      addr.String(),
		},
		{
			MetadataId: 3,
			Owner:      addr2.String(),
		},
		{
			MetadataId: 4,
			Owner:      addr2.String(),
		},
	}

	for _, nft := range nfts {
		suite.app.NFTKeeper.SetNFT(suite.ctx, nft)
	}

	// check nfts existance by id
	for _, nft := range nfts {
		n, err := suite.app.NFTKeeper.GetNFTById(suite.ctx, nft.Id())
		suite.Require().NoError(err)
		suite.Require().Equal(nft, n)
	}

	// check nfts by owner after set
	nftsByOwner = suite.app.NFTKeeper.GetNFTsByOwner(suite.ctx, addr)
	suite.Require().Len(nftsByOwner, 2)

	// check nfts for whole export
	allNfts = suite.app.NFTKeeper.GetAllNFTs(suite.ctx)
	suite.Require().Len(allNfts, 4)
	suite.Require().Equal(nfts, allNfts)
}
