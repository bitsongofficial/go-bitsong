package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

type Attributes map[string]string

type KeyValue struct {
	Key   string
	Value string
}

func (a Attributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string(a))
}

func (a *Attributes) UnmarshalJSON(data []byte) error {
	var value map[string]string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*a = value
	return nil
}

func (a Attributes) MarshalAmino() ([]KeyValue, error) {
	fieldKeys := make([]string, len(a))
	i := 0
	for key := range a {
		fieldKeys[i] = key
		i++
	}

	sort.Stable(sort.StringSlice(fieldKeys))

	p := make([]KeyValue, len(a))
	for i, key := range fieldKeys {
		p[i] = KeyValue{
			Key:   key,
			Value: a[key],
		}
	}

	return p, nil
}

// UnmarshalAmino transforms the key/value array to a OptionalData.
func (a *Attributes) UnmarshalAmino(keyValues []KeyValue) error {
	tempMap := make(map[string]string, len(keyValues))
	for _, p := range keyValues {
		tempMap[p.Key] = p.Value
	}

	*a = tempMap

	return nil
}

func (a Attributes) String() string {
	buff := new(bytes.Buffer)
	for key, value := range a {
		fmt.Fprintf(buff, "%s: %s\n", key, value)
	}
	return buff.String()
}
