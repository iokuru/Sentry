package risk

import (
	"regexp"
	"strings"
)

var authorityRegex = regexp.MustCompile(`(?i)\b(internal\s*revenue\s*service|irs|federal\s*trade\s*commission|ftc|social\s*security\s*administration|ssa|hmrc|australian\s*taxation\s*office|ato|federal\s*police|sheriff|police\s*department|fbi|customs\s*and\s*border|microsoft|apple\s*support|windows\s*defender|geek\s*squad|arrest\s*warrant|criminal\s*charges?|threat\s*of\s*arrest|court\s*summons|asset\s*freeze|account\s*closure|jail\s*time|legal\s*action\s*taken)\b`)

func checkAuthorityImpersonation(in Input) (bool, string) {
	combined := strings.ToLower(in.Text + " " + in.Context)
	if match := authorityRegex.FindString(combined); match != "" {
		return true, "claims official government or legal authority with threats of arrest or account closure"
	}
	return false, ""
}
