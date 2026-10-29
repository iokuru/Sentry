package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	comprehendsvc "github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	textractsvc "github.com/aws/aws-sdk-go-v2/service/textract"
	textracttypes "github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/iokuru/sentry/internal/costguard"
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
}

func main() {
	fmt.Println("================================================================")
	fmt.Println("SENTRY PREFLIGHT CHECK")
	fmt.Println("================================================================")

	var failures []string
	var warnings []string

	// 1. check inputs.local.json
	configFile := "inputs.local.json"
	data, err := os.ReadFile(configFile)
	if err != nil {
		failures = append(failures, "missing inputs.local.json: copy inputs.example.json to inputs.local.json and fill in values")
		printSummaryAndExit(failures, warnings)
		return
	}

	var inputs Inputs
	if err := json.Unmarshal(data, &inputs); err != nil {
		failures = append(failures, fmt.Sprintf("invalid inputs.local.json: %v", err))
		printSummaryAndExit(failures, warnings)
		return
	}

	if strings.TrimSpace(inputs.AWSRegion) == "" {
		failures = append(failures, "aws_region is required in inputs.local.json (e.g. us-east-1)")
	}
	if strings.TrimSpace(inputs.ProjectName) == "" {
		failures = append(failures, "project_name is required in inputs.local.json (e.g. sentry)")
	}
	if strings.TrimSpace(inputs.Environment) == "" {
		failures = append(failures, "environment is required in inputs.local.json (e.g. prod)")
	}
	if !strings.Contains(inputs.CaregiverEmail, "@") {
		failures = append(failures, "valid caregiver_email required in inputs.local.json")
	}
	if !strings.Contains(inputs.BudgetAlertEmail, "@") {
		failures = append(failures, "valid budget_alert_email required in inputs.local.json")
	}

	compMode := strings.ToLower(strings.TrimSpace(inputs.UseComprehend))
	if compMode == "" {
		compMode = "auto"
	}
	if compMode != "auto" && compMode != "on" && compMode != "off" {
		failures = append(failures, "use_comprehend must be auto, on, or off")
	}

	ocrMode := strings.ToLower(strings.TrimSpace(inputs.OcrMode))
	if ocrMode == "" {
		ocrMode = "auto"
	}
	if ocrMode != "auto" && ocrMode != "textract" && ocrMode != "browser" {
		failures = append(failures, "ocr_mode must be auto, textract, or browser")
	}

	if len(failures) == 0 {
		fmt.Printf("[PASS] Configuration: inputs.local.json verified (%s / %s)\n", inputs.ProjectName, inputs.Environment)
	}

	// 2. check terraform binary
	tfPath, err := exec.LookPath("terraform")
	if err != nil {
		failures = append(failures, "terraform binary not found on PATH: install terraform to continue")
	} else {
		fmt.Printf("[PASS] Terraform: binary found at %s\n", tfPath)
	}

	// 3. check costguard zero-dollar constraints
	violations, err := costguard.ValidateSourceDirectory("infra")
	if err != nil {
		failures = append(failures, fmt.Sprintf("costguard static audit failed: %v", err))
	} else if len(violations) > 0 {
		for _, v := range violations {
			failures = append(failures, fmt.Sprintf("costguard violation: %s - %s", v.Resource, v.Reason))
		}
	} else {
		fmt.Println("[PASS] Costguard: all infra definitions comply with zero dollar free tier constraints")
	}

	// 4. check aws credentials and region
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
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
		failures = append(failures, fmt.Sprintf("failed to load aws config: %v", err))
	} else {
		stsClient := sts.NewFromConfig(awsCfg)
		caller, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			failures = append(failures, fmt.Sprintf("aws authentication failed: %v (run 'aws configure' or check credentials)", err))
		} else {
			account := "unknown"
			if caller.Account != nil {
				account = *caller.Account
			}
			fmt.Printf("[PASS] AWS Credentials: authenticated as account %s in %s\n", account, awsCfg.Region)

			// 5. probe comprehend
			if compMode == "off" {
				fmt.Println("[INFO] Comprehend: probe skipped (use_comprehend set to off)")
			} else {
				compClient := comprehendsvc.NewFromConfig(awsCfg)
				probeText := "safe test verification probe"
				_, err := compClient.DetectDominantLanguage(ctx, &comprehendsvc.DetectDominantLanguageInput{
					Text: aws.String(probeText),
				})
				if err != nil {
					if compMode == "on" {
						failures = append(failures, fmt.Sprintf("comprehend probe failed (mode on): %v", err))
					} else {
						warnings = append(warnings, fmt.Sprintf("comprehend unavailable in this account (%v); continuing in rules fallback mode", err))
					}
				} else {
					fmt.Println("[PASS] Comprehend: service active and accessible")
				}
			}

			// 6. probe textract
			if ocrMode == "browser" {
				fmt.Println("[INFO] Textract: probe skipped (ocr_mode set to browser)")
			} else {
				textractClient := textractsvc.NewFromConfig(awsCfg)
				onePixelPNG := []byte{
					0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
					0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
					0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
					0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
					0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
					0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
				}
				probeCtx, pCancel := context.WithTimeout(ctx, 4*time.Second)
				defer pCancel()
				_, err := textractClient.DetectDocumentText(probeCtx, &textractsvc.DetectDocumentTextInput{
					Document: &textracttypes.Document{
						Bytes: onePixelPNG,
					},
				})
				if err != nil {
					if ocrMode == "textract" {
						failures = append(failures, fmt.Sprintf("textract probe failed (mode textract): %v", err))
					} else {
						warnings = append(warnings, fmt.Sprintf("textract unavailable in this account (%v); continuing in browser ocr fallback mode", err))
					}
				} else {
					fmt.Println("[PASS] Textract: service active and accessible")
				}
			}
		}
	}

	// 7. check telegram configuration
	if strings.TrimSpace(inputs.TelegramBotToken) != "" {
		tgClient := telegram.NewClient()
		user, err := tgClient.GetMe(ctx, inputs.TelegramBotToken)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("telegram token verification warning: %v", err))
		} else {
			fmt.Printf("[PASS] Telegram: connected to bot @%s\n", user.Username)
		}
	} else {
		warnings = append(warnings, "telegram bot token empty in inputs.local.json; telegram alerts will be skipped")
	}

	printSummaryAndExit(failures, warnings)
}

func printSummaryAndExit(failures, warnings []string) {
	fmt.Println("================================================================")
	if len(warnings) > 0 {
		fmt.Println("NOTICE:")
		for _, w := range warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	if len(failures) > 0 {
		fmt.Println("PREFLIGHT FAILED:")
		for _, f := range failures {
			fmt.Printf("  - %s\n", f)
		}
		fmt.Println("================================================================")
		os.Exit(1)
	}

	fmt.Println("Preflight passed. Ready for deployment.")
	fmt.Println("================================================================")
}
