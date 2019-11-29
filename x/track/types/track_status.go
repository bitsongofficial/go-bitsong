package types

import (
	"encoding/json"
	"fmt"
)

type (
	// TrackStatus is a type alias that represents an track status as a byte
	TrackStatus byte
)

//nolint
const (
	StatusNil      TrackStatus = 0x00
	StatusVerified TrackStatus = 0x01
	StatusRejected TrackStatus = 0x02
	StatusFailed   TrackStatus = 0x03
)

// TrackStatusFromString turns a string into a AlbumStatus
func TrackStatusFromString(str string) (TrackStatus, error) {
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
		return TrackStatus(0xff), fmt.Errorf("'%s' is not a valid track status", str)
	}
}

// Valid tells if the track status can be used
func (status TrackStatus) Valid() bool {
	if status == StatusNil ||
		status == StatusVerified ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status TrackStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *TrackStatus) Unmarshal(data []byte) error {
	*status = TrackStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status TrackStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *TrackStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := TrackStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status TrackStatus) String() string {
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
func (status TrackStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
