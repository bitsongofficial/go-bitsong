package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestProcessRoyalties() {
	owner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	owner2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	creator2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	coins := sdk.Coins{sdk.NewInt64Coin("ubtsg", 1000000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, owner, coins)

	tests := []struct {
		testCase     string
		metadata     nfttypes.Metadata
		distributor  sdk.AccAddress
		distrAmount  uint64
		expectPass   bool
		checkAccount sdk.AccAddress
		checkAmount  uint64
	}{
		{
			"distribution to single creator",
			nfttypes.Metadata{
				CollId:               1,
				Id:                   1,
				MetadataAuthority:    owner.String(),
				MintAuthority:        owner.String(),
				Name:                 "NewPUNK1",
				SellerFeeBasisPoints: 10,
				Creators: []types.Creator{
					{Address: creator1.String(), Verified: true, Share: 100},
				},
				PrimarySaleHappened: true,
				IsMutable:           false,
			},
			owner,
			10000,
			true,
			creator1,
			1000,
		},
		{
			"distribution to multiple creators",
			nfttypes.Metadata{
				CollId:               1,
				Id:                   2,
				MetadataAuthority:    owner.String(),
				MintAuthority:        owner.String(),
				Name:                 "NewPUNK2",
				SellerFeeBasisPoints: 10,
				Creators: []types.Creator{
					{Address: creator1.String(), Verified: true, Share: 100},
					{Address: creator2.String(), Verified: true, Share: 100},
				},
				PrimarySaleHappened: true,
				IsMutable:           false,
			},
			owner,
			10000,
			true,
			creator1,
			500,
		},
		{
			"dust amount distribution cheeck",
			nfttypes.Metadata{
				CollId:               1,
				Id:                   2,
				MetadataAuthority:    owner.String(),
				MintAuthority:        owner.String(),
				Name:                 "NewPUNK2",
				SellerFeeBasisPoints: 10,
				Creators: []types.Creator{
					{Address: creator1.String(), Verified: true, Share: 1000000},
					{Address: creator2.String(), Verified: true, Share: 1},
				},
				PrimarySaleHappened: true,
				IsMutable:           false,
			},
			owner,
			1,
			true,
			creator1,
			0,
		},
		{
			"distribution to no creator",
			nfttypes.Metadata{
				CollId:               1,
				Id:                   3,
				MetadataAuthority:    owner.String(),
				MintAuthority:        owner.String(),
				Name:                 "NewPUNK3",
				SellerFeeBasisPoints: 10,
				Creators:             []types.Creator{},
				PrimarySaleHappened:  true,
				IsMutable:            false,
			},
			owner,
			10000,
			true,
			creator1,
			0,
		},
		{
			"distribution to invalid metadata - SellerFeeBasisPoints",
			nfttypes.Metadata{
				CollId:               1,
				Id:                   3,
				MetadataAuthority:    owner.String(),
				MintAuthority:        owner.String(),
				Name:                 "NewPUNK3",
				SellerFeeBasisPoints: 1000,
				Creators:             []types.Creator{},
				PrimarySaleHappened:  true,
				IsMutable:            false,
			},
			owner,
			10000,
			false,
			creator1,
			0,
		},
		{
			"distribution by zero balance account",
			nfttypes.Metadata{
				CollId:               1,
				Id:                   3,
				MetadataAuthority:    owner.String(),
				MintAuthority:        owner.String(),
				Name:                 "NewPUNK3",
				SellerFeeBasisPoints: 1000,
				Creators: []types.Creator{
					{Address: creator1.String(), Verified: true, Share: 100},
				},
				PrimarySaleHappened: true,
				IsMutable:           false,
			},
			owner2,
			10000,
			false,
			creator1,
			0,
		},
	}

	for _, tc := range tests {
		creatorOldBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.checkAccount, "ubtsg")
		err := suite.app.MarketplaceKeeper.ProcessRoyalties(suite.ctx, tc.metadata, owner, "ubtsg", tc.distrAmount)
		if tc.expectPass {
			suite.Require().NoError(err)
			creatorNewBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.checkAccount, "ubtsg")
			suite.Require().Equal(creatorOldBalance.Amount.Int64()+int64(tc.checkAmount), creatorNewBalance.Amount.Int64())
		} else {
			suite.Require().Error(err)
		}
	}
}
