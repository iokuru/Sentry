package risk

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/iokuru/sentry/internal/model"
)

var (
	bannedWords = []string{
		"ai",
		"algorithm",
		"model",
		"confidence score",
		"false positive",
		"probability",
	}

	falseGuarantees = []string{
		"100% safe",
		"definitely a scam",
		"guaranteed",
	}

	victimBlaming = []string{
		"you shouldn't have",
		"you should not have",
		"be more careful",
		"your fault",
	}

	hedgingContradictions = []string{
		"this is safe but could be dangerous",
		"safe but dangerous",
	}

	allowedAcronyms = map[string]bool{
		"OTP":  true,
		"PIN":  true,
		"SMS":  true,
		"UPI":  true,
		"KYC":  true,
		"PAN":  true,
		"URL":  true,
		"UK":   true,
		"US":   true,
		"IRS":  true,
		"HMRC": true,
		"CBI":  true,
		"ATO":  true,
		"FTC":  true,
		"IC3":  true,
		"DHL":  true,
		"RX":   true,
		"CVS":  true,
		"SBI":  true,
		"HDFC": true,
		"ATM":  true,
		"BTC":  true,
		"ETH":  true,
		"USDT": true,
	}

	wordRegex = regexp.MustCompile(`\b[A-Za-z]+\b`)
)

func ValidateResultWording(out model.SeniorOutput) []string {
	var errors []string

	checkField := func(fieldName, text string, maxLen int) {
		if len(text) > maxLen {
			errors = append(errors, fmt.Sprintf("%s exceeds max length of %d (actual: %d)", fieldName, maxLen, len(text)))
		}

		lower := strings.ToLower(text)

		for _, banned := range bannedWords {
			pattern := fmt.Sprintf(`(?i)\b%s\b`, regexp.QuoteMeta(banned))
			if regexp.MustCompile(pattern).MatchString(text) {
				errors = append(errors, fmt.Sprintf("%s contains banned word: %q", fieldName, banned))
			}
		}

		for _, fg := range falseGuarantees {
			if strings.Contains(lower, fg) {
				errors = append(errors, fmt.Sprintf("%s contains false guarantee: %q", fieldName, fg))
			}
		}

		for _, vb := range victimBlaming {
			if strings.Contains(lower, vb) {
				errors = append(errors, fmt.Sprintf("%s contains victim-blaming phrase: %q", fieldName, vb))
			}
		}

		for _, hc := range hedgingContradictions {
			if strings.Contains(lower, hc) {
				errors = append(errors, fmt.Sprintf("%s contains contradictory hedging: %q", fieldName, hc))
			}
		}

		if strings.Contains(text, "!") {
			errors = append(errors, fmt.Sprintf("%s contains exclamation mark", fieldName))
		}

		if strings.Contains(text, "—") {
			errors = append(errors, fmt.Sprintf("%s contains em dash", fieldName))
		}

		words := wordRegex.FindAllString(text, -1)
		for _, w := range words {
			if len(w) >= 2 && isAllUpper(w) {
				if !allowedAcronyms[w] {
					errors = append(errors, fmt.Sprintf("%s contains unauthorized all-caps word: %q", fieldName, w))
				}
			}
		}
	}

	checkField("headline", out.Headline, 80)
	checkField("what_this_is", out.WhatThisIs, 220)
	checkField("what_to_do_next", out.WhatToDoNext, 220)
	checkField("hold_off", out.HoldOff, 220)
	checkField("safe_alternative", out.SafeAlternative, 220)

	return errors
}

func isAllUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
