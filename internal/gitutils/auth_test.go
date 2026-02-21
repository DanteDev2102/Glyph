package gitutils

import (
	"testing"
)

func TestValidateNotEmpty(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
	}{
		{"valid", false},
		{"  valid  ", false},
		{"", true},
		{"   ", true},
		{"\t\n", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := validateNotEmpty(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNotEmpty(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
