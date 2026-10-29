package logging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/iokuru/sentry/internal/logging"
)

func TestStructuredLoggingWithRequestID(t *testing.T) {
	buf := new(bytes.Buffer)
	logging.SetOutput(buf)

	ctx := logging.WithRequestID(context.Background(), "req-test-one-two-three")
	slog.InfoContext(ctx, "processing event", slog.String("incident_id", "inc-abc-def"))

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse structured json log: %v", err)
	}

	if entry["msg"] != "processing event" {
		t.Errorf("expected msg 'processing event', got %v", entry["msg"])
	}
	if entry["level"] != "INFO" {
		t.Errorf("expected level 'INFO', got %v", entry["level"])
	}
	if entry["request_id"] != "req-test-one-two-three" {
		t.Errorf("expected request_id 'req-test-one-two-three', got %v", entry["request_id"])
	}
	if entry["incident_id"] != "inc-abc-def" {
		t.Errorf("expected incident_id 'inc-abc-def', got %v", entry["incident_id"])
	}
}
