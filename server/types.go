package server

type ValidatorInfo struct {
	Val       string `json:"val"`
	NumDels   int    `json:"num_dels"`
	NumTokens int    `json:"num_tokens"`
	Jailed    bool   `json:"jailed"`
}
