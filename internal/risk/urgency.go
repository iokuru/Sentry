package risk

import (
	"regexp"
	"strings"
)

var urgencyRegex = regexp.MustCompile(`(?i)\b(urgent|urgently|immediate|immediately|right now|within 24 hours|within 1 hour|act now|hurry|final notice|last warning|expires today|deadline|action required|suspended today|account suspension|immediate action)\b`)

func checkUrgencyWording(in Input) (bool, string) {
	combined := strings.ToLower(in.Text + " " + in.Context)
	if match := urgencyRegex.FindString(combined); match != "" {
		return true, "demands immediate action or threatens quick consequences"
	}
	return false, ""
}
