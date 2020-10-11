package types

type URLs map[string]string
type EIDs map[string]string

type ID string

func (id ID) String() string {
	return string(id)
}
