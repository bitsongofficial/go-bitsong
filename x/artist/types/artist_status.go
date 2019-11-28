package types

import (
	"encoding/json"
	"fmt"
)

type (
	// ArtistStatus is a type alias that represents an artist status as a byte
	ArtistStatus byte
)

//nolint
const (
	StatusNil      ArtistStatus = 0x00
	StatusVerified ArtistStatus = 0x01
	StatusRejected ArtistStatus = 0x02
	StatusFailed   ArtistStatus = 0x03
)

// ArtistStatusFromString turns a string into a ArtistStatus
func ArtistStatusFromString(str string) (ArtistStatus, error) {
	switch str {
	case "Verified":
		return StatusVerified, nil

	case "Rejected":
		return StatusRejected, nil

	case "Failed":
		return StatusFailed, nil

	case "":
		return StatusNil, nil

	default:
		return ArtistStatus(0xff), fmt.Errorf("'%s' is not a valid artist status", str)
	}
}

// Valid tells if the artist status can be used
func (status ArtistStatus) Valid() bool {
	if status == StatusNil ||
		status == StatusVerified ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ArtistStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ArtistStatus) Unmarshal(data []byte) error {
	*status = ArtistStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ArtistStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ArtistStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ArtistStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status ArtistStatus) String() string {
	switch status {
	case StatusVerified:
		return "Verified"

	case StatusRejected:
		return "Rejected"

	case StatusFailed:
		return "Failed"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (status ArtistStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
