package comprehend

import (
	"context"
	"testing"

	comprehendsvc "github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
)

type mockClient struct {
	lang string
}

func (m *mockClient) DetectDominantLanguage(ctx context.Context, params *comprehendsvc.DetectDominantLanguageInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectDominantLanguageOutput, error) {
	code := m.lang
	return &comprehendsvc.DetectDominantLanguageOutput{
		Languages: []types.DominantLanguage{
			{LanguageCode: &code},
		},
	}, nil
}

func (m *mockClient) DetectSentiment(ctx context.Context, params *comprehendsvc.DetectSentimentInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectSentimentOutput, error) {
	return &comprehendsvc.DetectSentimentOutput{Sentiment: types.SentimentTypeNeutral}, nil
}

func (m *mockClient) DetectKeyPhrases(ctx context.Context, params *comprehendsvc.DetectKeyPhrasesInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectKeyPhrasesOutput, error) {
	return &comprehendsvc.DetectKeyPhrasesOutput{}, nil
}

func (m *mockClient) DetectPiiEntities(ctx context.Context, params *comprehendsvc.DetectPiiEntitiesInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectPiiEntitiesOutput, error) {
	return &comprehendsvc.DetectPiiEntitiesOutput{}, nil
}

func TestAnalyzeWithFallback(t *testing.T) {
	ctx := context.Background()

	t.Run("nil client returns fallback", func(t *testing.T) {
		res := AnalyzeWithFallback(ctx, nil, "some message")
		if res.Language != "en" {
			t.Fatalf("expected en, got %s", res.Language)
		}
	})

	t.Run("non-english returns rule fallback", func(t *testing.T) {
		client := &mockClient{lang: "es"}
		res := AnalyzeWithFallback(ctx, client, "urgente por favor transferir dinero")
		if res.Language != "es" {
			t.Fatalf("expected es, got %s", res.Language)
		}
		if res.Entities != nil {
			t.Fatalf("expected nil entities on fallback, got %v", res.Entities)
		}
	})
}
