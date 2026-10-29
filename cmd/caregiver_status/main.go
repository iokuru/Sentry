package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/iokuru/sentry/internal/logging"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/store"
)

type StatusUpdateRequest struct {
	IncidentID string               `json:"incident_id"`
	Status     model.IncidentStatus `json:"status"`
}

var incidentStore *store.Store

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

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err == nil {
		dbClient := dynamodb.NewFromConfig(cfg)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)
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

func handleStatusUpdate(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ctx = logging.WithRequestID(ctx, req.RequestContext.RequestID)

	caregiverID := extractCaregiverID(req)
	if caregiverID == "" {
		return errorResponse(http.StatusUnauthorized, "unauthorized: valid cognito jwt required"), nil
	}

	var input StatusUpdateRequest
	if req.Body != "" {
		if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
			return errorResponse(http.StatusBadRequest, "invalid request json"), nil
		}
	}

	if id, ok := req.PathParameters["id"]; ok && id != "" {
		input.IncidentID = id
	}

	if input.IncidentID == "" {
		return errorResponse(http.StatusBadRequest, "incident_id is required"), nil
	}

	if input.Status != model.StatusAcknowledged && input.Status != model.StatusResolved && input.Status != model.StatusAnalyzed {
		return errorResponse(http.StatusBadRequest, "status must be acknowledged, resolved, or analyzed"), nil
	}

	if incidentStore == nil {
		return errorResponse(http.StatusInternalServerError, "store unavailable"), nil
	}

	inc, err := incidentStore.GetIncident(ctx, input.IncidentID)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if inc == nil {
		return errorResponse(http.StatusNotFound, "incident not found"), nil
	}

	if inc.SeniorID != "" {
		links, err := incidentStore.ListLinksForCaregiver(ctx, caregiverID)
		if err != nil {
			return errorResponse(http.StatusInternalServerError, err.Error()), nil
		}
		isLinked := false
		for _, l := range links {
			if l.SeniorID == inc.SeniorID {
				isLinked = true
				break
			}
		}
		if !isLinked {
			return errorResponse(http.StatusForbidden, "forbidden: senior not linked to caregiver"), nil
		}
	}

	if err := incidentStore.UpdateIncidentStatus(ctx, input.IncidentID, input.Status); err != nil {
		slog.ErrorContext(ctx, "failed to update status", slog.String("incident_id", input.IncidentID), slog.String("error", err.Error()))
		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}

	slog.InfoContext(ctx, "status updated", slog.String("incident_id", input.IncidentID), slog.String("status", string(input.Status)))

	resBody, _ := json.Marshal(map[string]string{
		"incident_id": input.IncidentID,
		"status":      string(input.Status),
	})

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  getCorsOrigin(),
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "POST,PATCH,OPTIONS",
		},
		Body: string(resBody),
	}, nil
}

func getCorsOrigin() string {
	if o := os.Getenv("CORS_ORIGIN"); o != "" {
		return o
	}
	return "*"
}

func errorResponse(code int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": getCorsOrigin(),
		},
		Body: string(body),
	}
}

func main() {
	lambda.Start(handleStatusUpdate)
}
