package artist

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/client"
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
	NewArtistVerifyProposalHandler = keeper.NewArtistVerifyProposalHandler
	HandleVerifyArtistProposal     = keeper.HandleVerifyArtistProposal
	ProposalHandler                = client.ProposalHandler
)

type (
	// Keeper
	Keeper = keeper.Keeper

	// Types
	ArtistStatus         = types.ArtistStatus
	Artist               = types.Artist
	Artists              = types.Artists
	ArtistVerifyProposal = types.ArtistVerifyProposal

	// Msgs
	MsgCreateArtist = types.MsgCreateArtist
)
