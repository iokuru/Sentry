package telegram

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Config struct {
	Token         string
	ChatID        string
	WebhookSecret string
	Enabled       bool
}

func LoadConfigFromSSM(ctx context.Context, ssmClient *ssm.Client) *Config {
	prefix := os.Getenv("SSM_PREFIX")
	if prefix == "" {
		prefix = "/sentry"
	}
	prefix = strings.TrimSuffix(prefix, "/")

	tokenParam := prefix + "/telegram_token"
	chatParam := prefix + "/telegram_chat_id"
	secretParam := prefix + "/telegram_webhook_secret"

	cfg := &Config{Enabled: false}

	if ssmClient == nil {
		return cfg
	}

	out, err := ssmClient.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          []string{tokenParam, chatParam, secretParam},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return cfg
	}

	for _, p := range out.Parameters {
		val := aws.ToString(p.Value)
		switch aws.ToString(p.Name) {
		case tokenParam:
			cfg.Token = val
		case chatParam:
			cfg.ChatID = val
		case secretParam:
			cfg.WebhookSecret = val
		}
	}

	if isValidValue(cfg.Token) && isValidValue(cfg.ChatID) && isValidValue(cfg.WebhookSecret) {
		cfg.Enabled = true
	}

	return cfg
}

func isValidValue(v string) bool {
	trimmed := strings.TrimSpace(v)
	return trimmed != "" && !strings.EqualFold(trimmed, "PLACEHOLDER")
}
