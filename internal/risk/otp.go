package risk

import (
	"regexp"
)

var (
	advisoryRegex   = regexp.MustCompile(`(?i)\b(do\s*not|don'?t|never|no\s*one\s*will\s*ask\s*you\s*to)\s*(share|send|tell|give|provide|disclose)\s*(this|the|your|ur)?\s*(6[- ]digit\s*)?(code|otp(\s*pin)?|pin|passcode|digits)\b`)
	otpRequestRegex = regexp.MustCompile(`(?i)\b((share|send|tell|give|provide|read)\s*(me|us|back)?\s*(the|your|ur|this)?\s*(6[- ]digit\s*)?(code|otp(\s*pin)?|pin|passcode|digits)|enter\s*(your\s*)?(secret\s*)?(upi\s*)?pin\s*to\s*(receive|credit|claim|cancel|reverse)|code\s*(you\s*)?just\s*received\s*to\s*(cancel|reverse|verify|stop)|security\s*passcode\s*immediately\s*to\s*(reverse|cancel)|provide\s*to\s*representative\s*to\s*verify)\b`)
)

func checkOTPAndPIN(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	sanitized := advisoryRegex.ReplaceAllString(combined, "")
	if otpRequestRegex.MatchString(sanitized) {
		return true, "requests a private one-time verification code, pin, or password"
	}
	return false, ""
}
