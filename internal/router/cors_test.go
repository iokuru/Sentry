package router_test

import (
	"testing"

	"github.com/iokuru/sentry/internal/router"
)

func TestMatchOrigin(t *testing.T) {
	allowed := "https://sentry-security-krisnakk.netlify.app"

	tests := []struct {
		name    string
		allowed string
		request string
		want    string
	}{
		{
			name:    "exact match",
			allowed: allowed,
			request: allowed,
			want:    allowed,
		},
		{
			name:    "uppercase request origin",
			allowed: allowed,
			request: "HTTPS://SENTRY-SECURITY-KRISNAKK.NETLIFY.APP",
			want:    allowed,
		},
		{
			name:    "whitespace trimmed",
			allowed: "  " + allowed + "  ",
			request: " " + allowed + " ",
			want:    allowed,
		},
		{
			name:    "different origin",
			allowed: allowed,
			request: "https://malicious.example.com",
			want:    "",
		},
		{
			name:    "empty request origin",
			allowed: allowed,
			request: "",
			want:    "",
		},
		{
			name:    "wildcard allowed",
			allowed: "*",
			request: "https://any.example.com",
			want:    "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := router.MatchOrigin(tt.allowed, tt.request)
			if got != tt.want {
				t.Errorf("MatchOrigin(%q, %q) = %q; want %q", tt.allowed, tt.request, got, tt.want)
			}
		})
	}
}

func TestRequestOrigin(t *testing.T) {
	headers := map[string]string{
		"Origin": "https://example.com",
	}
	if got := router.RequestOrigin(headers); got != "https://example.com" {
		t.Errorf("RequestOrigin() = %q; want %q", got, "https://example.com")
	}

	headersLower := map[string]string{
		"origin": "https://example.org",
	}
	if got := router.RequestOrigin(headersLower); got != "https://example.org" {
		t.Errorf("RequestOrigin() = %q; want %q", got, "https://example.org")
	}
}
