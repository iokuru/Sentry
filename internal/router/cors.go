package router

import "strings"

func MatchOrigin(allowedOrigin, requestOrigin string) string {
	allowed := strings.ToLower(strings.TrimSpace(allowedOrigin))
	if allowed == "" || allowed == "*" {
		return "*"
	}
	req := strings.ToLower(strings.TrimSpace(requestOrigin))
	if req != "" && req == allowed {
		return allowed
	}
	return ""
}

func RequestOrigin(headers map[string]string) string {
	for k, v := range headers {
		if strings.EqualFold(k, "origin") {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
