package costguard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var forbiddenResourceTypes = []string{
	"aws_nat_gateway",
	"aws_db_instance",
	"aws_rds_cluster",
	"aws_rds_cluster_instance",
	"aws_bedrock_custom_model",
	"aws_bedrockagent",
	"aws_opensearch_domain",
	"aws_elasticsearch_domain",
	"aws_elasticache_cluster",
}

type Violation struct {
	Resource string
	Reason   string
}

// PlanRepresentation parses terraform plan json or configuration
type TerraformPlan struct {
	PlannedValues struct {
		RootModule struct {
			Resources    []PlanResource `json:"resources"`
			ChildModules []struct {
				Resources []PlanResource `json:"resources"`
			} `json:"child_modules"`
		} `json:"root_module"`
	} `json:"planned_values"`
}

type PlanResource struct {
	Address string                 `json:"address"`
	Type    string                 `json:"type"`
	Values  map[string]interface{} `json:"values"`
}

func ValidatePlanJSON(planJSON []byte) ([]Violation, error) {
	var plan TerraformPlan
	if err := json.Unmarshal(planJSON, &plan); err != nil {
		return nil, fmt.Errorf("unmarshal plan json: %w", err)
	}

	var allResources []PlanResource
	allResources = append(allResources, plan.PlannedValues.RootModule.Resources...)
	for _, cm := range plan.PlannedValues.RootModule.ChildModules {
		allResources = append(allResources, cm.Resources...)
	}

	var violations []Violation
	for _, res := range allResources {
		// check forbidden resource types
		for _, forbidden := range forbiddenResourceTypes {
			if strings.HasPrefix(res.Type, forbidden) {
				violations = append(violations, Violation{
					Resource: res.Address,
					Reason:   fmt.Sprintf("forbidden paid resource type: %s", res.Type),
				})
			}
		}

		// check dynamodb billing mode
		if res.Type == "aws_dynamodb_table" {
			billingMode, ok := res.Values["billing_mode"].(string)
			if !ok || billingMode != "PAY_PER_REQUEST" {
				violations = append(violations, Violation{
					Resource: res.Address,
					Reason:   "dynamodb table must use PAY_PER_REQUEST billing mode to prevent provisioned capacity charges",
				})
			}
		}

		// check sns sms protocol
		if res.Type == "aws_sns_topic_subscription" {
			protocol, _ := res.Values["protocol"].(string)
			if strings.EqualFold(protocol, "sms") {
				violations = append(violations, Violation{
					Resource: res.Address,
					Reason:   "sns sms subscriptions incur cost; only email is permitted",
				})
			}
		}
	}

	return violations, nil
}

// ValidateSourceDirectory inspects terraform files directly
func ValidateSourceDirectory(dir string) ([]Violation, error) {
	var violations []Violation

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".tf") {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		text := string(content)

		for _, forbidden := range forbiddenResourceTypes {
			re := regexp.MustCompile(`resource\s+"` + forbidden + `[a-z_]*"\s+"[^"]+"`)
			if loc := re.FindString(text); loc != "" {
				violations = append(violations, Violation{
					Resource: fmt.Sprintf("%s (%s)", loc, filepath.Base(path)),
					Reason:   fmt.Sprintf("forbidden paid resource type: %s", forbidden),
				})
			}
		}

		// check for provisioned capacity in dynamodb
		if strings.Contains(text, `aws_dynamodb_table`) {
			if strings.Contains(text, `read_capacity`) || strings.Contains(text, `write_capacity`) {
				violations = append(violations, Violation{
					Resource: filepath.Base(path),
					Reason:   "dynamodb provisioned capacity detected; must use PAY_PER_REQUEST",
				})
			}
		}

		return nil
	})

	return violations, err
}
