package hotline

import (
	"strings"
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
		{"en-CA", "CA"},
		{"en-NZ", "NZ"},
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

func TestHotlineRegistryEntries(t *testing.T) {
	countries := []string{"US", "GB", "IN", "AU", "CA", "NZ"}

	for _, c := range countries {
		t.Run(c, func(t *testing.T) {
			h := Lookup(c)
			if h.Country != c {
				t.Errorf("expected country %q, got %q", c, h.Country)
			}
			if h.Name == "" {
				t.Errorf("expected non-empty name for %q", c)
			}
			if h.TelURI != "" && !strings.HasPrefix(h.TelURI, "tel:") {
				t.Errorf("expected TelURI to start with 'tel:', got %q", h.TelURI)
			}
			if !strings.HasPrefix(h.URL, "https://") {
				t.Errorf("expected URL to start with 'https://', got %q", h.URL)
			}
		})
	}
}
