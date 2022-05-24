package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"os"
	"sort"
	"strconv"
	"time"
)

func parseTime(timeStr string) (time.Time, error) {
	var startTime time.Time
	if timeStr == "" { // empty start time
		startTime = time.Unix(0, 0)
	} else if timeUnix, err := strconv.ParseInt(timeStr, 10, 64); err == nil { // unix time
		startTime = time.Unix(timeUnix, 0)
	} else if timeRFC, err := time.Parse(time.RFC3339, timeStr); err == nil { // RFC time
		startTime = timeRFC
	} else { // invalid input
		return startTime, errors.New("invalid start time format")
	}

	return startTime, nil
}

type Account struct {
	address sdk.AccAddress
	amount  sdk.Int
}

type ClaimInfo struct {
	Index  uint64   `json:"index"`
	Amount string   `json:"amount"`
	Proof  []string `json:"proof"`
}

func AccountsFromMap(accMap map[string]string) ([]*Account, error) {
	i := 0

	accsMap := make([]*Account, len(accMap))

	for strAddr, strAmt := range accMap {
		amt, ok := sdk.NewIntFromString(strAmt)
		if !ok {
			return nil, fmt.Errorf("could not cast %s to sdk.Int", strAddr)
		}

		addr, err := sdk.AccAddressFromBech32(strAddr)
		if err != nil {
			return nil, fmt.Errorf("could not cast %s to sdk.AccAddress", strAddr)
		}

		accsMap[i] = &Account{
			address: addr,
			amount:  amt,
		}
		i++
	}

	return accsMap, nil
}

func CreateDistributionList(accounts []*Account) (Tree, map[string]ClaimInfo, sdk.Int, error) {
	// sort lists by coin amount
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].amount.LT(accounts[j].amount)
	})

	totalAmt := sdk.ZeroInt()

	nodes := make([][]byte, len(accounts))
	for i, acc := range accounts {
		indexStr := strconv.FormatUint(uint64(i), 10)
		nodes[i] = []byte(fmt.Sprintf("%s%s%s", indexStr, acc.address.String(), acc.amount.String()))
		totalAmt = totalAmt.Add(acc.amount)
	}

	tree := NewTree(nodes...)

	addrToProof := make(map[string]ClaimInfo, len(accounts))

	for i, acc := range accounts {
		proof := ProofBytesToString(tree.Proof(crypto.Sha256(nodes[i])))

		addrToProof[acc.address.String()] = ClaimInfo{
			Index:  uint64(i),
			Amount: acc.amount.String(),
			Proof:  proof,
		}
	}

	return tree, addrToProof, totalAmt, nil
}

func ProofBytesToString(proof [][]byte) []string {
	str := make([]string, len(proof)-1)
	for i, p := range proof {
		if i == len(proof)-1 {
			continue
		}
		str[i] = fmt.Sprintf("%x", p)
	}
	return str
}

func createFile(filename string, contents interface{}) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("could not create file: %v", err)
	}
	totalBytes, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not marshal data: %v", err)
	}
	if _, err := file.Write(totalBytes); err != nil {
		return nil, fmt.Errorf("could not write data: %v", err)
	}
	return file, nil
}
