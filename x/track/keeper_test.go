package track

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"strconv"

	"github.com/BitSongOfficial/go-bitsong/x/track/types"

	"testing"
)

var (
	Addrs = createTestAddrs(500)
	PKs   = createTestPubKeys(500)

	addrDels = []sdk.AccAddress{
		Addrs[0],
		Addrs[1],
	}
	addrVals = []sdk.ValAddress{
		sdk.ValAddress(Addrs[2]),
		sdk.ValAddress(Addrs[3]),
		sdk.ValAddress(Addrs[4]),
		sdk.ValAddress(Addrs[5]),
		sdk.ValAddress(Addrs[6]),
	}
)

type TestInput struct {
	cdc            *codec.Codec
	ctx            sdk.Context
	trackKeeper    Keeper
	stackingKeeper staking.Keeper
}

func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)

	// Register AppAccount
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/staking/BaseAccount", nil)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func SetupTestInput(t *testing.T) TestInput {
	trackCapKey := sdk.NewKVStoreKey("trackCapKey")
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	stakingCapKey := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(stakingtypes.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(trackCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(stakingCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-ID"}, false, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := MakeTestCodec()

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(stakingtypes.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	bk := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:          nil,
		stakingtypes.NotBondedPoolName: {supply.Burner, supply.Staking},
		stakingtypes.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, stakingCapKey, accountKeeper, bk, maccPerms)

	initTokens := sdk.TokensFromConsensusPower(10)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(int64(len(Addrs)))))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	stakingKeeper := staking.NewKeeper(cdc, stakingCapKey, tkeyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace), stakingtypes.DefaultCodespace)
	stakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())

	// set module accounts
	err = notBondedPool.SetCoins(totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range Addrs {
		_, err := bk.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	trackSubspace := pk.Subspace(types.DefaultParamspace)

	trackKeeper := NewKeeper(trackCapKey, cdc, trackSubspace, stakingKeeper, supplyKeeper)
	trackKeeper.SetParams(ctx, types.DefaultParams())
	trackKeeper.SetInitialTrackID(ctx, types.DefaultStartingTrackID)
	initialFeePool := Pool{
		Rewards: sdk.NewCoins(sdk.NewInt64Coin("ubtsg", 100000)),
	}
	trackKeeper.SetFeePlayPool(ctx, initialFeePool)

	// Create validator
	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var validators [3]stakingtypes.Validator
	for i, amt := range amts {
		validators[i] = stakingtypes.NewValidator(addrVals[i], PKs[i], stakingtypes.Description{})
		validators[i], _ = validators[i].AddTokensFromDel(amt)
		stakingKeeper.SetValidator(ctx, validators[i])
		stakingKeeper.SetValidatorByPowerIndex(ctx, validators[i])
		stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	// Add delegation
	delegation := stakingtypes.NewDelegation(addrDels[0], addrVals[0], sdk.NewDec(9))
	stakingKeeper.SetDelegation(ctx, delegation)
	delegation2 := stakingtypes.NewDelegation(addrDels[0], addrVals[1], sdk.NewDec(10))
	stakingKeeper.SetDelegation(ctx, delegation2)
	delegation3 := stakingtypes.NewDelegation(addrDels[1], addrVals[0], sdk.NewDec(20))
	stakingKeeper.SetDelegation(ctx, delegation3)

	return TestInput{cdc: cdc, ctx: ctx, trackKeeper: trackKeeper}
}

func TestPublishTrack(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	_, err := trackKeeper.PublishTrack(ctx, "My second track on BitSong feat. Angelo Recca", addrDels[0], "ipfs hash", sdk.NewDecWithPrec(10, 2))
	require.NoError(t, err)

	_, ok := trackKeeper.GetTrack(ctx, 1)
	require.True(t, ok)

	_, err = trackKeeper.PublishTrack(ctx, "My third track on BitSong feat. Angelo Recca", addrDels[0], "ipfs hash", sdk.NewDecWithPrec(7, 2))
	require.NoError(t, err)

	_, ok = trackKeeper.GetTrack(ctx, 2)
	require.True(t, ok)
}

func TestGetSetTrack(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	newTrack := types.Track{
		TrackID:                 1,
		Owner:                   addrDels[0],
		Title:                   "My first track on BitSong feat. Angelo Recca",
		Content:                 "ipfs hash",
		TotalReward:             sdk.NewInt(0),
		RedistributionSplitRate: sdk.NewDecWithPrec(5, 2),
		CreateTime:              ctx.BlockHeader().Time,
	}

	err := trackKeeper.SetTrack(ctx, newTrack)
	require.NoError(t, err)

	_, ok := trackKeeper.GetTrack(ctx, 1)
	require.True(t, ok)

	_, ok = trackKeeper.GetTrack(ctx, 2)
	require.False(t, ok)
}

func TestGetSetPlay(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	userPower := trackKeeper.GetAccPower(ctx, addrDels[0])

	newPlay := types.Play{
		AccAddress: addrDels[0],
		TrackId:    1,
		Shares:     userPower,
		Streams:    sdk.NewInt(1),
		CreateTime: ctx.BlockHeader().Time,
	}

	err := trackKeeper.SetPlay(ctx, newPlay)
	require.NoError(t, err)

	play, ok := trackKeeper.GetAccPlay(ctx, addrDels[0], 1)
	require.True(t, ok)
	fmt.Printf("%s", play)

	_, ok = trackKeeper.GetAccPlay(ctx, addrDels[1], 1)
	require.False(t, ok)
}

func TestSavePlay(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	play, ok := trackKeeper.PlayTrack(ctx, addrDels[0], 1)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(1))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[0], 1)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(2))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[1], 1)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(1))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[1], 1)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(2))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[0], 2)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(1))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[1], 2)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(1))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[1], 2)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(2))
}

func TestPlayIterator(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	play, ok := trackKeeper.PlayTrack(ctx, addrDels[0], 1)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(1))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[0], 1)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(2))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[0], 2)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(1))

	play, ok = trackKeeper.PlayTrack(ctx, addrDels[1], 1)
	require.True(t, ok)
	require.Equal(t, play.Streams, sdk.NewInt(1))

	plays := trackKeeper.GetAllPlays(ctx)

	fmt.Printf("%s", plays)
}

func TestAllocateTokensToAccount(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	tokens := sdk.DecCoins{{sdk.DefaultBondDenom, sdk.NewDec(10000)}}
	trackKeeper.AllocateTokensToAccount(ctx, addrDels[0], tokens)
}

func createTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, TesterAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

func createTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") //base pubkey string
		buffer.WriteString(numString)                                                       //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}

func TesterAddr(addr string, bech string) sdk.AccAddress {

	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res

}
