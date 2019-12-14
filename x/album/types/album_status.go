package types

import (
	"encoding/json"
	"fmt"
)

type (
	// AlbumStatus is a type alias that represents an album status as a byte
	AlbumStatus byte
)

//nolint
const (
	StatusNil           AlbumStatus = 0x00
	StatusDepositPeriod AlbumStatus = 0x01
	StatusVerified      AlbumStatus = 0x02
	StatusRejected      AlbumStatus = 0x03
	StatusFailed        AlbumStatus = 0x04
)

// AlbumStatusFromString turns a string into a AlbumStatus
func AlbumStatusFromString(str string) (AlbumStatus, error) {
	switch str {
	case "DepositPeriod":
		return StatusDepositPeriod, nil

	case "Verified":
		return StatusVerified, nil

	case "Rejected":
		return StatusRejected, nil

	case "Failed":
		return StatusFailed, nil

	case "":
		return StatusNil, nil

	default:
		return AlbumStatus(0xff), fmt.Errorf("'%s' is not a valid album status", str)
	}
}

// Valid tells if the album status can be used
func (status AlbumStatus) Valid() bool {
	if status == StatusNil ||
		status == StatusDepositPeriod ||
		status == StatusVerified ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status AlbumStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *AlbumStatus) Unmarshal(data []byte) error {
	*status = AlbumStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status AlbumStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *AlbumStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := AlbumStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status AlbumStatus) String() string {
	switch status {
	case StatusDepositPeriod:
		return "DepositPeriod"

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
func (status AlbumStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
