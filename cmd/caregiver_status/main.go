package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/store"
)

type StatusUpdateRequest struct {
	IncidentID string               `json:"incident_id"`
	Status     model.IncidentStatus `json:"status"`
}

var incidentStore *store.Store

func init() {
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

func handleStatusUpdate(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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

	if input.Status != model.StatusAcknowledged && input.Status != model.StatusResolved {
		return errorResponse(http.StatusBadRequest, "status must be acknowledged or resolved"), nil
	}

	if incidentStore == nil {
		return errorResponse(http.StatusInternalServerError, "store unavailable"), nil
	}

	if err := incidentStore.UpdateIncidentStatus(ctx, input.IncidentID, input.Status); err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}

	resBody, _ := json.Marshal(map[string]string{
		"incident_id": input.IncidentID,
		"status":      string(input.Status),
	})

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "POST,PATCH,OPTIONS",
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
	lambda.Start(handleStatusUpdate)
}
