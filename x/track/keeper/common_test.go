package keeper

import (
	"encoding/hex"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
	"time"
)

var (
	mockTitle          = "The Show Must Go On"
	mockEmptyTitle     = ""
	mockEmptyAttribute = map[string]string{}
	mockContent        = types.Content{
		Path:        "/ipfs/Qm....",
		ContentType: "",
		Duration:    0,
		Attributes:  nil,
	}
	mockEmptyContent      = types.Content{}
	mockTrackMediaNoVideo = types.TrackMedia{
		Audio: mockContent,
		Video: mockEmptyContent,
		Image: mockContent,
	}
	mockEmptyTrackMedia     = types.TrackMedia{}
	mockRightHolder1        = types.NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder1"))), 100)
	mockRightHolder2        = types.NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder2"))), 25)
	mockRightHolder3        = types.NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder3"))), 25)
	mockRightHolder4        = types.NewRightHolder(sdk.AccAddress(crypto.AddressHash([]byte("rightHolder4"))), 50)
	mockRightsHoldersSingle = types.RightsHolders{
		mockRightHolder1,
	}
	mockRightsHoldersMultiple = types.RightsHolders{
		mockRightHolder2,
		mockRightHolder3,
		mockRightHolder4,
	}
	mockRewards = types.TrackRewards{
		Users:     10,
		Playlists: 10,
	}
	mockOwner         = sdk.AccAddress(crypto.AddressHash([]byte("owner")))
	mockHexDecoded, _ = hex.DecodeString("B0FA2953B126722264F67828AF7443144C85D867")
	mockTrackAddr1    = crypto.Address(mockHexDecoded)
	mockTrack         = types.Track{
		Title:         mockTitle,
		Attributes:    mockEmptyAttribute,
		Media:         mockTrackMediaNoVideo,
		Rewards:       mockRewards,
		RightsHolders: mockRightsHoldersSingle,
		SubmitTime:    time.Time{},
		Owner:         nil,
	}
)

func SetupTestInput() (sdk.Context, Keeper) {
	// define store keys
	trackKey := sdk.NewKVStoreKey("track")

	// create an in-memory db
	memDB := db.NewMemDB()
	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(trackKey, sdk.StoreTypeIAVL, memDB)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	// create a Cdc and a context
	cdc := testCodec()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return ctx, NewKeeper(cdc, trackKey)
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	// register the different types
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	types.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}
