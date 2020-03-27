package types

import (
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace string = ModuleName

	CodeInvalidGenesis          uint32 = 0
	CodeInvalidTrackTitle       uint32 = 1
	CodeUnknownTrack            uint32 = 2
	CodeInvalidTrackStatus      uint32 = 3
	CodeInvalidTrackMetadataURI uint32 = 4
)

func ErrInvalidGenesis(codespace string, msg string) error {
	return sdkerr.New(codespace, CodeInvalidGenesis, msg)
}

func ErrInvalidTrackTitle(codespace string, msg string) error {
	return sdkerr.New(codespace, CodeInvalidTrackTitle, msg)
}

func ErrUnknownTrack(codespace string, msg string) error {
	return sdkerr.New(codespace, CodeUnknownTrack, msg)
}

func ErrInvalidTrackStatus(codespace string, msg string) error {
	return sdkerr.New(codespace, CodeInvalidTrackStatus, msg)
}

func ErrInvalidTrackMetadataURI(codespace string, msg string) error {
	return sdkerr.New(codespace, CodeInvalidTrackMetadataURI, msg)
}
