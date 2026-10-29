package awsconfig_test

import (
	"context"
	"testing"

	"github.com/iokuru/sentry/internal/awsconfig"
)

func TestLoad(t *testing.T) {
	cfg, err := awsconfig.Load(context.Background())
	if err != nil {
		t.Fatalf("failed to load aws config: %v", err)
	}

	if cfg.Retryer == nil {
		t.Errorf("expected custom retryer to be configured, got nil")
	}
}
