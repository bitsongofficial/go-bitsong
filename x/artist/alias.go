package artist

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/keeper"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
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
	NewArtist = types.NewArtist
)

type (
	// Keeper
	Keeper = keeper.Keeper

	// Types
	ArtistStatus = types.ArtistStatus
	Artist       = types.Artist
	Artists      = types.Artists

	Deposits      = types.Deposits
	DepositParams = types.DepositParams

	// Msgs
	MsgCreateArtist = types.MsgCreateArtist
)
