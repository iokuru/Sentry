package hotline

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/iokuru/sentry/internal/model"
)

var (
	registryMu sync.RWMutex
	registry   = map[string]model.HotlineInfo{
		"US": {
			Country: "US",
			Name:    "US FTC Fraud Hotline",
			Number:  "1-877-382-4357",
			TelURI:  "tel:+18773824357",
			URL:     "https://reportfraud.ftc.gov",
		},
		"GB": {
			Country: "GB",
			Name:    "UK Action Fraud",
			Number:  "0300 123 2040",
			TelURI:  "tel:+443001232040",
			URL:     "https://www.actionfraud.police.uk",
		},
		"IN": {
			Country: "IN",
			Name:    "India National Cyber Crime Helpline",
			Number:  "1930",
			TelURI:  "tel:1930",
			URL:     "https://cybercrime.gov.in",
		},
		"AU": {
			Country: "AU",
			Name:    "Australia Scamwatch",
			Number:  "1300 795 995",
			TelURI:  "tel:+611300795995",
			URL:     "https://www.scamwatch.gov.au",
		},
	}
)

func Lookup(country string) model.HotlineInfo {
	code := strings.ToUpper(strings.TrimSpace(country))
	registryMu.RLock()
	defer registryMu.RUnlock()

	if h, ok := registry[code]; ok {
		return h
	}
	// default to us ftc when country is unrecognised
	return registry["US"]
}

func ApplyOverrides(rawJSON string) error {
	if strings.TrimSpace(rawJSON) == "" || rawJSON == "{}" {
		return nil
	}
	var overrides map[string]model.HotlineInfo
	if err := json.Unmarshal([]byte(rawJSON), &overrides); err != nil {
		return err
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	for k, v := range overrides {
		registry[strings.ToUpper(k)] = v
	}
	return nil
}
