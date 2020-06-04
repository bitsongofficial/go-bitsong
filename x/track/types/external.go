package types

import (
	"encoding/json"
	"sort"
)

type Externals map[string]string

// KeyValue is a simple key/value representation of one field of a OptionalData.
type KeyValue struct {
	Key   string
	Value string
}

// MarshalAmino transforms the OptionalData to an array of key/value.
func (m Externals) MarshalAmino() ([]KeyValue, error) {
	fieldKeys := make([]string, len(m))
	i := 0
	for key := range m {
		fieldKeys[i] = key
		i++
	}

	sort.Stable(sort.StringSlice(fieldKeys))

	p := make([]KeyValue, len(m))
	for i, key := range fieldKeys {
		p[i] = KeyValue{
			Key:   key,
			Value: m[key],
		}
	}

	return p, nil
}

// UnmarshalAmino transforms the key/value array to a Externals.
func (m *Externals) UnmarshalAmino(keyValues []KeyValue) error {
	tempMap := make(map[string]string, len(keyValues))
	for _, p := range keyValues {
		tempMap[p.Key] = p.Value
	}

	*m = tempMap

	return nil
}

// MarshalJSON implements encode.Marshaler
func (m Externals) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string(m))
}

// UnmarshalJSON implements decode.Unmarshaler
func (m *Externals) UnmarshalJSON(data []byte) error {
	var value map[string]string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*m = value
	return nil
}
