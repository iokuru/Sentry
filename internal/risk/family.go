package risk

import (
	"regexp"
	"strings"
)

var familyGreetingRegex = regexp.MustCompile(`(?i)\b(hi\s*mum|hi\s*mom|hey\s*mum|hey\s*mom|hello\s*mum|hello\s*mom|hey\s*dad|hi\s*dad|grandma|grandpa|grandchild|daughter|son)\b`)
var newNumberRegex = regexp.MustCompile(`(?i)\b(new\s*number|broke\s*my\s*phone|dropped\s*my\s*phone|lost\s*my\s*phone|temporary\s*number|save\s*this\s*number|text\s*me\s*on\s*this|new\s*phone|can't\s*call\s*right\s*now|don'?t\s*call)\b`)

func checkNewNumberFamilyClaims(in Input) (bool, string) {
	combined := strings.ToLower(in.Text + " " + in.Context)
	if familyGreetingRegex.MatchString(combined) && newNumberRegex.MatchString(combined) {
		return true, "claims to be a family member texting from an unfamiliar or replacement number"
	}
	return false, ""
}
