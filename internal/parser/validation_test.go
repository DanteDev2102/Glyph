package parser

import "testing"

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid-name", false},
		{"ValidName123", false},
		{"valid_name", false},
		{"", true},
		{"this-name-is-exactly-fifty-characters-long-123456", false},
		{"this-name-is-too-long-because-it-has-more-than-fifty-chars", true},
		{"invalid name", true},
		{"invalid@name", true},
		{"config", true},
		{"HELP", true},
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
				t.Errorf("ValidateTemplateName(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
