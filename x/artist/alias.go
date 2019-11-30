package artist

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/keeper"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
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

	// Msgs
	NewMsgCreateArtist             = types.NewMsgCreateArtist
	NewArtistVerifyProposalHandler = keeper.NewArtistVerifyProposalHandler
)

type (
	// Keeper
	Keeper = keeper.Keeper

	// Types
	ArtistStatus = types.ArtistStatus
	Artist       = types.Artist
	Artists      = types.Artists

	// Msgs
	MsgCreateArtist = types.MsgCreateArtist
)
