package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/iokuru/sentry/internal/telegram"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		// retrieve deployed api url from terraform output
		cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", "api_url")
		out, err := cmd.Output()
		if err == nil {
			apiUrl = strings.TrimSpace(string(out))
		}
	}

	if apiUrl == "" {
		fmt.Fprintf(os.Stderr, "api_url not found. Deploy infra first or set API_URL environment variable.\n")
		os.Exit(1)
	}

	webhookUrl := strings.TrimSuffix(apiUrl, "/") + "/telegram/hook"

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load aws config: %v\n", err)
		os.Exit(1)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	tgCfg := telegram.LoadConfigFromSSM(ctx, ssmClient)

	if !tgCfg.Enabled {
		fmt.Fprintf(os.Stderr, "telegram is not configured in SSM /sentry/telegram_*. Set token, chat_id, and webhook_secret in SSM first.\n")
		os.Exit(1)
	}

	reqPayload := map[string]string{
		"url":          webhookUrl,
		"secret_token": tgCfg.WebhookSecret,
	}
	body, _ := json.Marshal(reqPayload)

	tgEndpoint := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", tgCfg.Token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tgEndpoint, bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "create request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "setWebhook call failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("Telegram setWebhook response: %s\n", string(respBytes))
}
