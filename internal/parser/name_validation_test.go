package parser

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"validname", false},
		{"valid-name", false},
		{"valid_name", false},
		{"valid123", false},
		{"config", true},
		{"my template", true},
		{"name!", true},
		{"", true},
		{"../name", true},
		{"/", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateName(tt.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
