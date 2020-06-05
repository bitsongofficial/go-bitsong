package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the track store
type Keeper struct {
	bankKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}

// NewKeeper creates a content keeper
func NewKeeper(bk bank.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		bankKeeper: bk,
		storeKey:   key,
		cdc:        cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetTrack(ctx sdk.Context, c string) (track types.Track, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTrackKey(c))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &track)
	return track, true
}

func (k Keeper) SetTrack(ctx sdk.Context, track *types.Track) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&track)
	store.Set(types.GetTrackKey(track.Cid), bz)
}

func (k Keeper) IterateTracks(ctx sdk.Context, fn func(track types.Track) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TrackKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var track types.Track
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &track)
		if fn(track) {
			break
		}
	}
}

func (k Keeper) GetTracks(ctx sdk.Context) []types.Track {
	var tracks []types.Track
	k.IterateTracks(ctx, func(track types.Track) (stop bool) {
		tracks = append(tracks, track)
		return false
	})
	return tracks
}

func (k Keeper) Add(ctx sdk.Context, track *types.Track) (string, error) {
	_, found := k.GetTrack(ctx, track.Cid)
	if found {
		return "", fmt.Errorf("track %s is already stored", track.Cid)
	}

	pref := cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}

	cid, err := pref.Sum([]byte(track.Title)) // TODO: add more data
	if err != nil {
		return "", err
	}
	track.Cid = cid.String()

	trackID := k.autoIncrementID(ctx, types.KeyLastTrackID)
	track.TrackID = trackID

	trackAddr := k.generateTrackAddress(ctx, trackID)
	track.Address = trackAddr

	trackUri := k.generateTrackUri(ctx, trackID)
	track.Uri = trackUri

	//content.CreatedAt = ctx.BlockHeader().Time
	k.SetTrack(ctx, track)

	return cid.String(), nil
}

func (k Keeper) autoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(lastIDKey)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(lastIDKey, bz)
	return id
}

func (k Keeper) generateTrackAddress(ctx sdk.Context, trackID uint64) crypto.Address {
	addr := make([]byte, 20)
	addr[0] = 'T'
	binary.PutUvarint(addr[1:], trackID<<32)
	//return sdk.AccAddress(crypto.AddressHash(addr))
	return crypto.AddressHash(addr)
}

func (k Keeper) generateTrackUri(ctx sdk.Context, trackID uint64) string {
	return fmt.Sprintf("bitsong:track:%d", trackID)
}

/*func allocateDaoFunds(cnt types.Content, coin sdk.Coin) types.Content {
	// allocate dao funds
	for i, rh := range cnt.RightsHolders {
		price := sdk.NewDecCoinFromCoin(coin)
		allocation := price.Amount.Quo(sdk.NewDec(100).Quo(rh.Quota))
		cnt.RightsHolders[i].Rewards = rh.Rewards.Add(sdk.NewDecCoinFromDec(btsg.BondDenom, allocation))
	}

	return cnt
}*/
