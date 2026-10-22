package risk

import (
	"regexp"
	"strings"
)

var piiRegex = regexp.MustCompile(`(?i)\b(social\s*security\s*number|ssn|cvv|cvc|card\s*number|expiration\s*date|exp\s*date|driver'?s\s*licen[sc]e|passport\s*number|routing\s*number|mothers?\s*maiden\s*name|date\s*of\s*birth|dob)\b`)

func checkPIIRequested(in Input) (bool, string) {
	for _, entity := range in.Entities {
		switch entity.Type {
		case "SSN", "BANK_ACCOUNT_NUMBER", "CREDIT_DEBIT_NUMBER", "PASSPORT_NUMBER", "DRIVER_ID", "PIN", "PASSWORD":
			return true, "requests sensitive identity, banking, or credit credentials"
		}
	}
	combined := strings.ToLower(in.Text + " " + in.Context)
	if match := piiRegex.FindString(combined); match != "" {
		return true, "requests sensitive identity, banking, or credit credentials"
	}
	return false, ""
}
