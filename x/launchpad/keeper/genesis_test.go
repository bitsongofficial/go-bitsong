package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestInitExportGenesis() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	genState := suite.app.LaunchPadKeeper.ExportGenesis(suite.ctx)
	suite.Require().Equal(genState, types.DefaultGenesisState())

	genState.MintableMetadataIds = []types.MintableMetadataIds{
		{
			CollectionId:        1,
			MintableMetadataIds: []uint64{1, 2, 3, 4, 5},
		},
	}
	genState.Launchpads = []types.LaunchPad{
		{
			CollId:               1,
			Price:                0,
			Treasury:             addr1.String(),
			Denom:                "ubtsg",
			GoLiveDate:           1659870342,
			EndTimestamp:         0,
			MaxMint:              5,
			Minted:               0,
			Authority:            addr1.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		},
	}

	suite.app.LaunchPadKeeper.InitGenesis(suite.ctx, *genState)
	savedGenState := suite.app.LaunchPadKeeper.ExportGenesis(suite.ctx)
	suite.Require().Equal(*genState, *savedGenState)
}
