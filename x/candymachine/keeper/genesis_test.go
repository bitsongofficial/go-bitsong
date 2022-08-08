package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// TODO: test InitGenesis
// TODO: test ExportGenesis

func (suite *KeeperTestSuite) TestInitExportGenesis() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	genState := suite.app.CandyMachineKeeper.ExportGenesis(suite.ctx)
	suite.Require().Equal(genState, types.DefaultGenesisState())

	genState.Candymachines = []types.CandyMachine{
		{
			CollId:     1,
			Price:      0,
			Treasury:   addr1.String(),
			Denom:      "ubtsg",
			GoLiveDate: 1659870342,
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Mint,
				Value:   1000,
			},
			Minted:               0,
			Authority:            addr1.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		},
	}

	suite.app.CandyMachineKeeper.InitGenesis(suite.ctx, *genState)
	savedGenState := suite.app.CandyMachineKeeper.ExportGenesis(suite.ctx)
	suite.Require().Equal(*genState, *savedGenState)
}
