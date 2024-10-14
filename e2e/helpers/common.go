package helpers

import "testing"

func debugOutput(t *testing.T, stdout string) {
	if true {
		t.Log(stdout)
	}
}
