package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/iokuru/sentry/internal/hotline"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/store"
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
	sqsClient        *sqs.Client
	analysisQueueURL string
	incidentStore    *store.Store
)

func init() {
	analysisQueueURL = os.Getenv("ANALYSIS_QUEUE_URL")
	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	linksTable := os.Getenv("LINKS_TABLE")

	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err == nil {
		sqsClient = sqs.NewFromConfig(cfg)
		dbClient := dynamodb.NewFromConfig(cfg)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)
	}
}

func handleIntake(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var input IntakeRequest
	if req.Body != "" {
		if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
			return errorResponse(http.StatusBadRequest, "invalid request json"), nil
		}
	}

	if strings.TrimSpace(input.ImageKey) == "" && strings.TrimSpace(input.UserContext) == "" {
		return errorResponse(http.StatusBadRequest, "either screenshot or context message is required"), nil
	}

	country := input.Country
	if country == "" {
		country = hotline.ResolveCountry(input.Locale)
	}

	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	incidentID := fmt.Sprintf("inc-%d-%s", time.Now().Unix(), hex.EncodeToString(randBytes))

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
			return errorResponse(http.StatusInternalServerError, "failed to record incident"), nil
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
			return errorResponse(http.StatusInternalServerError, "failed to enqueue analysis job"), nil
		}
	}

	res := IntakeResponse{
		IncidentID: incidentID,
		Status:     string(model.StatusPending),
	}
	resBody, _ := json.Marshal(res)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusAccepted,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "POST,OPTIONS",
		},
		Body: string(resBody),
	}, nil
}

func errorResponse(code int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}
}

func main() {
	lambda.Start(handleIntake)
}
