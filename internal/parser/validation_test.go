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
		{"a", false},
		{"", true},                                     // Empty name
		{strings.Repeat("a", 51), true},               // Too long
		{"invalid name", true},                         // Space
		{"invalid@name", true},                         // Special char
		{"config", true},                               // Reserved word
		{"help", true},                                 // Reserved word
		{"init", true},                                 // Reserved word
		{"create", true},                               // Reserved word
		{"rm", true},                                   // Reserved word
		{"set", true},                                  // Reserved word
		{"list", true},                                 // Reserved word
		{"glyph", true},                                // Reserved word
		{"completion", true},                           // Reserved word
		{"valid-Name-123_456", false},                  // Mixed
	}

	for _, tt := range tests {
		err := ValidateTemplateName(tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateTemplateName(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
