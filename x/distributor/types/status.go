package types

import (
	"encoding/json"
	"fmt"
)

type (
	DistributorStatus byte
)

//nolint
const (
	StatusNil      DistributorStatus = 0x00
	StatusVerified DistributorStatus = 0x01
	StatusRejected DistributorStatus = 0x02
	StatusFailed   DistributorStatus = 0x03
)

func DistributorStatusFromString(str string) (DistributorStatus, error) {
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
		return DistributorStatus(0xff), fmt.Errorf("'%s' is not a valid distributor status", str)
	}
}

func (status DistributorStatus) Valid() bool {
	if status == StatusNil ||
		status == StatusVerified ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status DistributorStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *DistributorStatus) Unmarshal(data []byte) error {
	*status = DistributorStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status DistributorStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *DistributorStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := DistributorStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status DistributorStatus) String() string {
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
func (status DistributorStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
