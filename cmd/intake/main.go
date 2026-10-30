package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	textractsvc "github.com/aws/aws-sdk-go-v2/service/textract"
	textracttypes "github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/iokuru/sentry/internal/awsconfig"
	"github.com/iokuru/sentry/internal/hotline"
	"github.com/iokuru/sentry/internal/logging"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/router"
	"github.com/iokuru/sentry/internal/store"
	"github.com/iokuru/sentry/internal/textract"
)

type IntakeRequest struct {
	ImageKey    string `json:"image_key"`
	UserContext string `json:"user_context"`
	Locale      string `json:"locale"`
	Country     string `json:"country"`
	SeniorID    string `json:"senior_id"`
}

type IntakeResponse struct {
	IncidentID string `json:"incident_id"`
	Status     string `json:"status"`
}

type AnalysisQueueMessage struct {
	IncidentID  string `json:"incident_id"`
	SeniorID    string `json:"senior_id,omitempty"`
	ImageKey    string `json:"image_key"`
	UserContext string `json:"user_context"`
	Country     string `json:"country"`
}

var (
	sqsClient         *sqs.Client
	textractClient    *textractsvc.Client
	analysisQueueURL  string
	incidentStore     *store.Store
	ocrMode           string
	textractBlocked   bool
	textractProbeOnce sync.Once
)

var onePixelPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
	0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

func init() {
	logging.Init()
	analysisQueueURL = os.Getenv("ANALYSIS_QUEUE_URL")
	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	linksTable := os.Getenv("LINKS_TABLE")
	ocrMode = strings.ToLower(os.Getenv("OCR_MODE"))
	if ocrMode == "" {
		ocrMode = "auto"
	}
	if ocrMode == "browser" {
		textractBlocked = true
	}

	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	cfg, err := awsconfig.Load(context.Background())
	if err == nil {
		sqsClient = sqs.NewFromConfig(cfg)
		dbClient := dynamodb.NewFromConfig(cfg)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)
		if ocrMode != "browser" {
			textractClient = textractsvc.NewFromConfig(cfg)
		}
	}
}

func isTextractBlocked(ctx context.Context) bool {
	if ocrMode == "browser" {
		return true
	}
	if ocrMode == "textract" {
		return false
	}
	textractProbeOnce.Do(func() {
		if textractClient == nil {
			textractBlocked = true
			return
		}
		probeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		_, err := textractClient.DetectDocumentText(probeCtx, &textractsvc.DetectDocumentTextInput{
			Document: &textracttypes.Document{
				Bytes: onePixelPNG,
			},
		})
		if err != nil && textract.IsAccessOrSubscriptionError(err) {
			textractBlocked = true
			slog.WarnContext(ctx, "textract probe failed, switching to browser ocr handoff", slog.String("error", err.Error()))
		}
	})
	return textractBlocked
}

func handleIntake(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ctx = logging.WithRequestID(ctx, req.RequestContext.RequestID)
	var input IntakeRequest
	if req.Body != "" {
		if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
			return errorResponse(req, http.StatusBadRequest, "invalid request json"), nil
		}
	}

	if strings.TrimSpace(input.SeniorID) == "" {
		return errorResponse(req, http.StatusBadRequest, "device id is required"), nil
	}

	if strings.TrimSpace(input.ImageKey) == "" && strings.TrimSpace(input.UserContext) == "" {
		return errorResponse(req, http.StatusBadRequest, "either screenshot or context message is required"), nil
	}

	if strings.TrimSpace(input.ImageKey) != "" && strings.TrimSpace(input.UserContext) == "" {
		if isTextractBlocked(ctx) {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Headers:    corsHeaders(req),
				Body:       `{"code":"TEXTRACT_UNAVAILABLE","error":"textract unavailable: use browser ocr"}`,
			}, nil
		}
	}

	country := input.Country
	if country == "" {
		country = hotline.ResolveCountry(input.Locale)
	}

	var incidentID string
	if strings.TrimSpace(input.ImageKey) != "" {
		hash := sha256.Sum256([]byte(strings.TrimSpace(input.ImageKey)))
		incidentID = fmt.Sprintf("inc-%s", hex.EncodeToString(hash[:10]))
	} else {
		randBytes := make([]byte, 8)
		_, _ = rand.Read(randBytes)
		incidentID = fmt.Sprintf("inc-%d-%s", time.Now().Unix(), hex.EncodeToString(randBytes))
	}

	if incidentStore != nil {
		existing, err := incidentStore.GetIncident(ctx, incidentID)
		if err == nil && existing != nil {
			res := IntakeResponse{
				IncidentID: existing.ID,
				Status:     string(existing.Status),
			}
			resBody, _ := json.Marshal(res)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusOK,
				Headers:    corsHeaders(req),
				Body:       string(resBody),
			}, nil
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	inc := &model.Incident{
		ID:          incidentID,
		SeniorID:    input.SeniorID,
		ImageKey:    input.ImageKey,
		UserContext: input.UserContext,
		Status:      model.StatusPending,
		Country:     country,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if incidentStore != nil {
		if err := incidentStore.SaveIncident(ctx, inc); err != nil {
			return errorResponse(req, http.StatusInternalServerError, "failed to record incident"), nil
		}
	}

	if sqsClient != nil && analysisQueueURL != "" {
		msg := AnalysisQueueMessage{
			IncidentID:  incidentID,
			SeniorID:    input.SeniorID,
			ImageKey:    input.ImageKey,
			UserContext: input.UserContext,
			Country:     country,
		}
		msgBody, _ := json.Marshal(msg)
		_, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(analysisQueueURL),
			MessageBody: aws.String(string(msgBody)),
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to enqueue analysis job", slog.String("incident_id", incidentID), slog.String("error", err.Error()))
			return errorResponse(req, http.StatusInternalServerError, "failed to enqueue analysis job"), nil
		}
	}

	slog.InfoContext(ctx, "incident queued for analysis", slog.String("incident_id", incidentID), slog.String("country", country))

	res := IntakeResponse{
		IncidentID: incidentID,
		Status:     string(model.StatusPending),
	}
	resBody, _ := json.Marshal(res)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusAccepted,
		Headers:    corsHeaders(req),
		Body:       string(resBody),
	}, nil
}

func corsHeaders(req events.APIGatewayV2HTTPRequest) map[string]string {
	h := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "POST,OPTIONS",
	}
	allowed := os.Getenv("CORS_ORIGIN")
	if origin := router.MatchOrigin(allowed, router.RequestOrigin(req.Headers)); origin != "" {
		h["Access-Control-Allow-Origin"] = origin
	}
	return h
}

func errorResponse(req events.APIGatewayV2HTTPRequest, code int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers:    corsHeaders(req),
		Body:       string(body),
	}
}

func main() {
	lambda.Start(handleIntake)
}
