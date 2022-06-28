package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func (suite *KeeperTestSuite) TestKeeper_GetAllIndexById() {
	suite.SetupTest()
	ctx := suite.Ctx
	mk := suite.App.MerkledropKeeper

	// set merkledrop
	merkledropID := uint64(1)
	index := uint64(0)
	denom := "ubtsg"
	owner := suite.TestAccs[0]

	merkledrop := types.Merkledrop{
		Id:          merkledropID,
		MerkleRoot:  "sdsd",
		StartHeight: int64(10),
		EndHeight:   int64(20),
		Denom:       denom,
		Amount:      sdk.NewInt(100),
		Claimed:     sdk.ZeroInt(),
		Owner:       owner.String(),
	}

	mk.SetMerkleDrop(ctx, merkledrop)

	merkledrop2 := types.Merkledrop{
		Id:          merkledropID + 1,
		MerkleRoot:  "sdsd",
		StartHeight: int64(10),
		EndHeight:   int64(20),
		Denom:       denom,
		Amount:      sdk.NewInt(100),
		Claimed:     sdk.ZeroInt(),
		Owner:       owner.String(),
	}
	mk.SetMerkleDrop(ctx, merkledrop2)

	// check isClaimed => should be false
	isClaimed := mk.IsClaimed(ctx, merkledropID, index)
	assert.False(suite.T(), isClaimed)

	// set isClaimed
	mk.SetClaimed(ctx, merkledropID, index)

	// check isClaimed => should be true
	isClaimed = mk.IsClaimed(ctx, merkledropID, index)
	assert.True(suite.T(), isClaimed)

	// set fake claimed
	mk.SetClaimed(ctx, merkledropID, 10)
	mk.SetClaimed(ctx, merkledropID, 28)

	// get all indexes by merkledrop id
	indexes := mk.GetAllIndexesByMerkledropID(ctx, merkledropID)
	assert.Equal(suite.T(), []uint64{0, 10, 28}, indexes)

	// set fake claimed
	mk.SetClaimed(ctx, merkledropID+1, 30)
	mk.SetClaimed(ctx, merkledropID+1, 40)

	// get all indexes
	allindexes := mk.GetAllIndexes(ctx)
	assert.Equal(suite.T(), 2, len(allindexes))
	assert.Equal(suite.T(), 3, len(allindexes[0].Index))
	assert.Equal(suite.T(), 2, len(allindexes[1].Index))
}
