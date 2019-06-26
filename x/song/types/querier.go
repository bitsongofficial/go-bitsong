package types

import "strings"

// Query Result Payload for a names query
type QueryResTitles []string

// implement fmt.Stringer
func (n QueryResTitles) String() string {
	return strings.Join(n[:], "\n")
}