package risk

import (
	"regexp"
	"strings"
)

var shortenerRegex = regexp.MustCompile(`(?i)\b(bit\.ly|tinyurl\.com|t\.co|goo\.gl|is\.gd|cutt\.ly|rb\.gy|ow\.ly|buff\.ly|rebrand\.ly)/[a-zA-Z0-9_\-]+`)
var lookalikeRegex = regexp.MustCompile(`(?i)\b[a-zA-Z0-9\-]+-(bank|paypal|usps|auspost|ups|fedex|amazon|apple|netflix|chase|wellsfargo|verify|security|login|update)\.(com|net|org|xyz|top|link|live|info)\b`)
var suspiciousTldRegex = regexp.MustCompile(`(?i)https?://[a-zA-Z0-9.\-]+\.(xyz|top|tk|ml|ga|cf|click|loan|work|buzz|club|online)/`)

func checkLinkShortenersAndDomains(in Input) (bool, string) {
	combined := strings.ToLower(in.Text + " " + in.Context)
	if shortenerRegex.MatchString(combined) {
		return true, "uses a hidden shortened link to obscure destination"
	}
	if lookalikeRegex.MatchString(combined) || suspiciousTldRegex.MatchString(combined) {
		return true, "uses an imitation domain or suspicious website extension"
	}
	return false, ""
}
