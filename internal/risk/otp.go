package risk

import (
	"regexp"
	"strings"
)

var otpRegex = regexp.MustCompile(`(?i)\b(otp|one-time\s*passcode|verification\s*code|6-digit\s*code|security\s*code|secret\s*pin|pin|passcode|password|passwords|share\s*the\s*code|send\s*the\s*digits|read\s*the\s*code|two-factor\s*code|2fa\s*code|password\s*reset\s*code)\b`)

func checkOTPAndPIN(in Input) (bool, string) {
	combined := strings.ToLower(in.Text + " " + in.Context)
	if match := otpRegex.FindString(combined); match != "" {
		return true, "requests a private one-time verification code, pin, or password"
	}
	return false, ""
}
