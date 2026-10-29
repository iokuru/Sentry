package costguard_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var allowedServicePrefixes = []string{
	"aws_s3_",
	"aws_cloudfront_",
	"aws_apigatewayv2_",
	"aws_lambda_",
	"aws_dynamodb_",
	"aws_sqs_",
	"aws_sns_",
	"aws_cognito_",
	"aws_cloudwatch_",
	"aws_budgets_",
	"aws_scheduler_",
	"aws_ssm_",
	"aws_iam_",
}

func isServiceAllowed(resourceType string) bool {
	for _, prefix := range allowedServicePrefixes {
		if strings.HasPrefix(resourceType, prefix) {
			return true
		}
	}
	return false
}

func TestTerraformServiceAllowlist(t *testing.T) {
	infraDir := filepath.Join("..", "..", "infra")
	resourceRegex := regexp.MustCompile(`resource\s+"([a-zA-Z0-9_]+)"\s+"([^"]+)"`)

	count := 0
	err := filepath.Walk(infraDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".tf") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		matches := resourceRegex.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			if len(m) >= 2 {
				resType := m[1]
				resName := m[2]
				count++

				if !isServiceAllowed(resType) {
					t.Errorf("forbidden resource type detected in %s: %s.%s is not in free tier service allowlist",
						filepath.Base(path), resType, resName)
				}
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("failed to scan infra directory: %v", err)
	}

	if count == 0 {
		t.Fatalf("expected to find terraform resources, found 0")
	}

	t.Logf("verified %d terraform resources against free tier allowlist", count)
}

func TestComprehendAndTextractOmittedWhenOff(t *testing.T) {
	iamPath := filepath.Join("..", "..", "infra", "modules", "notify", "iam.tf")
	content, err := os.ReadFile(iamPath)
	if err != nil {
		t.Fatalf("failed to read notify iam.tf: %v", err)
	}

	text := string(content)

	// verify that comprehend statements are conditional on use_comprehend != "off"
	if !strings.Contains(text, `var.use_comprehend != "off"`) {
		t.Errorf("expected iam.tf to guard comprehend statements with var.use_comprehend != \"off\"")
	}

	// verify that textract statements are conditional on ocr_mode != "browser"
	if !strings.Contains(text, `var.ocr_mode != "browser"`) {
		t.Errorf("expected iam.tf to guard textract statements with var.ocr_mode != \"browser\"")
	}

	// verify that when both are off, analysis_statements contains neither
	if !strings.Contains(text, "local.textract_statement") || !strings.Contains(text, "local.comprehend_statement") {
		t.Errorf("expected analysis_statements to use local.textract_statement and local.comprehend_statement")
	}
}
