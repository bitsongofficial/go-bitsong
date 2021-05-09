package types

import "testing"

func TestValidateDenom(t *testing.T) {
	type args struct {
		denom string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "right case",
			wantErr: false,
			args:    args{denom: "btc"},
		},
		{
			name:    "start with a capital letter",
			wantErr: true,
			args:    args{denom: "Btc"},
		},
		{
			name:    "contain a capital letter",
			wantErr: true,
			args:    args{denom: "bTc"},
		},
		{
			name:    "less than 3 characters in length",
			wantErr: true,
			args:    args{denom: "ht"},
		},
		{
			name:    "equal 64 characters in length",
			wantErr: false,
			args:    args{denom: "btc1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct1"},
		},
		{
			name:    "more than 64 characters in length",
			wantErr: true,
			args:    args{denom: "btc1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct12"},
		},
		{
			name:    "contain peg",
			wantErr: true,
			args:    args{denom: "pegeth"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDenom(tt.args.denom); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDenom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateKeywords(t *testing.T) {
	type args struct {
		denom string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "right case", args: args{denom: "stake"}, wantErr: false},
		{name: "denom is peg", args: args{denom: "peg"}, wantErr: true},
		{name: "denom is Peg", args: args{denom: "Peg"}, wantErr: false},
		{name: "denom begin with peg", args: args{denom: "pegtoken"}, wantErr: true},
		{name: "denom is ibc", args: args{denom: "ibc"}, wantErr: true},
		{name: "denom is IBC", args: args{denom: "Peg"}, wantErr: false},
		{name: "denom begin with ibc", args: args{denom: "ibctoken"}, wantErr: true},
		{name: "denom is swap", args: args{denom: "swap"}, wantErr: true},
		{name: "denom is SWAP", args: args{denom: "SWAP"}, wantErr: false},
		{name: "denom begin with swap", args: args{denom: "swaptoken"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateKeywords(tt.args.denom); (err != nil) != tt.wantErr {
				t.Errorf("CheckKeywords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
