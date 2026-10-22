package comprehend

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/iokuru/sentry/internal/risk"
)

type Client interface {
	DetectDominantLanguage(ctx context.Context, params *comprehend.DetectDominantLanguageInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectDominantLanguageOutput, error)
	DetectSentiment(ctx context.Context, params *comprehend.DetectSentimentInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectSentimentOutput, error)
	DetectKeyPhrases(ctx context.Context, params *comprehend.DetectKeyPhrasesInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectKeyPhrasesOutput, error)
	DetectPiiEntities(ctx context.Context, params *comprehend.DetectPiiEntitiesInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectPiiEntitiesOutput, error)
}

type AnalysisResult struct {
	Language   string
	Sentiment  string
	KeyPhrases []string
	Entities   []risk.Entity
}

func Analyze(ctx context.Context, client Client, text string) (*AnalysisResult, error) {
	if text == "" {
		return &AnalysisResult{Language: "en"}, nil
	}

	// cap text at 5000 bytes to stay within comprehend limits
	sample := text
	if len(sample) > 4900 {
		sample = sample[:4900]
	}

	res := &AnalysisResult{Language: "en"}

	langOut, err := client.DetectDominantLanguage(ctx, &comprehend.DetectDominantLanguageInput{
		Text: &sample,
	})
	if err == nil && len(langOut.Languages) > 0 && langOut.Languages[0].LanguageCode != nil {
		res.Language = *langOut.Languages[0].LanguageCode
	}

	// comprehend pii and sentiment in free tier only operate on english
	if res.Language != "en" {
		return res, nil
	}

	langCode := types.LanguageCodeEn

	sentOut, err := client.DetectSentiment(ctx, &comprehend.DetectSentimentInput{
		LanguageCode: langCode,
		Text:         &sample,
	})
	if err == nil {
		res.Sentiment = string(sentOut.Sentiment)
	}

	keyOut, err := client.DetectKeyPhrases(ctx, &comprehend.DetectKeyPhrasesInput{
		LanguageCode: langCode,
		Text:         &sample,
	})
	if err == nil {
		for _, phrase := range keyOut.KeyPhrases {
			if phrase.Text != nil {
				res.KeyPhrases = append(res.KeyPhrases, *phrase.Text)
			}
		}
	}

	piiOut, err := client.DetectPiiEntities(ctx, &comprehend.DetectPiiEntitiesInput{
		LanguageCode: langCode,
		Text:         &sample,
	})
	if err != nil {
		return nil, fmt.Errorf("detect pii entities: %w", err)
	}

	for _, pii := range piiOut.Entities {
		score := 0.0
		if pii.Score != nil {
			score = float64(*pii.Score)
		}
		res.Entities = append(res.Entities, risk.Entity{
			Type:  string(pii.Type),
			Score: score,
		})
	}

	return res, nil
}
