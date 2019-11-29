package types

import (
	"encoding/json"
	"fmt"
)

type (
	// AlbumType is a type alias that represents an album type as a byte
	AlbumType byte
)

//nolint
const (
	TypeAlbum       AlbumType = 0x00
	TypeSingle      AlbumType = 0x01
	TypeCompilation AlbumType = 0x02
)

// AlbumTypeFromString turns a string into a AlbumType
func AlbumTypeFromString(str string) (AlbumType, error) {
	switch str {
	case "Album":
		return TypeAlbum, nil
	case "Single":
		return TypeSingle, nil
	case "Compilation":
		return TypeCompilation, nil
	default:
		return AlbumType(0xff), fmt.Errorf("'%s' is not a valid artist status", str)
	}
}

// Valid tells if the album type can be used
func (at AlbumType) Valid() bool {
	if at == TypeAlbum ||
		at == TypeSingle ||
		at == TypeCompilation {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (at AlbumType) Marshal() ([]byte, error) {
	return []byte{byte(at)}, nil
}

// Unmarshal needed for protobuf compatibility
func (at *AlbumType) Unmarshal(data []byte) error {
	*at = AlbumType(data[0])
	return nil
}

// Marshals to JSON using string
func (at AlbumType) MarshalJSON() ([]byte, error) {
	return json.Marshal(at.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (at *AlbumType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := AlbumTypeFromString(s)
	if err != nil {
		return err
	}

	*at = bz2
	return nil
}

// String implements the Stringer interface.
func (at AlbumType) String() string {
	switch at {
	case TypeAlbum:
		return "Album"
	case TypeSingle:
		return "Single"
	case TypeCompilation:
		return "Compilation"
	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (at AlbumType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(at.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(at))))
	}
}
