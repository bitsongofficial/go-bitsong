package types

import "testing"

func TestValidateSymbol(t *testing.T) {
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
			wantErr: false,
			args:    args{denom: "ht"},
		},
		{
			name:    "equal 64 characters in length",
			wantErr: true,
			args:    args{denom: "btc1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct1"},
		},
		{
			name:    "more than 64 characters in length",
			wantErr: true,
			args:    args{denom: "btc1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct12"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSymbol(tt.args.denom); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbol() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
