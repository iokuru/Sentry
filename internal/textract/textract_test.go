package textract_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	internaltextract "github.com/iokuru/sentry/internal/textract"
)

type mockTextractClient struct {
	output *textract.DetectDocumentTextOutput
	err    error
}

func (m *mockTextractClient) DetectDocumentText(ctx context.Context, params *textract.DetectDocumentTextInput, optFns ...func(*textract.Options)) (*textract.DetectDocumentTextOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.output, nil
}

func TestExtractText(t *testing.T) {
	tests := []struct {
		name        string
		client      *mockTextractClient
		expected    string
		expectError bool
	}{
		{
			name: "empty blocks",
			client: &mockTextractClient{
				output: &textract.DetectDocumentTextOutput{
					Blocks: []types.Block{},
				},
			},
			expected:    "",
			expectError: false,
		},
		{
			name: "single line block",
			client: &mockTextractClient{
				output: &textract.DetectDocumentTextOutput{
					Blocks: []types.Block{
						{
							BlockType:  types.BlockTypeLine,
							Text:       aws.String("URGENT: Your account has been suspended"),
							Confidence: aws.Float32(99.4),
						},
					},
				},
			},
			expected:    "URGENT: Your account has been suspended",
			expectError: false,
		},
		{
			name: "multiple lines with mixed block types",
			client: &mockTextractClient{
				output: &textract.DetectDocumentTextOutput{
					Blocks: []types.Block{
						{
							BlockType: types.BlockTypePage,
						},
						{
							BlockType:  types.BlockTypeLine,
							Text:       aws.String("Bank Alert: suspicious login detected"),
							Confidence: aws.Float32(98.2),
						},
						{
							BlockType:  types.BlockTypeWord,
							Text:       aws.String("Bank"),
							Confidence: aws.Float32(99.0),
						},
						{
							BlockType:  types.BlockTypeLine,
							Text:       aws.String("Click here to verify your identity immediately"),
							Confidence: aws.Float32(95.6),
						},
					},
				},
			},
			expected:    "Bank Alert: suspicious login detected\nClick here to verify your identity immediately",
			expectError: false,
		},
		{
			name: "nil text in line block is ignored safely",
			client: &mockTextractClient{
				output: &textract.DetectDocumentTextOutput{
					Blocks: []types.Block{
						{
							BlockType: types.BlockTypeLine,
							Text:      nil,
						},
						{
							BlockType: types.BlockTypeLine,
							Text:      aws.String("Valid line text"),
						},
					},
				},
			},
			expected:    "Valid line text",
			expectError: false,
		},
		{
			name: "client error is propagated",
			client: &mockTextractClient{
				err: errors.New("access denied"),
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := internaltextract.ExtractText(context.Background(), tt.client, "test-bucket", "test-key.png")
			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, res)
			}
		})
	}
}
