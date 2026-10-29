package comprehend_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	comprehendsvc "github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/iokuru/sentry/internal/comprehend"
)

type fullMockClient struct {
	langOutput      *comprehendsvc.DetectDominantLanguageOutput
	langErr         error
	sentimentOutput *comprehendsvc.DetectSentimentOutput
	sentimentErr    error
	keyOutput       *comprehendsvc.DetectKeyPhrasesOutput
	keyErr          error
	piiOutput       *comprehendsvc.DetectPiiEntitiesOutput
	piiErr          error
}

func (m *fullMockClient) DetectDominantLanguage(ctx context.Context, params *comprehendsvc.DetectDominantLanguageInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectDominantLanguageOutput, error) {
	if m.langErr != nil {
		return nil, m.langErr
	}
	return m.langOutput, nil
}

func (m *fullMockClient) DetectSentiment(ctx context.Context, params *comprehendsvc.DetectSentimentInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectSentimentOutput, error) {
	if m.sentimentErr != nil {
		return nil, m.sentimentErr
	}
	return m.sentimentOutput, nil
}

func (m *fullMockClient) DetectKeyPhrases(ctx context.Context, params *comprehendsvc.DetectKeyPhrasesInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectKeyPhrasesOutput, error) {
	if m.keyErr != nil {
		return nil, m.keyErr
	}
	return m.keyOutput, nil
}

func (m *fullMockClient) DetectPiiEntities(ctx context.Context, params *comprehendsvc.DetectPiiEntitiesInput, optFns ...func(*comprehendsvc.Options)) (*comprehendsvc.DetectPiiEntitiesOutput, error) {
	if m.piiErr != nil {
		return nil, m.piiErr
	}
	return m.piiOutput, nil
}

func TestAnalyze_EmptyText(t *testing.T) {
	client := &fullMockClient{}
	res, err := comprehend.Analyze(context.Background(), client, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Language != "en" {
		t.Errorf("expected language en, got %s", res.Language)
	}
}

func TestAnalyze_EnglishScamFixtures(t *testing.T) {
	client := &fullMockClient{
		langOutput: &comprehendsvc.DetectDominantLanguageOutput{
			Languages: []types.DominantLanguage{
				{LanguageCode: aws.String("en")},
			},
		},
		sentimentOutput: &comprehendsvc.DetectSentimentOutput{
			Sentiment: types.SentimentTypeNegative,
		},
		keyOutput: &comprehendsvc.DetectKeyPhrasesOutput{
			KeyPhrases: []types.KeyPhrase{
				{Text: aws.String("fraud division")},
				{Text: aws.String("immediate wire transfer")},
			},
		},
		piiOutput: &comprehendsvc.DetectPiiEntitiesOutput{
			Entities: []types.PiiEntity{
				{
					Type:  types.PiiEntityTypeBankAccountNumber,
					Score: aws.Float32(0.99),
				},
				{
					Type:  types.PiiEntityTypeSsn,
					Score: aws.Float32(0.95),
				},
			},
		},
	}

	res, err := comprehend.Analyze(context.Background(), client, "Urgent notice from bank fraud division")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Language != "en" {
		t.Errorf("expected language en, got %s", res.Language)
	}
	if res.Sentiment != "NEGATIVE" {
		t.Errorf("expected NEGATIVE sentiment, got %s", res.Sentiment)
	}
	if len(res.KeyPhrases) != 2 {
		t.Errorf("expected 2 key phrases, got %d", len(res.KeyPhrases))
	}
	if len(res.Entities) != 2 {
		t.Fatalf("expected 2 pii entities, got %d", len(res.Entities))
	}
	if res.Entities[0].Type != "BANK_ACCOUNT_NUMBER" || res.Entities[0].Score < 0.9 {
		t.Errorf("unexpected first entity: %+v", res.Entities[0])
	}
}

func TestAnalyze_NonEnglishLanguage(t *testing.T) {
	client := &fullMockClient{
		langOutput: &comprehendsvc.DetectDominantLanguageOutput{
			Languages: []types.DominantLanguage{
				{LanguageCode: aws.String("fr")},
			},
		},
	}

	res, err := comprehend.Analyze(context.Background(), client, "Bonjour, votre compte est bloque")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Language != "fr" {
		t.Errorf("expected fr, got %s", res.Language)
	}
	if len(res.Entities) != 0 {
		t.Errorf("expected zero entities for non-english, got %d", len(res.Entities))
	}
}

func TestAnalyze_PiiErrorPropagation(t *testing.T) {
	client := &fullMockClient{
		langOutput: &comprehendsvc.DetectDominantLanguageOutput{
			Languages: []types.DominantLanguage{
				{LanguageCode: aws.String("en")},
			},
		},
		sentimentOutput: &comprehendsvc.DetectSentimentOutput{
			Sentiment: types.SentimentTypeNeutral,
		},
		keyOutput: &comprehendsvc.DetectKeyPhrasesOutput{},
		piiErr:    errors.New("throttling exception"),
	}

	_, err := comprehend.Analyze(context.Background(), client, "Test message")
	if err == nil {
		t.Fatalf("expected error from pii failure, got nil")
	}
}
