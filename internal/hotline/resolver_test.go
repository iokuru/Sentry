package hotline

import (
	"testing"
)

func TestResolveCountry(t *testing.T) {
	tests := []struct {
		locale   string
		expected string
	}{
		{"en-US", "US"},
		{"en_US", "US"},
		{"en-GB", "GB"},
		{"en-AU", "AU"},
		{"en-IN", "IN"},
		{"hi-IN", "IN"},
		{"hi", "IN"},
		{"UK", "GB"},
		{"AU", "AU"},
		{"", "US"},
		{"unknown-region", "US"},
		{"fr-FR", "US"},
	}

	for _, tt := range tests {
		t.Run(tt.locale, func(t *testing.T) {
			got := ResolveCountry(tt.locale)
			if got != tt.expected {
				t.Fatalf("ResolveCountry(%q) = %q, want %q", tt.locale, got, tt.expected)
			}
		})
	}
}
