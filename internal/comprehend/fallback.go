package comprehend

import (
	"context"
	"strings"
)

func FallbackAnalysis(text string, lang string) *AnalysisResult {
	// rule fallback when comprehend language is unsupported or api call fails
	return &AnalysisResult{
		Language:   lang,
		Sentiment:  "NEUTRAL",
		KeyPhrases: extractSimplePhrases(text),
		Entities:   nil,
	}
}

func AnalyzeWithFallback(ctx context.Context, client Client, text string) *AnalysisResult {
	if client == nil {
		return FallbackAnalysis(text, "en")
	}

	res, err := Analyze(ctx, client, text)
	if err != nil {
		// fall back to pure rule scoring on error
		return FallbackAnalysis(text, "en")
	}

	if res.Language != "en" {
		return FallbackAnalysis(text, res.Language)
	}

	return res
}

func extractSimplePhrases(text string) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var phrases []string
	for i := 0; i < len(words); i += 3 {
		end := i + 3
		if end > len(words) {
			end = len(words)
		}
		phrases = append(phrases, strings.Join(words[i:end], " "))
	}
	return phrases
}
