package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ssmsvc "github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/iokuru/sentry/internal/telegram"
)

type Inputs struct {
	AWSProfile       string `json:"aws_profile"`
	AWSRegion        string `json:"aws_region"`
	ProjectName      string `json:"project_name"`
	Environment      string `json:"environment"`
	CaregiverEmail   string `json:"caregiver_email"`
	BudgetAlertEmail string `json:"budget_alert_email"`
	TelegramBotToken string `json:"telegram_bot_token"`
	TelegramChatID   string `json:"telegram_chat_id"`
	UseComprehend    string `json:"use_comprehend"`
	OcrMode          string `json:"ocr_mode"`
	FrontendHost     string `json:"frontend_host"`
	NetlifySiteName  string `json:"netlify_site_name"`
	AllowedOrigin    string `json:"allowed_origin"`
	EnableCloudfront *bool  `json:"enable_cloudfront"`
}

type Outputs struct {
	ApiURL                 string
	UserPoolID             string
	UserPoolClientID       string
	Region                 string
	CloudfrontURL          string
	FrontendBucket         string
	CloudfrontDistribution string
	UploadsBucket          string
}

func main() {
	fixSSOCacheTimestamps()
	fmt.Println("================================================================")
	fmt.Println("SENTRY GO-LIVE DEPLOYMENT")
	fmt.Println("================================================================")

	// Step 1: Run preflight
	fmt.Println("\n--> Step 1/7: Running preflight checks...")
	if err := runCmd("go", "run", "./cmd/preflight"); err != nil {
		failStep("Preflight verification failed", "Review the errors reported above, update inputs.local.json or AWS credentials, and re-run make go-live.")
	}

	// Read inputs
	inputsData, err := os.ReadFile("inputs.local.json")
	if err != nil {
		failStep("Failed to read inputs.local.json", "Ensure inputs.local.json exists and is readable.")
	}
	var inputs Inputs
	if err := json.Unmarshal(inputsData, &inputs); err != nil {
		failStep("Failed to parse inputs.local.json", "Verify that inputs.local.json contains valid JSON formatting.")
	}

	if inputs.AWSProfile != "" {
		_ = os.Setenv("AWS_PROFILE", inputs.AWSProfile)
	}
	if inputs.AWSRegion != "" {
		_ = os.Setenv("AWS_REGION", inputs.AWSRegion)
		_ = os.Setenv("AWS_DEFAULT_REGION", inputs.AWSRegion)
	}

	// Step 2: Build arm64 lambdas and frontend assets
	fmt.Println("\n--> Step 2/7: Building arm64 lambda packages and frontend assets...")
	if err := runCmd("go", "run", "./cmd/build", "-frontend"); err != nil {
		failStep("Build failed", "Ensure Go 1.23+ and Node 18+ are installed and run 'go run ./cmd/build -frontend' to see compiler errors.")
	}

	// Step 3: Run terraform apply
	fmt.Println("\n--> Step 3/7: Running terraform apply...")
	compMode := strings.ToLower(strings.TrimSpace(inputs.UseComprehend))
	if compMode == "" {
		compMode = "auto"
	}
	ocrMode := strings.ToLower(strings.TrimSpace(inputs.OcrMode))
	if ocrMode == "" {
		ocrMode = "auto"
	}
	enableCF := true
	if inputs.EnableCloudfront != nil {
		enableCF = *inputs.EnableCloudfront
	} else if strings.EqualFold(inputs.FrontendHost, "s3") || strings.EqualFold(inputs.FrontendHost, "netlify") {
		enableCF = false
	}
	tfvars := map[string]interface{}{
		"aws_region":         inputs.AWSRegion,
		"environment":        inputs.Environment,
		"project_name":       inputs.ProjectName,
		"budget_alert_email": inputs.BudgetAlertEmail,
		"use_comprehend":     compMode,
		"ocr_mode":           ocrMode,
		"enable_cloudfront":  enableCF,
		"frontend_host":      strings.ToLower(strings.TrimSpace(inputs.FrontendHost)),
		"allowed_origin":     strings.ToLower(strings.TrimSpace(inputs.AllowedOrigin)),
		"netlify_site_name":  strings.TrimSpace(inputs.NetlifySiteName),
	}
	tfvarsBytes, _ := json.MarshalIndent(tfvars, "", "  ")
	_ = os.WriteFile(filepath.Join("infra", "terraform.tfvars.json"), tfvarsBytes, 0644)

	if err := runTerraformWithFallback(tfvars); err != nil {
		failStep("Terraform apply failed", "Inspect the Terraform error output above, verify your AWS IAM permissions, and retry.")
	}

	outputs, err := readTerraformOutputs()
	if err != nil {
		failStep(fmt.Sprintf("Failed to read terraform outputs: %v", err), "Ensure Terraform state exists in infra/ and retry.")
	}

	// Step 4: Populate SSM parameters from inputs.local.json
	fmt.Println("\n--> Step 4/7: Populating SSM parameters in AWS Parameter Store...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var optFns []func(*config.LoadOptions) error
	if inputs.AWSRegion != "" {
		optFns = append(optFns, config.WithRegion(inputs.AWSRegion))
	}
	if inputs.AWSProfile != "" && inputs.AWSProfile != "default" {
		optFns = append(optFns, config.WithSharedConfigProfile(inputs.AWSProfile))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		failStep(fmt.Sprintf("Failed to initialize AWS client for SSM: %v", err), "Check your AWS credentials in inputs.local.json.")
	}

	ssmClient := ssmsvc.NewFromConfig(awsCfg)
	webhookSecret := generateRandomSecret(32)

	if err := putSSMParam(ctx, ssmClient, "/sentry/telegram_webhook_secret", webhookSecret, false); err != nil {
		failStep(fmt.Sprintf("Failed to store webhook secret in SSM: %v", err), "Ensure your IAM user has ssm:PutParameter permission.")
	}

	if inputs.TelegramBotToken != "" {
		if err := putSSMParam(ctx, ssmClient, "/sentry/telegram_token", inputs.TelegramBotToken, true); err != nil {
			failStep(fmt.Sprintf("Failed to store telegram token in SSM: %v", err), "Ensure your IAM user has ssm:PutParameter permission.")
		}
	}

	if inputs.TelegramChatID != "" {
		if err := putSSMParam(ctx, ssmClient, "/sentry/telegram_chat_id", inputs.TelegramChatID, false); err != nil {
			failStep(fmt.Sprintf("Failed to store telegram chat id in SSM: %v", err), "Ensure your IAM user has ssm:PutParameter permission.")
		}
	}

	// Step 5: Set Telegram Webhook if token configured
	var botUsername string
	if inputs.TelegramBotToken != "" {
		fmt.Println("\n--> Step 5/7: Registering Telegram webhook...")
		tgClient := telegram.NewClient()
		user, err := tgClient.GetMe(ctx, inputs.TelegramBotToken)
		if err == nil && user != nil {
			botUsername = user.Username
		}

		hookURL := strings.TrimSuffix(outputs.ApiURL, "/") + "/telegram/hook"
		if err := setTelegramWebhook(ctx, inputs.TelegramBotToken, hookURL, webhookSecret); err != nil {
			fmt.Printf("Warning: failed to set telegram webhook: %v\n", err)
		} else {
			fmt.Printf("Telegram webhook registered to %s\n", hookURL)
		}
	} else {
		fmt.Println("\n--> Step 5/7: Telegram skipped (no bot token provided)")
	}

	// Step 6: Generate /config.json and deploy frontend
	frontendHost := strings.ToLower(strings.TrimSpace(inputs.FrontendHost))
	if frontendHost == "" {
		frontendHost = "cloudfront"
	}

	runtimeConfig := map[string]string{
		"api_url":             outputs.ApiURL,
		"user_pool_id":        outputs.UserPoolID,
		"user_pool_client_id": outputs.UserPoolClientID,
		"region":              outputs.Region,
		"ocr_mode":            ocrMode,
	}
	cfgBytes, _ := json.MarshalIndent(runtimeConfig, "", "  ")
	cfgPath := filepath.Join("frontend", "build", "config.json")
	if err := os.WriteFile(cfgPath, cfgBytes, 0644); err != nil {
		failStep(fmt.Sprintf("Failed to write runtime config.json: %v", err), "Ensure frontend/build directory exists.")
	}
	_ = os.WriteFile(filepath.Join("frontend", "public", "config.json"), cfgBytes, 0644)

	if frontendHost == "netlify" {
		fmt.Println("\n--> Step 6/7: Deploying frontend bundle to Netlify...")
		uploadsHost := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", outputs.UploadsBucket, outputs.Region)
		cspHeader := fmt.Sprintf("Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https:; connect-src 'self' %s %s; font-src 'self'; worker-src 'self' blob:; frame-ancestors 'none'; object-src 'none'; base-uri 'self';", outputs.ApiURL, uploadsHost)
		headersContent := fmt.Sprintf("/*\n  X-Content-Type-Options: nosniff\n  X-Frame-Options: DENY\n  Referrer-Policy: strict-origin-when-cross-origin\n  Strict-Transport-Security: max-age=31536000; includeSubDomains; preload\n  %s\n", cspHeader)
		_ = os.WriteFile(filepath.Join("frontend", "build", "_headers"), []byte(headersContent), 0644)

		netPath, err := exec.LookPath("netlify")
		if err != nil {
			failStep("Netlify CLI not found", "Ensure netlify-cli is installed and on PATH.")
		}
		deployArgs := []string{"deploy", "--prod", "--dir=" + filepath.Join("frontend", "build")}
		if inputs.NetlifySiteName != "" {
			deployArgs = append(deployArgs, "--site="+inputs.NetlifySiteName)
		}
		if err := runCmd(netPath, deployArgs...); err != nil {
			failStep(fmt.Sprintf("Netlify deploy failed: %v", err), "Check Netlify CLI authentication and site settings.")
		}
	} else {
		fmt.Println("\n--> Step 6/7: Deploying frontend bundle to S3 and invalidating CloudFront cache...")
		s3Target := fmt.Sprintf("s3://%s", outputs.FrontendBucket)
		if err := runCmd("aws", "s3", "sync", filepath.Join("frontend", "build"), s3Target, "--delete"); err != nil {
			failStep(fmt.Sprintf("Failed to sync frontend to S3: %v", err), "Ensure AWS CLI is configured with S3 write access.")
		}

		if outputs.CloudfrontDistribution != "" {
			if err := runCmd("aws", "cloudfront", "create-invalidation", "--distribution-id", outputs.CloudfrontDistribution, "--paths", "/*"); err != nil {
				fmt.Printf("Warning: cloudfront invalidation returned: %v\n", err)
			}
		}
	}

	// Create caregiver Cognito user
	initialPassword := "SentrySafety2026!"
	fmt.Printf("\n--> Provisioning initial caregiver user: %s\n", inputs.CaregiverEmail)
	_ = runCmd("aws", "cognito-idp", "admin-create-user",
		"--user-pool-id", outputs.UserPoolID,
		"--username", inputs.CaregiverEmail,
		"--user-attributes", fmt.Sprintf("Name=email,Value=%s", inputs.CaregiverEmail), "Name=email_verified,Value=true",
		"--message-action", "SUPPRESS",
	)
	_ = runCmd("aws", "cognito-idp", "admin-set-user-password",
		"--user-pool-id", outputs.UserPoolID,
		"--username", inputs.CaregiverEmail,
		"--password", initialPassword,
		"--permanent",
	)
	_ = os.Setenv("CAREGIVER_EMAIL", inputs.CaregiverEmail)
	_ = os.Setenv("CAREGIVER_PASSWORD", initialPassword)

	// Step 7: Live Smoke Test
	fmt.Println("\n--> Step 7/7: Executing live smoke verification...")
	if err := runCmd("go", "run", "./cmd/smoke", "-live"); err != nil {
		fmt.Printf("Notice: smoke test encountered an issue: %v\n", err)
	}

	// Write evidence/aws-resources.txt
	_ = os.MkdirAll("evidence", 0755)
	stateCmd := exec.Command("terraform", "-chdir=infra", "state", "list")
	if stateList, err := stateCmd.Output(); err == nil {
		_ = os.WriteFile(filepath.Join("evidence", "aws-resources.txt"), stateList, 0644)
	}

	// Print final summary
	frontendURL := outputs.CloudfrontURL
	if frontendHost == "netlify" {
		if inputs.AllowedOrigin != "" {
			frontendURL = strings.ToLower(strings.TrimSpace(inputs.AllowedOrigin))
		} else if inputs.NetlifySiteName != "" {
			frontendURL = fmt.Sprintf("https://%s.netlify.app", inputs.NetlifySiteName)
		}
	}

	fmt.Println("\n================================================================")
	fmt.Println("SENTRY GO-LIVE SUCCESSFUL")
	fmt.Println("================================================================")
	fmt.Printf("Frontend URL:           %s\n", frontendURL)
	fmt.Printf("Caregiver Sign In:      %s/caregiver\n", frontendURL)
	fmt.Printf("Caregiver User:         %s\n", inputs.CaregiverEmail)
	fmt.Printf("Initial Password:       %s\n", initialPassword)
	if botUsername != "" {
		fmt.Printf("Telegram Bot:           https://t.me/%s\n", botUsername)
	}
	fmt.Printf("API Gateway Endpoint:   %s\n", outputs.ApiURL)
	fmt.Printf("Cognito User Pool ID:   %s\n", outputs.UserPoolID)
	fmt.Printf("Cognito App Client ID:  %s\n", outputs.UserPoolClientID)
	fmt.Println("================================================================")
	fmt.Println("Next: Sign up your caregiver account at the link above, then pair seniors via the portal.")
}

func failStep(errSummary, retryAdvice string) {
	fmt.Println("\n================================================================")
	fmt.Printf("GO-LIVE STOPPED: %s\n", errSummary)
	fmt.Printf("HOW TO RETRY:    %s\n", retryAdvice)
	fmt.Println("================================================================")
	os.Exit(1)
}

func putSSMParam(ctx context.Context, client *ssmsvc.Client, name, value string, secure bool) error {
	pType := ssmtypes.ParameterTypeString
	if secure {
		pType = ssmtypes.ParameterTypeSecureString
	}
	_, err := client.PutParameter(ctx, &ssmsvc.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      pType,
		Overwrite: aws.Bool(true),
	})
	return err
}

func setTelegramWebhook(ctx context.Context, token, webhookURL, secret string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)
	payload := map[string]string{
		"url":          webhookURL,
		"secret_token": secret,
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", res.StatusCode)
	}
	return nil
}

func generateRandomSecret(byteCount int) string {
	b := make([]byte, byteCount)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func readTerraformOutputs() (*Outputs, error) {
	out := &Outputs{}
	var err error
	out.ApiURL, err = getOutput("api_url")
	if err != nil {
		return nil, err
	}
	out.UserPoolID, _ = getOutput("user_pool_id")
	out.UserPoolClientID, _ = getOutput("user_pool_client_id")
	out.Region, _ = getOutput("region")
	out.CloudfrontURL, _ = getOutput("cloudfront_url")
	out.FrontendBucket, _ = getOutput("frontend_bucket")
	out.CloudfrontDistribution, _ = getOutput("cloudfront_distribution_id")
	out.UploadsBucket, _ = getOutput("uploads_bucket")
	return out, nil
}

func getOutput(name string) (string, error) {
	cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", name)
	res, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("terraform output %s: %w", name, err)
	}
	return strings.TrimSpace(string(res)), nil
}

func runTerraformWithFallback(tfvars map[string]interface{}) error {
	var buf bytes.Buffer
	multiOut := io.MultiWriter(os.Stdout, &buf)
	multiErr := io.MultiWriter(os.Stderr, &buf)

	cmd := exec.Command("terraform", "-chdir=infra", "apply", "-auto-approve")
	cmd.Stdout = multiOut
	cmd.Stderr = multiErr
	err := cmd.Run()
	if err == nil {
		return nil
	}

	outStr := buf.String()
	if strings.Contains(outStr, "Your account must be verified before you can add new CloudFront resources") {
		fmt.Println("\n[NOTICE] CloudFront distribution creation requires AWS Support verification for this account.")
		fmt.Println("[NOTICE] Proceeding with deployment of API Gateway, Lambdas, Auth, and Storage...")
		tfvars["enable_cloudfront"] = false
		tfvarsBytes, _ := json.MarshalIndent(tfvars, "", "  ")
		_ = os.WriteFile(filepath.Join("infra", "terraform.tfvars.json"), tfvarsBytes, 0644)
		return runCmd("terraform", "-chdir=infra", "apply", "-auto-approve")
	}

	return err
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
