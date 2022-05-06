package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: test
// NFTsByOwner
// Metadata
// Collection

func (suite *KeeperTestSuite) TestGRPCNFTInfo() {
	// TODO: call a utility function to create an nft

	tests := []struct {
		testCase         string
		nftId            uint64
		expectPass       bool
		expectedNFT      types.NFT
		expectedMetadata types.Metadata
	}{
		{
			"not existing nft id query",
			0,
			false,
			types.NFT{},
			types.Metadata{},
		},
		// TODO: add case for querying by existing nftId
	}

	for _, tc := range tests {
		resp, err := suite.app.NFTKeeper.NFTInfo(sdk.WrapSDKContext(suite.ctx), &types.QueryNFTInfoRequest{})
		if tc.expectPass {
			suite.Require().NoError(err)
			suite.Require().Equal(resp.Nft, tc.expectedNFT)
			suite.Require().Equal(resp.Metadata, tc.expectedMetadata)
		} else {
			suite.Require().Error(err)
		}
	}
}
