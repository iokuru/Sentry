package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	sqssvc "github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/iokuru/sentry/internal/quota"
)

type Inputs struct {
	AWSProfile  string `json:"aws_profile"`
	AWSRegion   string `json:"aws_region"`
	ProjectName   string `json:"project_name"`
	Environment   string `json:"environment"`
	UseComprehend string `json:"use_comprehend"`
}

type StatusItem struct {
	Component string
	Status    string
	Details   string
}

func main() {
	localFlag := flag.Bool("local", false, "check local devserver health instead of deployed AWS resources")
	flag.Parse()

	fmt.Println("================================================================================")
	fmt.Println("SENTRY STATUS REPORT")
	fmt.Println("================================================================================")

	var items []StatusItem
	httpClient := &http.Client{Timeout: 5 * time.Second}

	if *localFlag {
		// local devserver check
		res, err := httpClient.Get("http://localhost:8080/api/health")
		if err != nil {
			items = append(items, StatusItem{
				Component: "DevServer",
				Status:    "OFFLINE",
				Details:   fmt.Sprintf("http://localhost:8080 unreachable: %v", err),
			})
		} else {
			defer res.Body.Close()
			items = append(items, StatusItem{
				Component: "DevServer",
				Status:    fmt.Sprintf("HEALTHY (%d)", res.StatusCode),
				Details:   "http://localhost:8080 responding",
			})
		}

		printTable(items)
		return
	}

	// Deployed AWS mode
	cloudfrontURL, _ := getTerraformOutput("cloudfront_url")
	apiURL, _ := getTerraformOutput("api_url")
	dlqURL, _ := getTerraformOutput("analysis_dlq_url")
	region, _ := getTerraformOutput("region")

	// 1. Check CloudFront
	if cloudfrontURL != "" {
		res, err := httpClient.Get(cloudfrontURL)
		if err != nil {
			items = append(items, StatusItem{
				Component: "CloudFront",
				Status:    "UNREACHABLE",
				Details:   err.Error(),
			})
		} else {
			defer res.Body.Close()
			items = append(items, StatusItem{
				Component: "CloudFront",
				Status:    fmt.Sprintf("HEALTHY (%d)", res.StatusCode),
				Details:   cloudfrontURL,
			})
		}
	} else {
		items = append(items, StatusItem{
			Component: "CloudFront",
			Status:    "UNKNOWN",
			Details:   "terraform output cloudfront_url empty",
		})
	}

	// 2. Check API Gateway Health Route
	if apiURL != "" {
		healthURL := strings.TrimSuffix(apiURL, "/") + "/api/health"
		res, err := httpClient.Get(healthURL)
		if err != nil {
			items = append(items, StatusItem{
				Component: "API Gateway",
				Status:    "UNREACHABLE",
				Details:   err.Error(),
			})
		} else {
			defer res.Body.Close()
			items = append(items, StatusItem{
				Component: "API Gateway",
				Status:    fmt.Sprintf("HEALTHY (%d)", res.StatusCode),
				Details:   healthURL,
			})
		}
	} else {
		items = append(items, StatusItem{
			Component: "API Gateway",
			Status:    "UNKNOWN",
			Details:   "terraform output api_url empty",
		})
	}

	// 3. Check SQS DLQ Depth and CloudWatch Alarm
	var inputs Inputs
	if data, err := os.ReadFile("inputs.local.json"); err == nil {
		_ = json.Unmarshal(data, &inputs)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var optFns []func(*config.LoadOptions) error
	if region != "" {
		optFns = append(optFns, config.WithRegion(region))
	} else if inputs.AWSRegion != "" {
		optFns = append(optFns, config.WithRegion(inputs.AWSRegion))
	}
	if inputs.AWSProfile != "" && inputs.AWSProfile != "default" {
		optFns = append(optFns, config.WithSharedConfigProfile(inputs.AWSProfile))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err == nil {
		if dlqURL != "" {
			sqsClient := sqssvc.NewFromConfig(awsCfg)
			attrOut, err := sqsClient.GetQueueAttributes(ctx, &sqssvc.GetQueueAttributesInput{
				QueueUrl: aws.String(dlqURL),
				AttributeNames: []sqstypes.QueueAttributeName{
					sqstypes.QueueAttributeNameApproximateNumberOfMessages,
				},
			})
			if err != nil {
				items = append(items, StatusItem{
					Component: "Analysis DLQ",
					Status:    "ERROR",
					Details:   fmt.Sprintf("failed to get queue attributes: %v", err),
				})
			} else {
				msgCountStr := attrOut.Attributes[string(sqstypes.QueueAttributeNameApproximateNumberOfMessages)]
				count, _ := strconv.Atoi(msgCountStr)
				statusStr := "OK (0 msgs)"
				if count > 0 {
					statusStr = fmt.Sprintf("ATTENTION (%d msgs)", count)
				}
				items = append(items, StatusItem{
					Component: "Analysis DLQ",
					Status:    statusStr,
					Details:   fmt.Sprintf("dead letter queue has %d pending messages", count),
				})
			}
		}

		prefix := "sentry-prod"
		if inputs.ProjectName != "" && inputs.Environment != "" {
			prefix = fmt.Sprintf("%s-%s", inputs.ProjectName, inputs.Environment)
		}
		alarmName := fmt.Sprintf("%s-analysis-dlq-depth", prefix)

		alarmCmd := exec.Command("aws", "cloudwatch", "describe-alarms", "--alarm-names", alarmName, "--query", "MetricAlarms[0].StateValue", "--output", "text")
		if alarmOut, err := alarmCmd.Output(); err == nil {
			stateVal := strings.TrimSpace(string(alarmOut))
			if stateVal != "" && stateVal != "None" {
				items = append(items, StatusItem{
					Component: "DLQ Depth Alarm",
					Status:    stateVal,
					Details:   fmt.Sprintf("CloudWatch alarm %s", alarmName),
				})
			}
		}
	}

	// 4. Check Textract and Comprehend Free Tier Quotas
	limiter := quota.NewLimiter()
	textractUsage := limiter.GetUsage("textract")
	comprehendUsage := limiter.GetUsage("comprehend")

	items = append(items, StatusItem{
		Component: "Textract Cap",
		Status:    "OK",
		Details:   fmt.Sprintf("%d used / 100 daily limit (%d remaining)", textractUsage, 100-textractUsage),
	})
	items = append(items, StatusItem{
		Component: "Comprehend Cap",
		Status:    "OK",
		Details:   fmt.Sprintf("%d used / 100 daily limit (%d remaining)", comprehendUsage, 100-comprehendUsage),
	})

	printTable(items)
}

func printTable(items []StatusItem) {
	fmt.Printf("%-18s %-16s %s\n", "COMPONENT", "STATUS", "DETAILS")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, it := range items {
		fmt.Printf("%-18s %-16s %s\n", it.Component, it.Status, it.Details)
	}
	fmt.Println("================================================================================")
}

func getTerraformOutput(name string) (string, error) {
	cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", name)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
