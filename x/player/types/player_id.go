package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

const (
	PlayerIDPrefix = "player"
)

type ID interface {
	String() string
	Uint64() uint64
	Bytes() []byte
	Prefix() string
	IsEqual(ID) bool
	MarshalJSON() ([]byte, error)
}

var (
	_ ID = PlayerID{}
)

type PlayerID []byte

func NewPlayerID(i uint64) PlayerID {
	return types.Uint64ToBigEndian(i)
}

func NewPlayerIDFromString(s string) (PlayerID, error) {
	if len(s) < 5 {
		return nil, fmt.Errorf("invalid player id length")
	}

	i, err := strconv.ParseUint(s[4:], 16, 64)
	if err != nil {
		return nil, err
	}

	return NewPlayerID(i), nil
}

func (id PlayerID) String() string {
	return fmt.Sprintf("%s%x", PlayerIDPrefix, id.Uint64())
}

func (id PlayerID) Uint64() uint64 {
	return binary.BigEndian.Uint64(id)
}

func (id PlayerID) Bytes() []byte {
	return id
}

func (id PlayerID) Prefix() string {
	return PlayerIDPrefix
}

func (id PlayerID) IsEqual(_id ID) bool {
	return id.String() == _id.String()
}

func (id PlayerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PlayerID) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	_id, err := NewPlayerIDFromString(s)
	if err != nil {
		return err
	}

	*id = _id

	return nil
}
