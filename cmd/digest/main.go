package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/iokuru/sentry/internal/store"
)

var (
	incidentStore  *store.Store
	snsClient      *sns.Client
	alertsTopicArn string
)

func init() {
	alertsTopicArn = os.Getenv("CAREGIVER_ALERTS_TOPIC_ARN")
	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	linksTable := os.Getenv("LINKS_TABLE")
	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err == nil {
		dbClient := dynamodb.NewFromConfig(cfg)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)
		snsClient = sns.NewFromConfig(cfg)
	}
}

func handleDigest(ctx context.Context) error {
	if incidentStore == nil {
		return fmt.Errorf("incident store unavailable")
	}

	unresolved, err := incidentStore.ListUnresolvedIncidents(ctx)
	if err != nil {
		return fmt.Errorf("list unresolved incidents: %w", err)
	}

	if len(unresolved) == 0 {
		return nil
	}

	if snsClient == nil || alertsTopicArn == "" {
		return nil
	}

	var items []string
	for _, inc := range unresolved {
		items = append(items, fmt.Sprintf(
			"- %s | Severity: %s (%d/100) | Category: %s | Status: %s | Time: %s",
			inc.ID, inc.RiskLevel, inc.Score, inc.Category, inc.Status, inc.CreatedAt,
		))
	}

	subject := fmt.Sprintf("[Sentry] Daily digest: %d unresolved incident(s)", len(unresolved))
	body := fmt.Sprintf(
		"Sentry Daily Summary - %s\n\n"+
			"The following %d incident(s) require caregiver review:\n\n"+
			"%s\n\n"+
			"Open the caregiver portal to acknowledge or resolve these items.",
		time.Now().UTC().Format("2006-01-02"), len(unresolved),
		strings.Join(items, "\n"),
	)

	_, err = snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(alertsTopicArn),
		Subject:  aws.String(subject),
		Message:  aws.String(body),
	})
	if err != nil {
		return fmt.Errorf("publish daily digest: %w", err)
	}

	return nil
}

func main() {
	lambda.Start(handleDigest)
}
