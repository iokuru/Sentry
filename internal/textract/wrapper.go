package textract

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type Client interface {
	DetectDocumentText(ctx context.Context, params *textract.DetectDocumentTextInput, optFns ...func(*textract.Options)) (*textract.DetectDocumentTextOutput, error)
}

func ExtractText(ctx context.Context, client Client, bucket, key string) (string, error) {
	input := &textract.DetectDocumentTextInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: &bucket,
				Name:   &key,
			},
		},
	}

	out, err := client.DetectDocumentText(ctx, input)
	if err != nil {
		return "", fmt.Errorf("textract detect text from %s/%s: %w", bucket, key, err)
	}

	var lines []string
	for _, block := range out.Blocks {
		if block.BlockType == types.BlockTypeLine && block.Text != nil {
			lines = append(lines, *block.Text)
		}
	}

	return strings.Join(lines, "\n"), nil
}

func IsAccessOrSubscriptionError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "subscriptionrequired") ||
		strings.Contains(s, "accessdenied") ||
		strings.Contains(s, "not authorized") ||
		strings.Contains(s, "optinrequired") ||
		strings.Contains(s, "unrecognizedclientexception") ||
		strings.Contains(s, "subscription")
}
