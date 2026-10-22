package risk

import (
	"testing"

	"github.com/iokuru/sentry/internal/model"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name      string
		input     Input
		wantScore int
		wantLevel model.RiskLevel
	}{
		{
			name: "high risk multi signal scam",
			input: Input{
				Text: "URGENT: This is officer Davis from the IRS. Your account will be seized within 24 hours. Pay via $500 Apple gift card immediately.",
			},
			wantScore: 70, // urgency(15) + authority(25) + untraceable_payment(30) = 70
			wantLevel: model.RiskHigh,
		},
		{
			name: "medium risk otp request",
			input: Input{
				Text: "Your delivery is waiting. Please share your 6-digit code with the courier.",
			},
			wantScore: 30, // otp_request(30) = 30
			wantLevel: model.RiskMedium,
		},
		{
			name: "low risk benign receipt",
			input: Input{
				Text: "Your order for gardening gloves has been received. Have a pleasant day.",
			},
			wantScore: 0,
			wantLevel: model.RiskLow,
		},
		{
			name: "cap score at 100",
			input: Input{
				Text: "URGENT from IRS police! Hi mum, broke my phone. Send $500 gift card, OTP code to bit.ly/test and your SSN.",
				Entities: []Entity{
					{Type: "SSN", Text: "111-22-3333"},
				},
			},
			wantScore: 100,
			wantLevel: model.RiskHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, level, _ := CalculateScore(tt.input)
			if score != tt.wantScore {
				t.Fatalf("CalculateScore() score = %d, want %d", score, tt.wantScore)
			}
			if level != tt.wantLevel {
				t.Fatalf("CalculateScore() level = %v, want %v", level, tt.wantLevel)
			}
		})
	}
}

func TestCategoryClassification(t *testing.T) {
	tests := []struct {
		text     string
		expected model.Category
	}{
		{"Your Chase bank debit card was charged $450", model.CategoryBank},
		{"USPS package redelivery fee required for shipment", model.CategoryDelivery},
		{"Electric company power disconnection notice for unpaid bill", model.CategoryUtility},
		{"You are our lucky winner of a $50,000 lottery jackpot", model.CategoryPrize},
		{"Windows defender computer virus detected call Microsoft support", model.CategoryTechSupport},
		{"Hi mum, broke my phone this morning and need emergency cash", model.CategoryRelativeInNeed},
		{"My darling sweetheart, please send money for my flight home from the oil rig", model.CategoryRomance},
		{"Notice of unpaid taxes and tax audit from the IRS revenue service", model.CategoryTax},
		{"Meeting for coffee on Friday morning", model.CategoryOther},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			got := ClassifyCategory(Input{Text: tt.text})
			if got != tt.expected {
				t.Fatalf("ClassifyCategory(%q) = %q, want %q", tt.text, got, tt.expected)
			}
		})
	}
}

func TestBuildSeniorOutput(t *testing.T) {
	out := BuildSeniorOutput(model.CategoryBank, model.RiskHigh, "US")
	if out.HoldOff == "" {
		t.Fatal("expected hold off advice, got empty string")
	}
	if out.SafeAlternative == "" {
		t.Fatal("expected safe alternative action, got empty string")
	}
	if out.Hotline.Number == "" || out.Hotline.TelURI == "" {
		t.Fatalf("expected valid hotline, got %+v", out.Hotline)
	}
	if out.Hotline.Country != "US" {
		t.Fatalf("expected country US, got %s", out.Hotline.Country)
	}
}
