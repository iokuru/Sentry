package costguard

import (
	"path/filepath"
	"testing"
)

func TestValidateSourceDirectoryClean(t *testing.T) {
	infraDir := filepath.Join("..", "..", "infra")
	violations, err := ValidateSourceDirectory(infraDir)
	if err != nil {
		t.Fatalf("scan infra dir: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf("expected 0 cost violations in sentry infra, got %d: %+v", len(violations), violations)
	}
}

func TestValidatePlanJSONForbiddenTypes(t *testing.T) {
	planJSON := []byte(`{
		"planned_values": {
			"root_module": {
				"resources": [
					{
						"address": "aws_nat_gateway.gw",
						"type": "aws_nat_gateway",
						"values": {}
					},
					{
						"address": "aws_db_instance.db",
						"type": "aws_db_instance",
						"values": {}
					},
					{
						"address": "aws_dynamodb_table.bad",
						"type": "aws_dynamodb_table",
						"values": {
							"billing_mode": "PROVISIONED"
						}
					}
				],
				"child_modules": []
			}
		}
	}`)

	violations, err := ValidatePlanJSON(planJSON)
	if err != nil {
		t.Fatalf("validate plan: %v", err)
	}

	if len(violations) != 3 {
		t.Fatalf("expected 3 violations, got %d: %+v", len(violations), violations)
	}
}

func TestValidatePlanJSONClean(t *testing.T) {
	planJSON := []byte(`{
		"planned_values": {
			"root_module": {
				"resources": [
					{
						"address": "aws_s3_bucket.uploads",
						"type": "aws_s3_bucket",
						"values": {}
					},
					{
						"address": "aws_dynamodb_table.incidents",
						"type": "aws_dynamodb_table",
						"values": {
							"billing_mode": "PAY_PER_REQUEST"
						}
					}
				],
				"child_modules": []
			}
		}
	}`)

	violations, err := ValidatePlanJSON(planJSON)
	if err != nil {
		t.Fatalf("validate plan: %v", err)
	}

	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}
