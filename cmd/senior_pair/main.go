package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/iokuru/sentry/internal/awsconfig"
	"github.com/iokuru/sentry/internal/logging"
	"github.com/iokuru/sentry/internal/router"
	"github.com/iokuru/sentry/internal/store"
)

type PairRequest struct {
	Code     string `json:"code"`
	DeviceID string `json:"device_id"`
}

type PairResponse struct {
	SeniorID   string `json:"senior_id"`
	SeniorName string `json:"senior_name"`
}

var linkStore *store.Store

func init() {
	logging.Init()
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
		linkStore = store.NewStore(dbClient, incidentsTable, linksTable)
	}
}

func handleSeniorPair(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ctx = logging.WithRequestID(ctx, req.RequestContext.RequestID)
	if linkStore == nil {
		return errorResponse(req, http.StatusInternalServerError, "store unavailable"), nil
	}

	var input PairRequest
	if req.Body != "" {
		if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
			return errorResponse(req, http.StatusBadRequest, "invalid request json"), nil
		}
	}

	code := strings.ToUpper(strings.TrimSpace(input.Code))
	deviceID := strings.TrimSpace(input.DeviceID)

	if code == "" || deviceID == "" {
		return errorResponse(req, http.StatusBadRequest, "code and device_id are required"), nil
	}

	link, err := linkStore.RedeemInvite(ctx, code, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrInviteNotFound):
			return errorResponse(req, http.StatusNotFound, "invite code not found"), nil
		case errors.Is(err, store.ErrInviteExpired):
			return errorResponse(req, http.StatusGone, "invite code has expired"), nil
		case errors.Is(err, store.ErrInviteUsed):
			return errorResponse(req, http.StatusConflict, "invite code has already been used"), nil
		default:
			return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
		}
	}

	slog.InfoContext(ctx, "senior paired successfully", slog.String("senior_id", link.SeniorID))

	body, _ := json.Marshal(PairResponse{
		SeniorID:   link.SeniorID,
		SeniorName: link.SeniorName,
	})
	return jsonResponse(req, http.StatusOK, body), nil
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

func jsonResponse(req events.APIGatewayV2HTTPRequest, status int, body []byte) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers:    corsHeaders(req),
		Body:       string(body),
	}
}

func errorResponse(req events.APIGatewayV2HTTPRequest, status int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return jsonResponse(req, status, body)
}

func main() {
	lambda.Start(handleSeniorPair)
}
