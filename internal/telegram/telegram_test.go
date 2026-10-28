package telegram

import (
	"strings"
	"testing"
)

func TestFormatAlertMessage(t *testing.T) {
	msg, buttons := FormatAlertMessage(
		"bank",
		"high",
		"impersonated bank asking for immediate payment",
		"US FTC",
		"1-877-382-4357",
		"inc-12345",
		"2026-10-29T03:00:00Z",
	)

	if !strings.Contains(msg, "*SENTRY ALERT*") {
		t.Errorf("expected alert title in message, got %s", msg)
	}
	if !strings.Contains(msg, "HIGH") {
		t.Errorf("expected HIGH level in message, got %s", msg)
	}
	if !strings.Contains(msg, "US FTC at 1-877-382-4357") {
		t.Errorf("expected plain text hotline in message, got %s", msg)
	}
	if !strings.Contains(msg, "inc-12345") {
		t.Errorf("expected incident id in message, got %s", msg)
	}

	if len(buttons) != 1 || len(buttons[0]) != 2 {
		t.Fatalf("expected 1 row of 2 buttons, got %v", buttons)
	}

	ackBtn := buttons[0][0]
	if ackBtn.Text != "Acknowledge" || ackBtn.CallbackData != "ack:inc-12345" {
		t.Errorf("unexpected acknowledge button: %+v", ackBtn)
	}

	resBtn := buttons[0][1]
	if resBtn.Text != "Resolve" || resBtn.CallbackData != "res:inc-12345" {
		t.Errorf("unexpected resolve button: %+v", resBtn)
	}
}

func TestParseCallbackData(t *testing.T) {
	tests := []struct {
		input      string
		wantAction string
		wantID     string
		wantOK     bool
	}{
		{"ack:inc-99", "ack", "inc-99", true},
		{"res:inc-42", "res", "inc-42", true},
		{"invalid", "", "", false},
		{"other:inc-1", "", "", false},
		{"", "", "", false},
	}

	for _, tt := range tests {
		action, id, ok := ParseCallbackData(tt.input)
		if ok != tt.wantOK || action != tt.wantAction || id != tt.wantID {
			t.Errorf("ParseCallbackData(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.input, action, id, ok, tt.wantAction, tt.wantID, tt.wantOK)
		}
	}
}

func TestVerifyWebhookSecret(t *testing.T) {
	tests := []struct {
		header   string
		expected string
		want     bool
	}{
		{"secret-token-123", "secret-token-123", true},
		{"secret-token-123", "wrong-secret", false},
		{"", "secret-token-123", false},
		{"secret-token-123", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got := VerifyWebhookSecret(tt.header, tt.expected)
		if got != tt.want {
			t.Errorf("VerifyWebhookSecret(%q, %q) = %v, want %v", tt.header, tt.expected, got, tt.want)
		}
	}
}

func TestIsValidValue(t *testing.T) {
	tests := []struct {
		val  string
		want bool
	}{
		{"PLACEHOLDER", false},
		{"placeholder", false},
		{"", false},
		{"  ", false},
		{"123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11", true},
		{"987654321", true},
	}

	for _, tt := range tests {
		got := isValidValue(tt.val)
		if got != tt.want {
			t.Errorf("isValidValue(%q) = %v, want %v", tt.val, got, tt.want)
		}
	}
}
