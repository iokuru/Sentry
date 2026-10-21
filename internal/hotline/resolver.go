package hotline

import (
	"strings"
)

func ResolveCountry(locale string) string {
	raw := strings.TrimSpace(locale)
	if raw == "" {
		return "US"
	}

	raw = strings.ReplaceAll(raw, "_", "-")
	parts := strings.Split(raw, "-")

	var candidate string
	if len(parts) >= 2 {
		candidate = strings.ToUpper(parts[1])
	} else {
		candidate = strings.ToUpper(parts[0])
	}

	switch candidate {
	case "UK":
		return "GB"
	case "US", "USA":
		return "US"
	case "GB", "GBR":
		return "GB"
	case "IN", "IND":
		return "IN"
	case "AU", "AUS":
		return "AU"
	}

	// check primary language tag if region was omitted
	primary := strings.ToLower(parts[0])
	switch primary {
	case "hi", "mr", "ta", "te", "bn", "gu":
		return "IN"
	}

	return "US"
}
