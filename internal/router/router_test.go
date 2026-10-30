package router

import "testing"

func TestRewriteURI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/", "/index.html"},
		{"/caregiver", "/caregiver.html"},
		{"/caregiver/", "/caregiver/index.html"},
		{"/config.json", "/config.json"},
		{"/_next/static/chunks/app.js", "/_next/static/chunks/app.js"},
		{"/tesseract/worker.min.js", "/tesseract/worker.min.js"},
	}

	for _, tc := range tests {
		actual := RewriteURI(tc.input)
		if actual != tc.expected {
			t.Errorf("RewriteURI(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}
