package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/risk"
	"github.com/iokuru/sentry/internal/store"
)

type Scenario struct {
	Name             string
	Category         model.Category
	Country          string
	Sender           string
	Text             string
	Filename         string
	Status           model.IncidentStatus
	MinExpectedScore int
	MaxExpectedScore int
}

var allScenarios = []Scenario{
	{
		Name:             "Chase Bank Account Suspension",
		Category:         model.CategoryBank,
		Country:          "US",
		Sender:           "Chase Security",
		Text:             "URGENT: Chase bank security alert. Your debit card has been locked due to unauthorized activity. Call 1-800-555-0199 or verify your PIN now at chase-verify-security.xyz to avoid permanent account closure.",
		Filename:         "bank_scam.png",
		Status:           model.StatusAnalyzed,
		MinExpectedScore: 60,
		MaxExpectedScore: 100,
	},
	{
		Name:             "USPS Package Redelivery",
		Category:         model.CategoryDelivery,
		Country:          "US",
		Sender:           "USPS Tracking",
		Text:             "USPS: Your package could not be delivered due to incomplete address. Please update tracking details at bit.ly/usps-redelivery-99 and enter the one-time verification code to pay $1.99 redelivery fee.",
		Filename:         "delivery_scam.png",
		Status:           model.StatusAnalyzed,
		MinExpectedScore: 30,
		MaxExpectedScore: 80,
	},
	{
		Name:             "Electric Utility Disconnection",
		Category:         model.CategoryUtility,
		Country:          "US",
		Sender:           "Electric Power Dept",
		Text:             "FINAL NOTICE: Electric service will be terminated today in 30 minutes due to unpaid utility bill of $428. Settle balance immediately using an Apple gift card or prepaid card.",
		Filename:         "utility_scam.png",
		Status:           model.StatusAcknowledged,
		MinExpectedScore: 30,
		MaxExpectedScore: 80,
	},
	{
		Name:             "Lottery Sweepstakes Winner",
		Category:         model.CategoryPrize,
		Country:          "US",
		Sender:           "Mega Millions Lottery",
		Text:             "Official Sweepstakes Winner! You have won $1,250,000 in the lucky draw lottery. Send $350 processing fee via Bitcoin or cryptocurrency to claim your prize jackpot.",
		Filename:         "prize_scam.png",
		Status:           model.StatusAnalyzed,
		MinExpectedScore: 30,
		MaxExpectedScore: 80,
	},
	{
		Name:             "Grandson Bail Emergency",
		Category:         model.CategoryRelativeInNeed,
		Country:          "AU",
		Sender:           "Unknown Number",
		Text:             "Hi mum, broke my phone this morning and got arrested for an accident. Text me on this new number right now, I need help with urgent bail money, please wire $1500.",
		Filename:         "relative_scam.png",
		Status:           model.StatusPending,
		MinExpectedScore: 50,
		MaxExpectedScore: 100,
	},
	{
		Name:             "IRS Warrant and Arrest Demand",
		Category:         model.CategoryTax,
		Country:          "US",
		Sender:           "IRS Legal Enforcement",
		Text:             "IRS final warning: A federal warrant has been issued for your arrest due to unpaid back taxes. Federal revenue officers are en route to your residence. Immediate payment required.",
		Filename:         "tax_scam.png",
		Status:           model.StatusResolved,
		MinExpectedScore: 30,
		MaxExpectedScore: 80,
	},
	{
		Name:             "Online Romance Customs Fee",
		Category:         model.CategoryRomance,
		Country:          "US",
		Sender:           "Soulmate Flight",
		Text:             "My darling soulmate, I want to meet you in person but I am stuck at customs on my flight. They are demanding $800 wire transfer for luggage release before letting me see you.",
		Filename:         "romance_scam.png",
		Status:           model.StatusAnalyzed,
		MinExpectedScore: 30,
		MaxExpectedScore: 80,
	},
	{
		Name:             "Microsoft Trojan Defender Spyware",
		Category:         model.CategoryTechSupport,
		Country:          "GB",
		Sender:           "Windows Defender Security",
		Text:             "Microsoft Windows Defender Alert: Trojan spyware detected on your computer. Your bank logins and passwords are compromised. Immediate action required, call support at 1-800-555-0144 to verify one-time code.",
		Filename:         "tech_support_scam.png",
		Status:           model.StatusAnalyzed,
		MinExpectedScore: 50,
		MaxExpectedScore: 100,
	},
}

func main() {
	if err := os.MkdirAll("fixtures", 0755); err != nil {
		log.Fatalf("failed to create fixtures dir: %v", err)
	}

	fmt.Println("Generating scam screenshot fixtures in fixtures/ ...")
	for _, sc := range allScenarios {
		filePath := filepath.Join("fixtures", sc.Filename)
		if err := generateScreenshotPNG(sc.Sender, sc.Text, filePath); err != nil {
			log.Printf("error generating %s: %v", filePath, err)
		}
	}

	fmt.Println("\nRunning fixtures through scoring engine...")
	fmt.Printf("%-18s | %-8s | %-6s | %-24s | %-22s\n", "Category", "Level", "Score", "Signals Matched", "Fixture File")
	fmt.Println(strings.Repeat("-", 88))

	var failedScenarios int
	for _, sc := range allScenarios {
		input := risk.Input{Text: sc.Text}
		score, level, matches := risk.CalculateScore(input)
		detectedCategory := risk.ClassifyCategory(input)

		var sigNames []string
		for _, m := range matches {
			sigNames = append(sigNames, m.Name)
		}
		sigSummary := strings.Join(sigNames, ", ")
		if len(sigSummary) > 24 {
			sigSummary = sigSummary[:21] + "..."
		}

		fmt.Printf("%-18s | %-8s | %-6d | %-24s | %-22s\n",
			detectedCategory, strings.ToUpper(string(level)), score, sigSummary, sc.Filename)

		if detectedCategory != sc.Category {
			fmt.Printf("  WARNING: Expected category %s, got %s\n", sc.Category, detectedCategory)
			failedScenarios++
		}
		if score < sc.MinExpectedScore || score > sc.MaxExpectedScore {
			fmt.Printf("  WARNING: Score %d out of expected range [%d, %d]\n", score, sc.MinExpectedScore, sc.MaxExpectedScore)
			failedScenarios++
		}
	}
	fmt.Println(strings.Repeat("-", 88))

	if failedScenarios == 0 {
		fmt.Println("All 8 categories verified deterministically by scoring engine.")
	}

	seedToRunningEnvironment()
}

func generateScreenshotPNG(sender, messageText, outPath string) error {
	const (
		width  = 400
		height = 680
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	bgColor := color.RGBA{R: 242, G: 242, B: 247, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	statusBar := image.Rect(0, 0, width, 36)
	draw.Draw(img, statusBar, &image.Uniform{color.RGBA{R: 44, G: 44, B: 46, A: 255}}, image.Point{}, draw.Src)

	appHeader := image.Rect(0, 36, width, 90)
	draw.Draw(img, appHeader, &image.Uniform{color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.Point{}, draw.Src)
	drawBorder(img, appHeader, color.RGBA{R: 209, G: 209, B: 214, A: 255})

	bubble := image.Rect(20, 110, width-20, 520)
	draw.Draw(img, bubble, &image.Uniform{color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.Point{}, draw.Src)
	drawBorder(img, bubble, color.RGBA{R: 200, G: 200, B: 205, A: 255})

	senderBar := image.Rect(20, 110, width-20, 150)
	draw.Draw(img, senderBar, &image.Uniform{color.RGBA{R: 225, G: 29, B: 72, A: 255}}, image.Point{}, draw.Src)

	lines := wrapText(messageText, 38)
	lineY := 170
	for _, line := range lines {
		if lineY > 500 {
			break
		}
		drawTextLine(img, 36, lineY, line)
		lineY += 26
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create png: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	return nil
}

func drawBorder(img *image.RGBA, r image.Rectangle, col color.RGBA) {
	for x := r.Min.X; x < r.Max.X; x++ {
		img.Set(x, r.Min.Y, col)
		img.Set(x, r.Max.Y-1, col)
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		img.Set(r.Min.X, y, col)
		img.Set(r.Max.X-1, y, col)
	}
}

func drawTextLine(img *image.RGBA, startX, startY int, text string) {
	col := color.RGBA{R: 28, G: 25, B: 23, A: 255}
	x := startX
	for _, ch := range text {
		if ch == ' ' {
			x += 6
			continue
		}
		for dx := 0; dx < 4; dx++ {
			for dy := 0; dy < 6; dy++ {
				if (dx+dy)%2 == 0 || dx == 0 || dy == 0 {
					img.Set(x+dx, startY+dy, col)
				}
			}
		}
		x += 7
		if x > 360 {
			break
		}
	}
}

func wrapText(text string, maxChars int) []string {
	words := strings.Fields(text)
	var lines []string
	var cur strings.Builder

	for _, w := range words {
		if cur.Len()+len(w)+1 > maxChars {
			lines = append(lines, cur.String())
			cur.Reset()
		}
		if cur.Len() > 0 {
			cur.WriteString(" ")
		}
		cur.WriteString(w)
	}
	if cur.Len() > 0 {
		lines = append(lines, cur.String())
	}
	return lines
}

func seedToRunningEnvironment() {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://localhost:8080/api/incidents")
	if err == nil && resp.StatusCode == http.StatusOK {
		fmt.Println("\nSeeding incidents to local devserver at http://localhost:8080 ...")
		for _, sc := range allScenarios {
			intakePayload := map[string]string{
				"image_key":    "fixtures/" + sc.Filename,
				"user_context": sc.Text,
				"country":      sc.Country,
			}
			body, _ := json.Marshal(intakePayload)
			intakeResp, postErr := client.Post("http://localhost:8080/api/intake", "application/json", bytes.NewReader(body))
			if postErr == nil {
				intakeResp.Body.Close()
			}
		}
		fmt.Println("Seeded all scenarios into local devserver memory store.")
		return
	}

	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	if incidentsTable == "" {
		return
	}
	linksTable := os.Getenv("LINKS_TABLE")
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	st := store.NewStore(dbClient, incidentsTable, linksTable)

	for i, sc := range allScenarios {
		input := risk.Input{Text: sc.Text}
		score, level, matches := risk.CalculateScore(input)
		category := risk.ClassifyCategory(input)
		output := risk.BuildSeniorOutput(category, level, sc.Country)

		inc := &model.Incident{
			ID:          fmt.Sprintf("seed-inc-%d", i+1),
			Score:       score,
			RiskLevel:   level,
			Category:    category,
			Signals:     matches,
			Output:      output,
			Status:      sc.Status,
			Country:     sc.Country,
			UserContext: sc.Text,
			CreatedAt:   time.Now().UTC().Add(-time.Duration(i*2) * time.Hour).Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		}
		_ = st.SaveIncident(ctx, inc)
	}
	fmt.Println("Seeded incidents directly to DynamoDB tables.")
}
