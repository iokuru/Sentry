package risk

import (
	"regexp"
	"strings"
)

var untraceablePaymentRegex = regexp.MustCompile(`(?i)\b(gift\s*cards?|apple\s*gift\s*card|itunes\s*card|google\s*play\s*card|steam\s*card|target\s*gift\s*card|vanilla\s*visa|bitcoin|crypto|cryptocurrency|btc|usdt|eth|bitcoin\s*atm|western\s*union|moneygram|wire(\s*transfer)?|cash\s*in\s*mail|courier\s*pickup|prepaid\s*card|green\s*dot)\b`)

func checkUntraceablePayment(in Input) (bool, string) {
	combined := strings.ToLower(in.Text + " " + in.Context)
	if match := untraceablePaymentRegex.FindString(combined); match != "" {
		return true, "asks for untraceable payment via gift card, cryptocurrency, or wire transfer"
	}
	return false, ""
}
