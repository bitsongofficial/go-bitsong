package album

import (
	"github.com/bitsongofficial/go-bitsong/x/album/client"
	"github.com/bitsongofficial/go-bitsong/x/album/keeper"
	"github.com/bitsongofficial/go-bitsong/x/album/types"
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
	NewAlbum                      = types.NewAlbum
	NewAlbumVerifyProposal        = types.NewAlbumVerifyProposal
	NewAlbumVerifyProposalHandler = keeper.NewAlbumVerifyProposalHandler
	ProposalHandler               = client.ProposalHandler

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

	// Msgs
	MsgCreateAlbum = types.MsgCreateAlbum
)
