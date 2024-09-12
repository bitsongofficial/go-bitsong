package merkledrop_test

import (
	"fmt"
	"testing"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/keeper"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func generateMerkledrop(owner sdk.AccAddress, accs map[string]string) (*types.MsgCreate, map[string]cli.ClaimInfo) {
	startHeight := int64(0)
	endHeight := int64(100)

	accMap, _ := cli.AccountsFromMap(accs)
	tree, claimInfo, totalAmt, _ := cli.CreateDistributionList(accMap)
	merkleRoot := fmt.Sprintf("%x", tree.Root())

	airdropCoins, _ := sdk.ParseCoinNormalized(fmt.Sprintf("%s%s", totalAmt.String(), sdk.DefaultBondDenom))

	return types.NewMsgCreate(
		owner,
		merkleRoot,
		startHeight,
		endHeight,
		airdropCoins,
	), claimInfo
}

func TestExpiredMerkledropHeight(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	msgSvr := keeper.NewMsgServerImpl(app.MerkledropKeeper)

	owner1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	owner2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	owner3 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	initCoins := sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000)),
	}
	accs := make(map[string]string, 2)
	accs[owner1.String()] = "10000"
	accs[owner2.String()] = "20000"
	accs[owner3.String()] = "30000"

	// mint coins for owner1
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	require.NoError(t, err)

	// send coins to owner1
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, owner1, initCoins)
	require.NoError(t, err)

	// mint coins for owner2
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	require.NoError(t, err)

	// send coins to owner2
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, owner2, initCoins)
	require.NoError(t, err)

	// mint coins for owner3
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	require.NoError(t, err)

	// send coins to owner3
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, owner3, initCoins)
	require.NoError(t, err)

	// merkledrops len should be zero
	merkledrops := app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 0)

	// owner1's balance should be initCoins
	balance := app.BankKeeper.GetBalance(ctx, owner1, sdk.DefaultBondDenom)
	require.Equal(t, initCoins.AmountOf(sdk.DefaultBondDenom), balance.Amount)

	// owner2's balance should be initCoins
	balance = app.BankKeeper.GetBalance(ctx, owner1, sdk.DefaultBondDenom)
	require.Equal(t, initCoins.AmountOf(sdk.DefaultBondDenom), balance.Amount)

	// owner3's balance should be initCoins
	balance = app.BankKeeper.GetBalance(ctx, owner1, sdk.DefaultBondDenom)
	require.Equal(t, initCoins.AmountOf(sdk.DefaultBondDenom), balance.Amount)

	// owner1 create a merkledrop
	msg1, claimInfo1 := generateMerkledrop(owner1, accs)
	res, err := msgSvr.Create(sdk.WrapSDKContext(ctx), msg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// owner2 create a merkledrop
	msg2, _ := generateMerkledrop(owner2, accs)
	res, err = msgSvr.Create(sdk.WrapSDKContext(ctx), msg2)
	require.NoError(t, err)
	require.NotNil(t, res)

	// owner3 create a merkledrop
	msg3, _ := generateMerkledrop(owner3, accs)
	res, err = msgSvr.Create(sdk.WrapSDKContext(ctx), msg3)
	require.NoError(t, err)
	require.NotNil(t, res)

	// check merkledrops length, should be 3
	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 3)

	// move to block 50
	merkledrop.EndBlocker(ctx, app.MerkledropKeeper)

	newHeader := ctx.BlockHeader()
	newHeader.Height = 50
	ctx = ctx.WithBlockHeader(newHeader)
	merkledrop.EndBlocker(ctx, app.MerkledropKeeper)

	// claim all balance from merkledrop owner1
	amt, ok := sdk.NewIntFromString(claimInfo1[owner1.String()].Amount)
	require.True(t, ok)

	claim := types.NewMsgClaim(
		claimInfo1[owner1.String()].Index, uint64(1), amt, claimInfo1[owner1.String()].Proof, owner1,
	)

	resClaim, err := msgSvr.Claim(sdk.WrapSDKContext(ctx), claim)
	require.NoError(t, err)
	require.NotNil(t, resClaim)

	amt, ok = sdk.NewIntFromString(claimInfo1[owner2.String()].Amount)
	require.True(t, ok)

	claim = types.NewMsgClaim(
		claimInfo1[owner2.String()].Index, uint64(1), amt, claimInfo1[owner2.String()].Proof, owner2,
	)

	resClaim, err = msgSvr.Claim(sdk.WrapSDKContext(ctx), claim)
	require.NoError(t, err)
	require.NotNil(t, resClaim)

	amt, ok = sdk.NewIntFromString(claimInfo1[owner3.String()].Amount)
	require.True(t, ok)

	claim = types.NewMsgClaim(
		claimInfo1[owner3.String()].Index, uint64(1), amt, claimInfo1[owner3.String()].Proof, owner3,
	)

	resClaim, err = msgSvr.Claim(sdk.WrapSDKContext(ctx), claim)
	require.NoError(t, err)
	require.NotNil(t, resClaim)

	// check merkledrops len, should be 2
	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 2)

	// move to block 100
	newHeader = ctx.BlockHeader()
	newHeader.Height = 100
	ctx = ctx.WithBlockHeader(newHeader)
	merkledrop.EndBlocker(ctx, app.MerkledropKeeper)

	// check merkledrops len, should be 0
	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 0)

	// move to block 110
	newHeader = ctx.BlockHeader()
	newHeader.Height = 110
	ctx = ctx.WithBlockHeader(newHeader)
	merkledrop.EndBlocker(ctx, app.MerkledropKeeper)

	// check merkledrops len, should be 0
	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 0)
}

func TestExpiredMerkledrop(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	msgSvr := keeper.NewMsgServerImpl(app.MerkledropKeeper)

	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	initCoins := sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000000)),
	}
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, owner, initCoins)
	require.NoError(t, err)

	merkledrops := app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 0)

	startHeight := int64(0)
	endHeight := int64(100)

	airdropCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))

	msg := types.NewMsgCreate(
		owner,
		"",
		startHeight,
		endHeight,
		airdropCoin,
	)

	balance := app.BankKeeper.GetBalance(ctx, owner, sdk.DefaultBondDenom)
	require.Equal(t, initCoins.AmountOf(sdk.DefaultBondDenom), balance.Amount)

	res, err := msgSvr.Create(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	balance = app.BankKeeper.GetBalance(ctx, owner, sdk.DefaultBondDenom)
	creationFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000_000))
	expectedAmt := initCoins.AmountOf(sdk.DefaultBondDenom).Sub(creationFee.Amount).Sub(airdropCoin.Amount)
	require.Equal(t, expectedAmt, balance.Amount)

	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 1)

	newHeader := ctx.BlockHeader()
	newHeader.Height = endHeight
	ctx = ctx.WithBlockHeader(newHeader)

	merkledrop.EndBlocker(ctx, app.MerkledropKeeper)

	merkledrops = app.MerkledropKeeper.GetAllMerkleDrops(ctx)
	require.Len(t, merkledrops, 0)

	balance = app.BankKeeper.GetBalance(ctx, owner, sdk.DefaultBondDenom)
	expectedAmt = initCoins.AmountOf(sdk.DefaultBondDenom).Sub(creationFee.Amount)
	require.Equal(t, expectedAmt, balance.Amount)
}
