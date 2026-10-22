package quota

import (
	"context"
	"testing"
)

func TestLimiterCap(t *testing.T) {
	ctx := context.Background()
	lim := NewLimiter()

	cap := 3
	for i := 1; i <= cap; i++ {
		allowed, count, err := lim.Allow(ctx, "textract", cap)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed || count != i {
			t.Fatalf("iteration %d: allowed=%v count=%d", i, allowed, count)
		}
	}

	// next request should be blocked by daily cap
	allowed, count, err := lim.Allow(ctx, "textract", cap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatalf("expected request to be rejected past cap, got allowed=true, count=%d", count)
	}
}
