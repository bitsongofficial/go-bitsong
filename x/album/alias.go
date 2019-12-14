package album

import (
	"github.com/bitsongofficial/go-bitsong/x/album/keeper"
	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	DefaultParamspace = types.DefaultParamspace
)

var (
	// Keeper methods
	NewKeeper  = keeper.NewKeeper
	NewHandler = keeper.NewHandler
	NewQuerier = keeper.NewQuerier

	// Codec
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	// Types
	NewAlbum = types.NewAlbum

	// Msgs
	NewMsgCreateAlbum = types.NewMsgCreateAlbum
)

type (
	// Keeper
	Keeper = keeper.Keeper

	// Types
	AlbumStatus = types.AlbumStatus
	Album       = types.Album
	Albums      = types.Albums

	Deposits      = types.Deposits
	DepositParams = types.DepositParams

	// Msgs
	MsgCreateAlbum = types.MsgCreateAlbum
)
