package store

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/iokuru/sentry/internal/model"
)

type mockTransport struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.handler(req)
}

func TestStoreOperations(t *testing.T) {
	ctx := context.Background()

	mock := &mockTransport{
		handler: func(req *http.Request) (*http.Response, error) {
			// return successful empty json responses for put and update
			body := "{}"
			header := req.Header.Get("X-Amz-Target")
			if header == "DynamoDB_20120810.GetItem" {
				body = `{"Item":{"id":{"S":"inc-123"},"status":{"S":"analyzed"},"score":{"N":"75"},"risk_level":{"S":"high"}}}`
			} else if header == "DynamoDB_20120810.Scan" {
				body = `{"Items":[{"id":{"S":"inc-1"},"status":{"S":"pending"},"created_at":{"S":"2026-10-22T10:00:00Z"}},{"id":{"S":"inc-2"},"status":{"S":"resolved"},"created_at":{"S":"2026-10-22T11:00:00Z"}}]}`
			} else if header == "DynamoDB_20120810.Query" {
				body = `{"Items":[{"caregiver_id":{"S":"cg-1"},"senior_id":{"S":"sen-1"},"senior_name":{"S":"Mom"},"telegram_chat_id":{"S":"12345"}}]}`
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}, nil
		},
	}

	cfg := aws.Config{
		Region: "us-east-1",
		HTTPClient: &http.Client{
			Transport: mock,
		},
	}
	client := dynamodb.NewFromConfig(cfg)
	st := NewStore(client, "sentry-incidents", "sentry-links")

	t.Run("save incident", func(t *testing.T) {
		inc := &model.Incident{
			ID:        "inc-123",
			Score:     75,
			RiskLevel: model.RiskHigh,
			Category:  model.CategoryBank,
			Status:    model.StatusPending,
		}
		if err := st.SaveIncident(ctx, inc); err != nil {
			t.Fatalf("SaveIncident() failed: %v", err)
		}
	})

	t.Run("get incident", func(t *testing.T) {
		inc, err := st.GetIncident(ctx, "inc-123")
		if err != nil {
			t.Fatalf("GetIncident() failed: %v", err)
		}
		if inc == nil || inc.ID != "inc-123" {
			t.Fatalf("expected inc-123, got %+v", inc)
		}
	})

	t.Run("list unresolved", func(t *testing.T) {
		unresolved, err := st.ListUnresolvedIncidents(ctx)
		if err != nil {
			t.Fatalf("ListUnresolvedIncidents() failed: %v", err)
		}
		if len(unresolved) != 1 || unresolved[0].ID != "inc-1" {
			t.Fatalf("expected 1 unresolved incident (inc-1), got %d items", len(unresolved))
		}
	})

	t.Run("update status", func(t *testing.T) {
		if err := st.UpdateIncidentStatus(ctx, "inc-123", model.StatusResolved); err != nil {
			t.Fatalf("UpdateIncidentStatus() failed: %v", err)
		}
	})

	t.Run("create and get invite", func(t *testing.T) {
		inv := &model.Invite{
			Code:        "CDFGHJ",
			CaregiverID: "cg-1",
			SeniorName:  "Mom",
		}
		if err := st.CreateInvite(ctx, inv); err != nil {
			t.Fatalf("CreateInvite() failed: %v", err)
		}
	})

	t.Run("link operations", func(t *testing.T) {
		link := &model.CaregiverLink{
			CaregiverID: "cg-1",
			SeniorID:    "sen-1",
			SeniorName:  "Mom",
		}
		if err := st.CreateLink(ctx, link); err != nil {
			t.Fatalf("CreateLink() failed: %v", err)
		}
		links, err := st.ListLinksForCaregiver(ctx, "cg-1")
		if err != nil {
			t.Fatalf("ListLinksForCaregiver() failed: %v", err)
		}
		if len(links) != 1 || links[0].SeniorID != "sen-1" {
			t.Fatalf("expected 1 link, got %+v", links)
		}
		seniorLinks, err := st.ListLinksForSenior(ctx, "sen-1")
		if err != nil {
			t.Fatalf("ListLinksForSenior() failed: %v", err)
		}
		if len(seniorLinks) != 1 {
			t.Fatalf("expected 1 senior link, got %+v", seniorLinks)
		}
		if err := st.SetTelegramChatID(ctx, "cg-1", "99999"); err != nil {
			t.Fatalf("SetTelegramChatID() failed: %v", err)
		}
		if err := st.DeleteLink(ctx, "cg-1", "sen-1"); err != nil {
			t.Fatalf("DeleteLink() failed: %v", err)
		}
	})
}
