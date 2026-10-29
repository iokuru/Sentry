package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	comprehendsvc "github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	textractsvc "github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/iokuru/sentry/internal/awsconfig"
	"github.com/iokuru/sentry/internal/comprehend"
	"github.com/iokuru/sentry/internal/logging"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/quota"
	"github.com/iokuru/sentry/internal/risk"
	"github.com/iokuru/sentry/internal/store"
	"github.com/iokuru/sentry/internal/telegram"
	"github.com/iokuru/sentry/internal/textract"
)

type AnalysisQueueMessage struct {
	IncidentID  string `json:"incident_id"`
	SeniorID    string `json:"senior_id,omitempty"`
	ImageKey    string `json:"image_key"`
	UserContext string `json:"user_context"`
	Country     string `json:"country"`
}

var (
	textractClient     *textractsvc.Client
	comprehendClient   *comprehendsvc.Client
	snsClient          *sns.Client
	sqsClient          *sqs.Client
	s3Client           *s3.Client
	telegramClient     *telegram.Client
	telegramCfg        *telegram.Config
	incidentStore      *store.Store
	rateLimiter        *quota.Limiter
	uploadBucket       string
	alertsTopicArn     string
	analysisDLQURL     string
	useComprehend      string
	comprehendDisabled bool
	comprehendLogOnce  sync.Once
)

func init() {
	logging.Init()
	uploadBucket = os.Getenv("UPLOADS_BUCKET")
	if uploadBucket == "" {
		uploadBucket = "sentry-uploads"
	}
	alertsTopicArn = os.Getenv("CAREGIVER_ALERTS_TOPIC_ARN")
	analysisDLQURL = os.Getenv("ANALYSIS_DLQ_URL")
	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	linksTable := os.Getenv("LINKS_TABLE")

	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	useComprehend = strings.ToLower(os.Getenv("USE_COMPREHEND"))
	if useComprehend == "" {
		useComprehend = "auto"
	}
	if useComprehend == "off" {
		comprehendDisabled = true
	}

	rateLimiter = quota.NewLimiter()
	telegramClient = telegram.NewClient()

	cfg, err := awsconfig.Load(context.Background())
	if err == nil {
		textractClient = textractsvc.NewFromConfig(cfg)
		comprehendClient = comprehendsvc.NewFromConfig(cfg)
		snsClient = sns.NewFromConfig(cfg)
		sqsClient = sqs.NewFromConfig(cfg)
		s3Client = s3.NewFromConfig(cfg)
		dbClient := dynamodb.NewFromConfig(cfg)
		ssmClient := ssm.NewFromConfig(cfg)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)

		// cached telegram config loaded once at cold start
		telegramCfg = telegram.LoadConfigFromSSM(context.Background(), ssmClient)
	}
}

func forwardPoisonToDLQ(ctx context.Context, body string, reason string) {
	if sqsClient != nil && analysisDLQURL != "" {
		_, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(analysisDLQURL),
			MessageBody: aws.String(body),
			MessageAttributes: map[string]sqstypes.MessageAttributeValue{
				"dlq_reason": {
					DataType:    aws.String("String"),
					StringValue: aws.String(reason),
				},
			},
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to forward poison message to dlq", slog.String("error", err.Error()))
		}
	}
}

func handleAnalysis(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		recCtx := logging.WithRequestID(ctx, record.MessageId)
		var msg AnalysisQueueMessage
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			slog.WarnContext(recCtx, "corrupt poison sqs message", slog.String("error", err.Error()))
			forwardPoisonToDLQ(recCtx, record.Body, fmt.Sprintf("unmarshal error: %v", err))
			continue
		}

		if strings.TrimSpace(msg.ImageKey) == "" && strings.TrimSpace(msg.UserContext) == "" {
			slog.WarnContext(recCtx, "invalid sqs message: empty image and context", slog.String("incident_id", msg.IncidentID))
			forwardPoisonToDLQ(recCtx, record.Body, "empty image key and user context")
			if incidentStore != nil && msg.IncidentID != "" {
				_ = incidentStore.UpdateIncidentStatus(recCtx, msg.IncidentID, model.StatusFailed)
			}
			continue
		}

		if incidentStore != nil && msg.IncidentID != "" {
			existing, err := incidentStore.GetIncident(recCtx, msg.IncidentID)
			if err == nil && existing != nil {
				if existing.Status == model.StatusAnalyzed || existing.Status == model.StatusAcknowledged || existing.Status == model.StatusResolved {
					slog.InfoContext(recCtx, "incident already processed, skipping retry", slog.String("incident_id", existing.ID), slog.String("status", string(existing.Status)))
					continue
				}
			}
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
		if combinedText != "" && comprehendClient != nil && !comprehendDisabled {
			// enforce daily request cap for comprehend free tier
			allowed, _, _ := rateLimiter.Allow(ctx, "comprehend", 100)
			if allowed {
				res, err := comprehend.Analyze(ctx, comprehendClient, combinedText)
				if err != nil {
					if comprehend.IsAccessOrSubscriptionError(err) {
						comprehendDisabled = true
						comprehendLogOnce.Do(func() {
							slog.WarnContext(recCtx, "comprehend disabled due to subscription or access error, falling back to rules", slog.String("error", err.Error()))
						})
					} else {
						slog.WarnContext(recCtx, "comprehend call failed, falling back to rules", slog.String("error", err.Error()))
					}
					compRes = comprehend.FallbackAnalysis(combinedText, "en")
				} else {
					if res.Language != "en" {
						compRes = comprehend.FallbackAnalysis(combinedText, res.Language)
					} else {
						compRes = res
					}
				}
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
			SeniorID:      msg.SeniorID,
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
			if err := incidentStore.SaveIncident(recCtx, inc); err != nil {
				return fmt.Errorf("store analyzed incident %s: %w", inc.ID, err)
			}
		}

		if s3Client != nil && msg.ImageKey != "" && uploadBucket != "" {
			_, err := s3Client.DeleteObject(recCtx, &s3.DeleteObjectInput{
				Bucket: aws.String(uploadBucket),
				Key:    aws.String(msg.ImageKey),
			})
			if err != nil {
				slog.WarnContext(recCtx, "failed to delete analyzed screenshot", slog.String("image_key", msg.ImageKey), slog.String("error", err.Error()))
			}
		}

		slog.InfoContext(recCtx, "analysis completed",
			slog.String("incident_id", inc.ID),
			slog.Int("score", score),
			slog.String("risk_level", string(level)),
			slog.String("category", string(category)),
		)

		if level == model.RiskHigh {
			notifyCaregivers(recCtx, inc, score, category, matches, output)
		}
	}

	return nil
}

func notifyCaregivers(ctx context.Context, inc *model.Incident, score int, category model.Category, matches []model.SignalMatch, output model.SeniorOutput) {
	var wg sync.WaitGroup

	if snsClient != nil && alertsTopicArn != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
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

			_, err := snsClient.Publish(ctx, &sns.PublishInput{
				TopicArn: aws.String(alertsTopicArn),
				Subject:  aws.String(subject),
				Message:  aws.String(body),
			})
			if err != nil {
				slog.ErrorContext(ctx, "sns alert failed", slog.String("error", err.Error()))
			}
		}()
	}

	if telegramCfg != nil && telegramCfg.Enabled && telegramClient != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			detectedOneLine := "Suspicious patterns detected"
			if len(matches) > 0 {
				detectedOneLine = matches[0].Description
			}

			text, buttons := telegram.FormatAlertMessage(
				string(category),
				string(inc.RiskLevel),
				detectedOneLine,
				output.Hotline.Name,
				output.Hotline.Number,
				inc.ID,
				inc.UpdatedAt,
			)

			targetChatIDs := make(map[string]struct{})
			if incidentStore != nil && inc.SeniorID != "" {
				links, err := incidentStore.ListLinksForSenior(ctx, inc.SeniorID)
				if err == nil {
					for _, l := range links {
						if l.TelegramChatID != "" {
							targetChatIDs[l.TelegramChatID] = struct{}{}
						}
					}
				}
			}

			if len(targetChatIDs) == 0 && telegramCfg.ChatID != "" {
				targetChatIDs[telegramCfg.ChatID] = struct{}{}
			}

			for chatID := range targetChatIDs {
				_, err := telegramClient.Send(ctx, telegramCfg.Token, chatID, text, buttons)
				if err != nil {
					slog.ErrorContext(ctx, "telegram alert failed", slog.String("chat_id", chatID), slog.String("error", err.Error()))
				}
			}
		}()
	}

	wg.Wait()
}

func main() {
	lambda.Start(handleAnalysis)
}
