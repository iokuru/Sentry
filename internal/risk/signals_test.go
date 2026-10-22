package risk

import (
	"testing"
)

func TestCheckUrgencyWording(t *testing.T) {
	tests := []struct {
		name    string
		input   Input
		wantHit bool
	}{
		{
			name:    "urgent suspension deadline",
			input:   Input{Text: "Your account is flagged! Take immediate action within 24 hours to avoid suspension."},
			wantHit: true,
		},
		{
			name:    "calm informational note",
			input:   Input{Text: "Your weekly grocery flyer is ready to review."},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := checkUrgencyWording(tt.input)
			if got != tt.wantHit {
				t.Fatalf("checkUrgencyWording() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestCheckUntraceablePayment(t *testing.T) {
	tests := []struct {
		name    string
		input   Input
		wantHit bool
	}{
		{
			name:    "apple gift card request",
			input:   Input{Text: "Go to the nearest pharmacy and buy a $500 Apple gift card right away."},
			wantHit: true,
		},
		{
			name:    "crypto payment",
			input:   Input{Text: "Deposit bitcoin at the nearest ATM machine to settle balance."},
			wantHit: true,
		},
		{
			name:    "standard store receipt",
			input:   Input{Text: "Thank you for shopping with us today at the supermarket."},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := checkUntraceablePayment(tt.input)
			if got != tt.wantHit {
				t.Fatalf("checkUntraceablePayment() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestCheckLinkShortenersAndDomains(t *testing.T) {
	tests := []struct {
		name    string
		input   Input
		wantHit bool
	}{
		{
			name:    "bitly link",
			input:   Input{Text: "Click bit.ly/3xXyZ to confirm your parcel."},
			wantHit: true,
		},
		{
			name:    "lookalike domain",
			input:   Input{Text: "Verify details at usps-tracking-update.xyz/auth"},
			wantHit: true,
		},
		{
			name:    "normal text without links",
			input:   Input{Text: "See you on Tuesday for morning tea."},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := checkLinkShortenersAndDomains(tt.input)
			if got != tt.wantHit {
				t.Fatalf("checkLinkShortenersAndDomains() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestCheckOTPAndPIN(t *testing.T) {
	tests := []struct {
		name    string
		input   Input
		wantHit bool
	}{
		{
			name:    "asking for 6 digit verification code",
			input:   Input{Text: "A 6-digit code was sent to your phone. Read the digits back to me."},
			wantHit: true,
		},
		{
			name:    "casual chat",
			input:   Input{Text: "The gate combination is 1234 when you arrive."},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := checkOTPAndPIN(tt.input)
			if got != tt.wantHit {
				t.Fatalf("checkOTPAndPIN() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestCheckAuthorityImpersonation(t *testing.T) {
	tests := []struct {
		name    string
		input   Input
		wantHit bool
	}{
		{
			name:    "irs arrest threat",
			input:   Input{Text: "This is officer Miller from the IRS. An arrest warrant has been issued."},
			wantHit: true,
		},
		{
			name:    "legitimate community newsletter",
			input:   Input{Text: "Local council recycling calendar for the upcoming month."},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := checkAuthorityImpersonation(tt.input)
			if got != tt.wantHit {
				t.Fatalf("checkAuthorityImpersonation() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestCheckNewNumberFamilyClaims(t *testing.T) {
	tests := []struct {
		name    string
		input   Input
		wantHit bool
	}{
		{
			name:    "hi mum broke my phone",
			input:   Input{Text: "Hi mum, broke my phone. This is my new number, text me here."},
			wantHit: true,
		},
		{
			name:    "hi mom only without new number",
			input:   Input{Text: "Hi mom, loved Sunday dinner thank you!"},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := checkNewNumberFamilyClaims(tt.input)
			if got != tt.wantHit {
				t.Fatalf("checkNewNumberFamilyClaims() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}

func TestCheckPIIRequested(t *testing.T) {
	tests := []struct {
		name    string
		input   Input
		wantHit bool
	}{
		{
			name: "entity with ssn",
			input: Input{
				Entities: []Entity{{Type: "SSN", Text: "123-45-6789"}},
			},
			wantHit: true,
		},
		{
			name:    "text with ssn request",
			input:   Input{Text: "Reply with your social security number and cvv to verify."},
			wantHit: true,
		},
		{
			name:    "innocuous text",
			input:   Input{Text: "Meeting at the park tomorrow at ten."},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := checkPIIRequested(tt.input)
			if got != tt.wantHit {
				t.Fatalf("checkPIIRequested() = %v, want %v", got, tt.wantHit)
			}
		})
	}
}
