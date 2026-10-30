package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/iokuru/sentry/internal/awsconfig"
	"github.com/iokuru/sentry/internal/logging"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/router"
	"github.com/iokuru/sentry/internal/store"
)

var (
	incidentStore *store.Store
	presignClient *s3.PresignClient
	uploadBucket  string
)

func init() {
	logging.Init()
	uploadBucket = os.Getenv("UPLOADS_BUCKET")
	if uploadBucket == "" {
		uploadBucket = "sentry-uploads"
	}

	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	linksTable := os.Getenv("LINKS_TABLE")
	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	cfg, err := awsconfig.Load(context.Background())
	if err == nil {
		dbClient := dynamodb.NewFromConfig(cfg)
		s3Client := s3.NewFromConfig(cfg)
		presignClient = s3.NewPresignClient(s3Client)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)
	}
}

func attachThumbnailURL(ctx context.Context, inc *model.Incident) {
	if presignClient != nil && inc.ImageKey != "" {
		presigned, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(uploadBucket),
			Key:    aws.String(inc.ImageKey),
		}, s3.WithPresignExpires(15*time.Minute))
		if err == nil {
			inc.ThumbnailURL = presigned.URL
		}
	}
}

func extractCaregiverID(req events.APIGatewayV2HTTPRequest) string {
	if req.RequestContext.Authorizer != nil && req.RequestContext.Authorizer.JWT != nil {
		claims := req.RequestContext.Authorizer.JWT.Claims
		if sub, ok := claims["sub"]; ok && sub != "" {
			return sub
		}
		if username, ok := claims["cognito:username"]; ok && username != "" {
			return username
		}
	}
	if auth := req.Headers["authorization"]; auth != "" || req.Headers["Authorization"] != "" {
		return "caregiver"
	}
	return ""
}

func handleList(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ctx = logging.WithRequestID(ctx, req.RequestContext.RequestID)

	if req.RequestContext.HTTP.Path == "/api/health" || req.RawPath == "/api/health" {
		return jsonResponse(req, http.StatusOK, []byte(`{"status":"healthy","service":"sentry"}`)), nil
	}

	if incidentStore == nil {
		return errorResponse(req, http.StatusInternalServerError, "store unavailable"), nil
	}

	// check if individual incident is queried via query params or path (senior intake polling)
	if id, ok := req.QueryStringParameters["id"]; ok && id != "" {
		inc, err := incidentStore.GetIncident(ctx, id)
		if err != nil {
			return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
		}
		if inc == nil {
			return errorResponse(req, http.StatusNotFound, "incident not found"), nil
		}
		attachThumbnailURL(ctx, inc)
		body, _ := json.Marshal(inc)
		return jsonResponse(req, http.StatusOK, body), nil
	}

	caregiverID := extractCaregiverID(req)
	if caregiverID == "" {
		return errorResponse(req, http.StatusUnauthorized, "unauthorized: valid cognito jwt required"), nil
	}

	links, err := incidentStore.ListLinksForCaregiver(ctx, caregiverID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to query caregiver links", slog.String("error", err.Error()))
		return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
	}

	linkedSeniors := make(map[string]bool)
	for _, l := range links {
		linkedSeniors[l.SeniorID] = true
	}

	allIncidents, err := incidentStore.ListIncidents(ctx, 50)
	if err != nil {
		return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
	}

	var incidents []model.Incident
	for _, inc := range allIncidents {
		if len(linkedSeniors) == 0 || linkedSeniors[inc.SeniorID] {
			attachThumbnailURL(ctx, &inc)
			incidents = append(incidents, inc)
		}
	}
	if incidents == nil {
		incidents = []model.Incident{}
	}

	slog.InfoContext(ctx, "listed incidents for caregiver", slog.String("caregiver_id", caregiverID), slog.Int("count", len(incidents)))

	body, _ := json.Marshal(incidents)
	return jsonResponse(req, http.StatusOK, body), nil
}

func corsHeaders(req events.APIGatewayV2HTTPRequest) map[string]string {
	h := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "GET,OPTIONS",
	}
	allowed := os.Getenv("CORS_ORIGIN")
	if origin := router.MatchOrigin(allowed, router.RequestOrigin(req.Headers)); origin != "" {
		h["Access-Control-Allow-Origin"] = origin
	}
	return h
}

func jsonResponse(req events.APIGatewayV2HTTPRequest, code int, body []byte) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers:    corsHeaders(req),
		Body:       string(body),
	}
}

func errorResponse(req events.APIGatewayV2HTTPRequest, code int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return jsonResponse(req, code, body)
}

func main() {
	lambda.Start(handleList)
}
