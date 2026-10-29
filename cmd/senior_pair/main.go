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
		return errorResponse(http.StatusInternalServerError, "store unavailable"), nil
	}

	var input PairRequest
	if req.Body != "" {
		if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
			return errorResponse(http.StatusBadRequest, "invalid request json"), nil
		}
	}

	code := strings.ToUpper(strings.TrimSpace(input.Code))
	deviceID := strings.TrimSpace(input.DeviceID)

	if code == "" || deviceID == "" {
		return errorResponse(http.StatusBadRequest, "code and device_id are required"), nil
	}

	link, err := linkStore.RedeemInvite(ctx, code, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrInviteNotFound):
			return errorResponse(http.StatusNotFound, "invite code not found"), nil
		case errors.Is(err, store.ErrInviteExpired):
			return errorResponse(http.StatusGone, "invite code has expired"), nil
		case errors.Is(err, store.ErrInviteUsed):
			return errorResponse(http.StatusConflict, "invite code has already been used"), nil
		default:
			return errorResponse(http.StatusInternalServerError, err.Error()), nil
		}
	}

	slog.InfoContext(ctx, "senior paired successfully", slog.String("senior_id", link.SeniorID))

	body, _ := json.Marshal(PairResponse{
		SeniorID:   link.SeniorID,
		SeniorName: link.SeniorName,
	})
	return jsonResponse(http.StatusOK, body), nil
}

func getCorsOrigin() string {
	if o := os.Getenv("CORS_ORIGIN"); o != "" {
		return o
	}
	return "*"
}

func jsonResponse(status int, body []byte) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  getCorsOrigin(),
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
		},
		Body: string(body),
	}
}

func errorResponse(status int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return jsonResponse(status, body)
}

func main() {
	lambda.Start(handleSeniorPair)
}
