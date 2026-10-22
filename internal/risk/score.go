package risk

import (
	"github.com/iokuru/sentry/internal/model"
)

type Rule struct {
	ID      string
	Name    string
	Weight  int
	Checker func(in Input) (bool, string)
}

var Rules = []Rule{
	{
		ID:      "untraceable_payment",
		Name:    "untraceable payment request",
		Weight:  30,
		Checker: checkUntraceablePayment,
	},
	{
		ID:      "otp_request",
		Name:    "one-time passcode or pin request",
		Weight:  30,
		Checker: checkOTPAndPIN,
	},
	{
		ID:      "authority_impersonation",
		Name:    "authority impersonation and threats",
		Weight:  25,
		Checker: checkAuthorityImpersonation,
	},
	{
		ID:      "family_impersonation",
		Name:    "new number family claim",
		Weight:  25,
		Checker: checkNewNumberFamilyClaims,
	},
	{
		ID:      "suspicious_link",
		Name:    "shortened or lookalike link",
		Weight:  20,
		Checker: checkLinkShortenersAndDomains,
	},
	{
		ID:      "pii_harvesting",
		Name:    "sensitive personal identity requested",
		Weight:  20,
		Checker: checkPIIRequested,
	},
	{
		ID:      "urgency_wording",
		Name:    "urgency and time pressure",
		Weight:  15,
		Checker: checkUrgencyWording,
	},
}

func CalculateScore(in Input) (int, model.RiskLevel, []model.SignalMatch) {
	total := 0
	var matches []model.SignalMatch

	for _, rule := range Rules {
		matched, desc := rule.Checker(in)
		if matched {
			total += rule.Weight
			matches = append(matches, model.SignalMatch{
				ID:          rule.ID,
				Name:        rule.Name,
				Weight:      rule.Weight,
				Description: desc,
			})
		}
	}

	score := total
	if score > 100 {
		score = 100
	}

	var level model.RiskLevel
	switch {
	case score >= 60:
		level = model.RiskHigh
	case score >= 30:
		level = model.RiskMedium
	default:
		level = model.RiskLow
	}

	return score, level, matches
}
