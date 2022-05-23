package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (suite *KeeperTestSuite) TestKeeper_GetAllIndexById() {
	suite.SetupTest()

	// set merkledrop
	merkledrop := types.Merkledrop{
		Id:         1,
		MerkleRoot: "sdsd",
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		Coin: sdk.Coin{
			Denom:  "ubtsg",
			Amount: sdk.NewInt(100),
		},
		Claimed:   sdk.Coin{Denom: "ubtsg", Amount: sdk.ZeroInt()},
		Owner:     suite.TestAccs[0].String(),
		Withdrawn: false,
	}
	suite.App.MerkledropKeeper.SetMerkleDrop(suite.Ctx, merkledrop)

	// set isClaimed

	// get all claimed
}
