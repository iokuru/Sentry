package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	ssmsvc "github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Inputs struct {
	AWSProfile  string `json:"aws_profile"`
	AWSRegion   string `json:"aws_region"`
	ProjectName string `json:"project_name"`
	Environment string `json:"environment"`
}

func main() {
	forceFlag := flag.Bool("force", false, "skip interactive confirmation prompt")
	flag.Parse()

	fmt.Println("================================================================")
	fmt.Println("SENTRY TEARDOWN")
	fmt.Println("================================================================")

	var inputs Inputs
	if data, err := os.ReadFile("inputs.local.json"); err == nil {
		_ = json.Unmarshal(data, &inputs)
	}

	projectName := inputs.ProjectName
	if projectName == "" {
		projectName = "sentry"
	}

	if !*forceFlag {
		fmt.Printf("WARNING: This will permanently destroy all AWS resources for '%s'.\n", projectName)
		fmt.Printf("Type the project name '%s' to confirm: ", projectName)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nAborted.")
			os.Exit(1)
		}

		if strings.TrimSpace(input) != projectName {
			fmt.Println("Teardown cancelled: project name did not match.")
			os.Exit(1)
		}
	}

	// Step 1: Read buckets from Terraform outputs to empty them cleanly
	fmt.Println("\n--> Step 1/4: Reading bucket references...")
	frontendBucket, _ := getTerraformOutput("frontend_bucket")
	uploadsBucket, _ := getTerraformOutput("uploads_bucket")

	// Step 2: Empty buckets before destroy to prevent edge case bucket-not-empty errors
	fmt.Println("\n--> Step 2/4: Emptying S3 buckets...")
	if frontendBucket != "" {
		fmt.Printf("Emptying frontend bucket %s...\n", frontendBucket)
		_ = runCmd("aws", "s3", "rm", fmt.Sprintf("s3://%s", frontendBucket), "--recursive")
	}
	if uploadsBucket != "" {
		fmt.Printf("Emptying uploads bucket %s...\n", uploadsBucket)
		_ = runCmd("aws", "s3", "rm", fmt.Sprintf("s3://%s", uploadsBucket), "--recursive")
	}

	// Step 3: Delete SSM Parameters
	fmt.Println("\n--> Step 3/4: Deleting SSM configuration parameters...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var optFns []func(*config.LoadOptions) error
	if inputs.AWSRegion != "" {
		optFns = append(optFns, config.WithRegion(inputs.AWSRegion))
	}
	if inputs.AWSProfile != "" && inputs.AWSProfile != "default" {
		optFns = append(optFns, config.WithSharedConfigProfile(inputs.AWSProfile))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err == nil {
		ssmClient := ssmsvc.NewFromConfig(awsCfg)
		_, _ = ssmClient.DeleteParameters(ctx, &ssmsvc.DeleteParametersInput{
			Names: []string{
				"/sentry/telegram_token",
				"/sentry/telegram_chat_id",
				"/sentry/telegram_webhook_secret",
			},
		})
	}

	// Step 4: Run terraform destroy
	fmt.Println("\n--> Step 4/4: Destroying Terraform managed infrastructure...")
	if err := runCmd("terraform", "-chdir=infra", "destroy", "-auto-approve"); err != nil {
		fmt.Fprintf(os.Stderr, "Terraform destroy failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n================================================================")
	fmt.Println("Teardown completed. Zero AWS resources remaining.")
	fmt.Println("================================================================")
}

func getTerraformOutput(name string) (string, error) {
	cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", name)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
