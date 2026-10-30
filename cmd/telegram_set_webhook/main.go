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
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/iokuru/sentry/internal/telegram"
)

type Inputs struct {
	AWSProfile       string `json:"aws_profile"`
	AWSRegion        string `json:"aws_region"`
	TelegramBotToken string `json:"telegram_bot_token"`
}

func main() {
	fixSSOCacheTimestamps()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
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

	var inputs Inputs
	if data, err := os.ReadFile("inputs.local.json"); err == nil {
		_ = json.Unmarshal(data, &inputs)
	}

	var optFns []func(*config.LoadOptions) error
	if inputs.AWSRegion != "" {
		optFns = append(optFns, config.WithRegion(inputs.AWSRegion))
	}
	if inputs.AWSProfile != "" && inputs.AWSProfile != "default" {
		optFns = append(optFns, config.WithSharedConfigProfile(inputs.AWSProfile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load aws config: %v\n", err)
		os.Exit(1)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	tgCfg := telegram.LoadConfigFromSSM(ctx, ssmClient)

	token := tgCfg.Token
	if token == "" && inputs.TelegramBotToken != "" {
		token = inputs.TelegramBotToken
	}

	if token == "" {
		fmt.Fprintf(os.Stderr, "telegram bot token not found in SSM or inputs.local.json.\n")
		os.Exit(1)
	}

	secret := tgCfg.WebhookSecret

	reqPayload := map[string]string{
		"url":          webhookUrl,
		"secret_token": secret,
	}
	body, _ := json.Marshal(reqPayload)

	tgEndpoint := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)
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

func fixSSOCacheTimestamps() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	cacheDir := filepath.Join(home, ".aws", "sso", "cache")
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			p := filepath.Join(cacheDir, entry.Name())
			data, err := os.ReadFile(p)
			if err == nil && strings.Contains(string(data), `UTC"`) {
				fixed := strings.ReplaceAll(string(data), `UTC"`, `Z"`)
				_ = os.WriteFile(p, []byte(fixed), 0644)
			}
		}
	}
}
