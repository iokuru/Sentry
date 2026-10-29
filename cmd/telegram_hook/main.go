package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/store"
	"github.com/iokuru/sentry/internal/telegram"
)

var (
	incidentStore  *store.Store
	telegramClient *telegram.Client
	telegramCfg    *telegram.Config
)

func init() {
	incidentsTable := os.Getenv("INCIDENTS_TABLE")
	linksTable := os.Getenv("LINKS_TABLE")
	if incidentsTable == "" {
		incidentsTable = "sentry-incidents"
	}
	if linksTable == "" {
		linksTable = "sentry-links"
	}

	telegramClient = telegram.NewClient()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err == nil {
		dbClient := dynamodb.NewFromConfig(cfg)
		incidentStore = store.NewStore(dbClient, incidentsTable, linksTable)
		ssmClient := ssm.NewFromConfig(cfg)
		telegramCfg = telegram.LoadConfigFromSSM(context.Background(), ssmClient)
	}
}

func handleWebhook(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if telegramCfg == nil || !telegramCfg.Enabled {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Body:       `{"ok":false,"reason":"telegram disabled"}`,
		}, nil
	}

	// constant time secret header verification
	headerSecret := req.Headers["x-telegram-bot-api-secret-token"]
	if headerSecret == "" {
		headerSecret = req.Headers["X-Telegram-Bot-Api-Secret-Token"]
	}

	if !telegram.VerifyWebhookSecret(headerSecret, telegramCfg.WebhookSecret) {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"error":"unauthorized"}`,
		}, nil
	}

	var update telegram.Update
	if err := json.Unmarshal([]byte(req.Body), &update); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"invalid json"}`,
		}, nil
	}

	if update.CallbackQuery != nil {
		handleCallback(ctx, update.CallbackQuery)
	} else if update.Message != nil {
		handleMessage(ctx, update.Message)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"ok":true}`,
	}, nil
}

func handleCallback(ctx context.Context, cb *telegram.CallbackQuery) {
	action, incidentID, ok := telegram.ParseCallbackData(cb.Data)
	if !ok {
		return
	}

	var newStatus model.IncidentStatus
	var statusLabel string
	if action == "ack" {
		newStatus = model.StatusAcknowledged
		statusLabel = "Acknowledged"
	} else if action == "res" {
		newStatus = model.StatusResolved
		statusLabel = "Resolved"
	} else {
		return
	}

	if incidentStore != nil {
		if err := incidentStore.UpdateIncidentStatus(ctx, incidentID, newStatus); err != nil {
			log.Printf("failed to update incident %s: %v", incidentID, err)
		}
	}

	// answer callback query to dismiss loading state
	_ = telegramClient.AnswerCallback(
		ctx,
		telegramCfg.Token,
		cb.ID,
		fmt.Sprintf("Incident %s marked as %s", incidentID, statusLabel),
	)

	// edit original message to update status and button state
	if cb.Message != nil {
		chatID := fmt.Sprintf("%d", cb.Message.Chat.ID)
		messageID := cb.Message.MessageID
		originalText := cb.Message.Text

		var updatedText string
		if strings.Contains(originalText, "\n\n*Status:*") {
			parts := strings.Split(originalText, "\n\n*Status:*")
			updatedText = parts[0] + "\n\n*Status:* " + statusLabel
		} else {
			updatedText = originalText + "\n\n*Status:* " + statusLabel
		}

		var remainingButtons [][]telegram.InlineKeyboardButton
		if action == "ack" {
			remainingButtons = [][]telegram.InlineKeyboardButton{
				{
					{Text: "Resolve", CallbackData: fmt.Sprintf("res:%s", incidentID)},
				},
			}
		}

		_ = telegramClient.EditMessage(ctx, telegramCfg.Token, chatID, messageID, updatedText, remainingButtons)
	}
}

func handleMessage(ctx context.Context, msg *telegram.Message) {
	if msg == nil || telegramCfg == nil || telegramClient == nil {
		return
	}

	chatID := fmt.Sprintf("%d", msg.Chat.ID)
	text := strings.TrimSpace(msg.Text)

	if !strings.HasPrefix(text, "/start") {
		return
	}

	parts := strings.Fields(text)
	if len(parts) < 2 {
		_, _ = telegramClient.Send(ctx, telegramCfg.Token, chatID, "Welcome to Sentry. Send /start <CODE> using your invite code to link your caregiver account.", nil)
		return
	}

	code := strings.ToUpper(strings.TrimSpace(parts[1]))
	if incidentStore == nil {
		_, _ = telegramClient.Send(ctx, telegramCfg.Token, chatID, "Sentry service temporarily unavailable.", nil)
		return
	}

	inv, err := incidentStore.GetInvite(ctx, code)
	if err != nil || inv == nil {
		_, _ = telegramClient.Send(ctx, telegramCfg.Token, chatID, "Invite code not found or expired. Please generate a new invite code in your caregiver portal.", nil)
		return
	}

	if err := incidentStore.SetTelegramChatID(ctx, inv.CaregiverID, chatID); err != nil {
		log.Printf("failed to associate telegram chat id: %v", err)
		_, _ = telegramClient.Send(ctx, telegramCfg.Token, chatID, "Failed to link Telegram. Please try again.", nil)
		return
	}

	displayName := inv.SeniorName
	if displayName == "" {
		displayName = "your senior"
	}

	confirmMsg := fmt.Sprintf("Sentry linked. You will receive high risk alerts for %s in this chat.", displayName)
	_, _ = telegramClient.Send(ctx, telegramCfg.Token, chatID, confirmMsg, nil)
}

func main() {
	lambda.Start(handleWebhook)
}
