package awsconfig

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Load returns an AWS SDK configuration with bounded retries and exponential jitter backoff.
func Load(ctx context.Context) (aws.Config, error) {
	loadCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return config.LoadDefaultConfig(loadCtx,
		config.WithRetryer(func() aws.Retryer {
			return retry.NewStandard(func(o *retry.StandardOptions) {
				o.MaxAttempts = 3
				o.Backoff = retry.NewExponentialJitterBackoff(200 * time.Millisecond)
			})
		}),
	)
}
