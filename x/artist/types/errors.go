package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const (
	ErrCodeArtistNotFound    = 1
	ErrCodeArtistCreateError = 2
)

var (
	DefaultCodespace = ModuleName

	ErrArtistNotFound     = sdkerrors.Register(ModuleName, ErrCodeArtistNotFound, "artist not found")
	ErrArtistiCreateError = sdkerrors.Register(ModuleName, ErrCodeArtistCreateError, "artist create error")
)
