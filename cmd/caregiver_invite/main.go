package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/iokuru/sentry/internal/awsconfig"
	"github.com/iokuru/sentry/internal/logging"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/pairing"
	"github.com/iokuru/sentry/internal/router"
	"github.com/iokuru/sentry/internal/store"
)

type CreateInviteRequest struct {
	SeniorName string `json:"senior_name"`
}

type CreateInviteResponse struct {
	Code      string `json:"code"`
	ExpiresAt int64  `json:"expires_at"`
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

func handleCaregiverInvite(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ctx = logging.WithRequestID(ctx, req.RequestContext.RequestID)
	if linkStore == nil {
		return errorResponse(req, http.StatusInternalServerError, "store unavailable"), nil
	}

	caregiverID := extractCaregiverID(req)
	if caregiverID == "" {
		return errorResponse(req, http.StatusUnauthorized, "unauthorized: valid cognito jwt required"), nil
	}
	method := req.RequestContext.HTTP.Method
	path := req.RequestContext.HTTP.Path

	switch {
	case method == http.MethodPost && strings.HasSuffix(path, "/invites"):
		var input CreateInviteRequest
		if req.Body != "" {
			_ = json.Unmarshal([]byte(req.Body), &input)
		}

		code, err := pairing.GenerateCode()
		if err != nil {
			return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
		}
		inv := &model.Invite{
			Code:        code,
			CaregiverID: caregiverID,
			SeniorName:  strings.TrimSpace(input.SeniorName),
		}

		if err := linkStore.CreateInvite(ctx, inv); err != nil {
			return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
		}

		slog.InfoContext(ctx, "created caregiver invite", slog.String("caregiver_id", caregiverID))

		body, _ := json.Marshal(CreateInviteResponse{
			Code:      code,
			ExpiresAt: inv.ExpiresAt,
		})
		return jsonResponse(req, http.StatusOK, body), nil

	case method == http.MethodGet && strings.HasSuffix(path, "/links"):
		links, err := linkStore.ListLinksForCaregiver(ctx, caregiverID)
		if err != nil {
			return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
		}
		if links == nil {
			links = []model.CaregiverLink{}
		}
		body, _ := json.Marshal(links)
		return jsonResponse(req, http.StatusOK, body), nil

	case method == http.MethodDelete && strings.Contains(path, "/links"):
		seniorID := req.PathParameters["senior_id"]
		if seniorID == "" {
			seniorID = req.QueryStringParameters["senior_id"]
		}
		if seniorID == "" && req.Body != "" {
			var b struct {
				SeniorID string `json:"senior_id"`
			}
			_ = json.Unmarshal([]byte(req.Body), &b)
			seniorID = b.SeniorID
		}
		if seniorID == "" {
			return errorResponse(req, http.StatusBadRequest, "senior_id is required"), nil
		}

		if err := linkStore.DeleteLink(ctx, caregiverID, seniorID); err != nil {
			return errorResponse(req, http.StatusInternalServerError, err.Error()), nil
		}
		return jsonResponse(req, http.StatusOK, []byte(`{"status":"deleted"}`)), nil

	default:
		return errorResponse(req, http.StatusNotFound, "route not found"), nil
	}
}

func corsHeaders(req events.APIGatewayV2HTTPRequest) map[string]string {
	h := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "GET,POST,DELETE,OPTIONS",
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
	lambda.Start(handleCaregiverInvite)
}
