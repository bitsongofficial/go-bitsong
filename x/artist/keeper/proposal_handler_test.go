package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func testProposal(title string, description string, artistID uint64) types.ArtistVerifyProposal {
	return types.ArtistVerifyProposal{
		Title:       title,
		Description: description,
		ArtistID:    artistID,
	}
}

func TestProposalHandlerPassed(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInputDefault(t, false, 10)
	owner := delAddr1

	// fix: set initialid
	keeper.SetArtistID(ctx, 1)

	// create artist
	_, err := keeper.CreateArtist(ctx, "Test", owner)
	require.NoError(t, err)

	// open proposal
	tp := testProposal("Verify my artist", "Please verify my artist", 1)
	hdlr := NewArtistVerifyProposalHandler(keeper)
	require.NoError(t, hdlr(ctx, tp))

	/*

		// add coins to the module account
		macc := keeper.GetDistributionAccount(ctx)
		err := macc.SetCoins(macc.GetCoins().Add(amount))
		require.NoError(t, err)

		supplyKeeper.SetModuleAccount(ctx, macc)

		account := accountKeeper.NewAccountWithAddress(ctx, recipient)
		require.True(t, account.GetCoins().IsZero())
		accountKeeper.SetAccount(ctx, account)

		feePool := keeper.GetFeePool(ctx)
		feePool.CommunityPool = sdk.NewDecCoins(amount)
		keeper.SetFeePool(ctx, feePool)

		tp := testProposal(recipient, amount)
		hdlr := NewCommunityPoolSpendProposalHandler(keeper)
		require.NoError(t, hdlr(ctx, tp))
		require.Equal(t, accountKeeper.GetAccount(ctx, recipient).GetCoins(), amount)*/
}
