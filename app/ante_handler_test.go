package app

import (
	"github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"testing"
)

func TestNewAnteHandler(t *testing.T) {
	type args struct {
		options HandlerOptions
	}
	tests := []struct {
		name    string
		args    args
		want    types.AnteHandler
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAnteHandler(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnteHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnteHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
