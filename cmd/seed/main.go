package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/risk"
	"github.com/iokuru/sentry/internal/store"
)

type SeedScenario struct {
	Text     string
	Country  string
	Status   model.IncidentStatus
	TimeAgo  time.Duration
}

var scenarios = []SeedScenario{
	{
		Text:    "URGENT: Your Chase bank account has been locked. Verify immediately at chase-verify-security.xyz or account closure will occur within 24 hours.",
		Country: "US",
		Status:  model.StatusPending,
		TimeAgo: 15 * time.Minute,
	},
	{
		Text:    "USPS: Package redelivery pending. Click bit.ly/usps-track-98 to pay redelivery fee and confirm 6-digit code sent to your phone.",
		Country: "US",
		Status:  model.StatusAnalyzed,
		TimeAgo: 2 * time.Hour,
	},
	{
		Text:    "Electric power shutoff notice: immediate action required. Settle past due utility balance today using an Apple gift card or prepaid card.",
		Country: "US",
		Status:  model.StatusAcknowledged,
		TimeAgo: 6 * time.Hour,
	},
	{
		Text:    "Hi mum, broke my phone this morning and dropped it in water. Text me on this new number right now, I need help with an urgent bill.",
		Country: "AU",
		Status:  model.StatusPending,
		TimeAgo: 45 * time.Minute,
	},
	{
		Text:    "HMRC final tax demand: legal action and arrest warrant pending for unpaid revenue audit. Call immediately.",
		Country: "GB",
		Status:  model.StatusResolved,
		TimeAgo: 24 * time.Hour,
	},
	{
		Text:    "Winner! You won 25,00,000 in the lucky draw lottery. Deposit 5,000 processing fee to claim prize immediately.",
		Country: "IN",
		Status:  model.StatusAnalyzed,
		TimeAgo: 4 * time.Hour,
	},
}

func main() {
	ctx := context.Background()
	log.Println("seeding fake scam incidents...")

	// check if local dev server is running
	devURL := "http://localhost:8080/api/intake"
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://localhost:8080/api/incidents")
	if err == nil && resp.StatusCode == http.StatusOK {
		seedViaHTTP(client, devURL)
		log.Println("successfully seeded incidents via local dev server")
		return
	}

	// fallback: seed directly to dynamodb if configured
	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	linksTable := os.Getenv("LINKS_TABLE")
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("note: dynamodb not configured, please start 'make dev' to seed locally: %v", err)
		return
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	st := store.NewStore(dbClient, incidentsTable, linksTable)

	for i, sc := range scenarios {
		input := risk.Input{Text: sc.Text}
		score, level, matches := risk.CalculateScore(input)
		category := risk.ClassifyCategory(input)
		output := risk.BuildSeniorOutput(category, level, sc.Country)

		createdAt := time.Now().UTC().Add(-sc.TimeAgo).Format(time.RFC3339)
		inc := &model.Incident{
			ID:          fmt.Sprintf("seed-%d", i+1),
			UserContext: sc.Text,
			Score:       score,
			RiskLevel:   level,
			Category:    category,
			Signals:     matches,
			Output:      output,
			Status:      sc.Status,
			Country:     sc.Country,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		}

		if err := st.SaveIncident(ctx, inc); err != nil {
			log.Printf("failed to save seed incident %d: %v", i+1, err)
		} else {
			log.Printf("seeded incident %d: %s [%s, %s]", i+1, inc.ID, category, level)
		}
	}
}

func seedViaHTTP(client *http.Client, url string) {
	for i, sc := range scenarios {
		payload, _ := json.Marshal(map[string]string{
			"user_context": sc.Text,
			"country":      sc.Country,
		})
		_, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("failed to seed scenario %d: %v", i+1, err)
		} else {
			log.Printf("seeded scenario %d (%s)", i+1, sc.Country)
		}
	}
}
