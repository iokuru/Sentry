package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PresignRequest struct {
	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
}

type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	ImageKey  string `json:"image_key"`
}

var (
	presignClient *s3.PresignClient
	uploadBucket  string
)

func init() {
	uploadBucket = os.Getenv("UPLOADS_BUCKET")
	if uploadBucket == "" {
		uploadBucket = "sentry-uploads"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err == nil {
		s3Client := s3.NewFromConfig(cfg)
		presignClient = s3.NewPresignClient(s3Client)
	}
}

func handlePresign(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var input PresignRequest
	if req.Body != "" {
		_ = json.Unmarshal([]byte(req.Body), &input)
	}

	contentType := input.ContentType
	if contentType == "" {
		contentType = "image/png"
	}

	ext := ".png"
	if input.Filename != "" {
		parsedExt := strings.ToLower(filepath.Ext(input.Filename))
		if parsedExt == ".jpg" || parsedExt == ".jpeg" || parsedExt == ".webp" || parsedExt == ".png" {
			ext = parsedExt
		}
	}

	randBytes := make([]byte, 16)
	_, _ = rand.Read(randBytes)
	key := fmt.Sprintf("screenshots/%d-%s%s", time.Now().Unix(), hex.EncodeToString(randBytes), ext)

	if presignClient == nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"s3 client unavailable"}`,
		}, nil
	}

	presignedReq, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(uploadBucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error":"failed to presign: %s"}`, err.Error()),
		}, nil
	}

	res := PresignResponse{
		UploadURL: presignedReq.URL,
		ImageKey:  key,
	}
	body, _ := json.Marshal(res)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "POST,OPTIONS",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handlePresign)
}
