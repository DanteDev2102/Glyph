package parser

import (
	"strings"
	"testing"
)

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid-name", false},
		{"valid_name_123", false},
		{"A", false},
		{"", true},
		{strings.Repeat("a", 51), true},
		{"invalid name", true},
		{"invalid@name", true},
		{"invalid.name", true},
		{"config", true},
		{"CONFIG", true},
		{"help", true},
		{"init", true},
		{"create", true},
		{"rm", true},
		{"set", true},
		{"list", true},
		{"glyph", true},
		{"completion", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateName(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateName(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
