package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	comprehendsvc "github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	textractsvc "github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/iokuru/sentry/internal/comprehend"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/quota"
	"github.com/iokuru/sentry/internal/risk"
	"github.com/iokuru/sentry/internal/store"
	"github.com/iokuru/sentry/internal/textract"
)

type AnalysisQueueMessage struct {
	IncidentID  string `json:"incident_id"`
	ImageKey    string `json:"image_key"`
	UserContext string `json:"user_context"`
	Country     string `json:"country"`
}

var (
	textractClient   *textractsvc.Client
	comprehendClient *comprehendsvc.Client
	snsClient        *sns.Client
	incidentStore    *store.Store
	rateLimiter      *quota.Limiter
	uploadBucket     string
	alertsTopicArn   string
)

func init() {
	uploadBucket = os.Getenv("UPLOADS_BUCKET")
	if uploadBucket == "" {
		uploadBucket = "sentry-uploads"
	}
	alertsTopicArn = os.Getenv("CAREGIVER_ALERTS_TOPIC_ARN")
	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	linksTable := os.Getenv("LINKS_TABLE")

	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	rateLimiter = quota.NewLimiter()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err == nil {
		textractClient = textractsvc.NewFromConfig(cfg)
		comprehendClient = comprehendsvc.NewFromConfig(cfg)
		snsClient = sns.NewFromConfig(cfg)
		dbClient := dynamodb.NewFromConfig(cfg)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)
	}
}

func handleAnalysis(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		var msg AnalysisQueueMessage
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			return fmt.Errorf("unmarshal sqs message: %w", err)
		}

		var extractedText string
		if msg.ImageKey != "" && textractClient != nil {
			// enforce daily request cap for textract free tier
			allowed, _, _ := rateLimiter.Allow(ctx, "textract", 100)
			if allowed {
				text, err := textract.ExtractText(ctx, textractClient, uploadBucket, msg.ImageKey)
				if err == nil {
					extractedText = text
				}
			}
		}

		combinedText := extractedText
		if msg.UserContext != "" {
			if combinedText != "" {
				combinedText += "\n" + msg.UserContext
			} else {
				combinedText = msg.UserContext
			}
		}

		var compRes *comprehend.AnalysisResult
		if combinedText != "" && comprehendClient != nil {
			// enforce daily request cap for comprehend free tier
			allowed, _, _ := rateLimiter.Allow(ctx, "comprehend", 100)
			if allowed {
				compRes = comprehend.AnalyzeWithFallback(ctx, comprehendClient, combinedText)
			} else {
				compRes = comprehend.FallbackAnalysis(combinedText, "en")
			}
		} else {
			compRes = comprehend.FallbackAnalysis(combinedText, "en")
		}

		input := risk.Input{
			Text:       extractedText,
			Context:    msg.UserContext,
			Sentiment:  compRes.Sentiment,
			KeyPhrases: compRes.KeyPhrases,
			Entities:   compRes.Entities,
		}

		score, level, matches := risk.CalculateScore(input)
		category := risk.ClassifyCategory(input)
		output := risk.BuildSeniorOutput(category, level, msg.Country)

		inc := &model.Incident{
			ID:            msg.IncidentID,
			ImageKey:      msg.ImageKey,
			ExtractedText: extractedText,
			UserContext:   msg.UserContext,
			Score:         score,
			RiskLevel:     level,
			Category:      category,
			Signals:       matches,
			Output:        output,
			Status:        model.StatusAnalyzed,
			Country:       msg.Country,
			UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		}

		if incidentStore != nil {
			if err := incidentStore.SaveIncident(ctx, inc); err != nil {
				return fmt.Errorf("store analyzed incident %s: %w", inc.ID, err)
			}
		}

		if level == model.RiskHigh && snsClient != nil && alertsTopicArn != "" {
			var signalList []string
			for _, m := range matches {
				signalList = append(signalList, fmt.Sprintf("- %s: %s", m.Name, m.Description))
			}

			subject := fmt.Sprintf("[Sentry Alert] High risk scam detected (score %d/100)", score)
			body := fmt.Sprintf(
				"Sentry detected a high risk incident.\n\n"+
					"Incident ID: %s\n"+
					"Risk Score: %d / 100\n"+
					"Category: %s\n"+
					"Time: %s\n\n"+
					"Signals:\n%s\n\n"+
					"Senior Guidance:\n"+
					"Hold off on: %s\n"+
					"Safe alternative: %s\n"+
					"Emergency hotline: %s (%s)\n\n"+
					"Sign in to the caregiver portal to acknowledge or resolve.",
				inc.ID, score, category, inc.UpdatedAt,
				strings.Join(signalList, "\n"),
				output.HoldOff, output.SafeAlternative, output.Hotline.Name, output.Hotline.Number,
			)

			_, _ = snsClient.Publish(ctx, &sns.PublishInput{
				TopicArn: aws.String(alertsTopicArn),
				Subject:  aws.String(subject),
				Message:  aws.String(body),
			})
		}
	}

	return nil
}

func main() {
	lambda.Start(handleAnalysis)
}
