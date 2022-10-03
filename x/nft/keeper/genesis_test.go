package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creators := []types.Creator{
		{
			Address:  addr1.String(),
			Verified: false,
			Share:    1,
		},
	}

	genesisState := types.GenesisState{
		Params: types.Params{
			IssuePrice: sdk.NewInt64Coin("ubtsg", 1000),
		},
		Metadata: []types.Metadata{
			{
				CollId:               1,
				Id:                   1,
				MetadataAuthority:    addr1.String(),
				MintAuthority:        addr1.String(),
				Name:                 "meta1",
				Uri:                  "uri1",
				SellerFeeBasisPoints: 10,
				Creators:             creators,
				PrimarySaleHappened:  false,
				IsMutable:            true,
			},
		},
		LastMetadataIds: []types.LastMetadataIdInfo{
			{
				CollId:         1,
				LastMetadataId: 1,
			},
		},
		Nfts: []types.NFT{
			{
				CollId:     1,
				MetadataId: 1,
				Seq:        0,
				Owner:      addr1.String(),
			},
		},
		Collections: []types.Collection{
			{
				Id:              1,
				Name:            "name1",
				Uri:             "uri1",
				UpdateAuthority: addr1.String(),
			},
		},
		LastCollectionId: 1,
	}

	suite.app.NFTKeeper.InitGenesis(suite.ctx, genesisState)

	exportedGenesis := suite.app.NFTKeeper.ExportGenesis(suite.ctx)
	suite.Require().Equal(genesisState, *exportedGenesis)
}
