package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestTextOnlyPath(t *testing.T) {
	oldStore := incidentStore
	oldSQS := sqsClient
	defer func() {
		incidentStore = oldStore
		sqsClient = oldSQS
	}()
	incidentStore = nil
	sqsClient = nil

	reqBody := `{"senior_id":"dev-123","user_context":"urgent IRS demand for payment via gift cards","country":"US"}`
	req := events.APIGatewayV2HTTPRequest{
		Body: reqBody,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "test-req-1",
		},
	}

	res, err := handleIntake(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d: %s", res.StatusCode, res.Body)
	}

	var resp IntakeResponse
	if err := json.Unmarshal([]byte(res.Body), &resp); err != nil {
		t.Fatalf("failed to parse response json: %v", err)
	}

	if !strings.HasPrefix(resp.IncidentID, "inc-") {
		t.Fatalf("expected incident id to have prefix inc-, got %s", resp.IncidentID)
	}

	if resp.Status != "pending" {
		t.Fatalf("expected status pending, got %s", resp.Status)
	}
}

func TestTypedErrorHandoffWhenTextractBlocked(t *testing.T) {
	oldMode := ocrMode
	oldBlocked := textractBlocked
	defer func() {
		ocrMode = oldMode
		textractBlocked = oldBlocked
	}()

	ocrMode = "browser"
	textractBlocked = true

	reqBody := `{"senior_id":"dev-123","image_key":"screenshots/test.png","country":"US"}`
	req := events.APIGatewayV2HTTPRequest{
		Body: reqBody,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "test-req-2",
		},
	}

	res, err := handleIntake(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d: %s", res.StatusCode, res.Body)
	}

	if !strings.Contains(res.Body, "TEXTRACT_UNAVAILABLE") {
		t.Fatalf("expected response body to contain TEXTRACT_UNAVAILABLE, got: %s", res.Body)
	}
}

func TestBrowserOCRSubmissionSucceeds(t *testing.T) {
	oldMode := ocrMode
	oldBlocked := textractBlocked
	oldStore := incidentStore
	oldSQS := sqsClient
	defer func() {
		ocrMode = oldMode
		textractBlocked = oldBlocked
		incidentStore = oldStore
		sqsClient = oldSQS
	}()

	ocrMode = "browser"
	textractBlocked = true
	incidentStore = nil
	sqsClient = nil

	// client sends image_key for thumbnail AND extracted text in user_context
	reqBody := `{"senior_id":"dev-123","image_key":"screenshots/test.png","user_context":"text extracted on device via tesseract worker","country":"US"}`
	req := events.APIGatewayV2HTTPRequest{
		Body: reqBody,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "test-req-3",
		},
	}

	res, err := handleIntake(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d: %s", res.StatusCode, res.Body)
	}

	var resp IntakeResponse
	if err := json.Unmarshal([]byte(res.Body), &resp); err != nil {
		t.Fatalf("failed to parse response json: %v", err)
	}

	if !strings.HasPrefix(resp.IncidentID, "inc-") {
		t.Fatalf("expected incident id with inc- prefix, got %s", resp.IncidentID)
	}
}
