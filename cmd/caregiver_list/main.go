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
	"github.com/iokuru/sentry/internal/store"
)

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

func handleList(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// check if individual incident is queried via query params or path
	if id, ok := req.QueryStringParameters["id"]; ok && id != "" {
		if incidentStore == nil {
			return errorResponse(http.StatusInternalServerError, "store unavailable"), nil
		}
		inc, err := incidentStore.GetIncident(ctx, id)
		if err != nil {
			return errorResponse(http.StatusInternalServerError, err.Error()), nil
		}
		if inc == nil {
			return errorResponse(http.StatusNotFound, "incident not found"), nil
		}
		body, _ := json.Marshal(inc)
		return jsonResponse(http.StatusOK, body), nil
	}

	if incidentStore == nil {
		return errorResponse(http.StatusInternalServerError, "store unavailable"), nil
	}

	incidents, err := incidentStore.ListIncidents(ctx, 50)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}

	body, _ := json.Marshal(incidents)
	return jsonResponse(http.StatusOK, body), nil
}

func jsonResponse(code int, body []byte) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "GET,OPTIONS",
		},
		Body: string(body),
	}
}

func errorResponse(code int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return jsonResponse(code, body)
}

func main() {
	lambda.Start(handleList)
}
