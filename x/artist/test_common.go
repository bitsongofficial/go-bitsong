package artist

import (
	"bytes"
	"log"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/mock"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

var (
	valTokens  = sdk.TokensFromConsensusPower(42)
	initTokens = sdk.TokensFromConsensusPower(100000)
	valCoins   = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

type testInput struct {
	mApp     *mock.App
	keeper   Keeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

func getMockApp(t *testing.T, numGenAccs int, genState GenesisState, genAccs []auth.Account) testInput {
	mApp := mock.NewApp()

	types.RegisterCodec(mApp.Cdc)

	keyGov := sdk.NewKVStoreKey(types.StoreKey)

	pk := mApp.ParamsKeeper

	/*rtr := NewRouter().
	AddRoute(RouterKey, ProposalHandler)*/
	keeper := NewKeeper(mApp.Cdc, keyGov, pk, types.DefaultCodespace)
	//keeper := NewKeeper(mApp.Cdc, keyGov, pk, pk.Subspace(DefaultParamspace), supplyKeeper, sk, DefaultCodespace, rtr)

	mApp.Router().AddRoute(types.RouterKey, NewHandler(keeper))
	mApp.QueryRouter().AddRoute(types.QuerierRoute, NewQuerier(keeper))

	mApp.SetInitChainer(getInitChainer(mApp, keeper, genState))

	require.NoError(t, mApp.CompleteSetup(keyGov))

	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs, valCoins)
	}

	mock.SetGenesis(mApp, genAccs)

	return testInput{mApp, keeper, addrs, pubKeys, privKeys}
}

// artist initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, genState GenesisState) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		if genState.IsEmpty() {
			InitGenesis(ctx, keeper, DefaultGenesisState())
		} else {
			InitGenesis(ctx, keeper, genState)
		}
		return abci.ResponseInitChain{}
	}
}

func SortAddresses(addrs []sdk.AccAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}
	SortByteArrays(byteAddrs)
	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// Public
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}

func testArtist() types.Meta {
	return types.NewGeneralMeta("Freddy Mercury")
}

// checks if two artists are equal (note: slow, for tests only)
func ArtistEqual(artistA types.Artist, artistB types.Artist) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(artistA),
		types.ModuleCdc.MustMarshalBinaryBare(artistB))
}
