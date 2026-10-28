package telegram

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type SendMessageRequest struct {
	ChatID      string                `json:"chat_id"`
	Text        string                `json:"text"`
	ParseMode   string                `json:"parse_mode,omitempty"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type EditMessageTextRequest struct {
	ChatID      string                `json:"chat_id"`
	MessageID   int                   `json:"message_id"`
	Text        string                `json:"text"`
	ParseMode   string                `json:"parse_mode,omitempty"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type AnswerCallbackQueryRequest struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
}

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Message *Message `json:"message,omitempty"`
	Data    string   `json:"data"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	Chat      Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			// 3 second timeout for telegram requests
			Timeout: 3 * time.Second,
		},
		baseURL: "https://api.telegram.org",
	}
}

func (c *Client) Send(ctx context.Context, token, chatID, text string, buttons [][]InlineKeyboardButton) (int, error) {
	url := fmt.Sprintf("%s/bot%s/sendMessage", c.baseURL, token)

	reqPayload := SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	}
	if len(buttons) > 0 {
		reqPayload.ReplyMarkup = &InlineKeyboardMarkup{InlineKeyboard: buttons}
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return 0, fmt.Errorf("marshal send request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("create send request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("post telegram message: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return 0, fmt.Errorf("telegram returned status %d", res.StatusCode)
	}

	var apiRes struct {
		OK     bool    `json:"ok"`
		Result Message `json:"result"`
	}
	if err := json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		return 0, fmt.Errorf("decode telegram response: %w", err)
	}

	return apiRes.Result.MessageID, nil
}

func (c *Client) EditMessage(ctx context.Context, token, chatID string, messageID int, text string, buttons [][]InlineKeyboardButton) error {
	url := fmt.Sprintf("%s/bot%s/editMessageText", c.baseURL, token)

	reqPayload := EditMessageTextRequest{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      text,
		ParseMode: "Markdown",
	}
	if len(buttons) > 0 {
		reqPayload.ReplyMarkup = &InlineKeyboardMarkup{InlineKeyboard: buttons}
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return fmt.Errorf("marshal edit request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create edit request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post telegram edit: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("telegram returned status %d", res.StatusCode)
	}

	return nil
}

func (c *Client) AnswerCallback(ctx context.Context, token, callbackQueryID, text string) error {
	url := fmt.Sprintf("%s/bot%s/answerCallbackQuery", c.baseURL, token)

	reqPayload := AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            text,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		return fmt.Errorf("marshal callback response: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create callback request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post callback answer: %w", err)
	}
	defer res.Body.Close()

	return nil
}

func FormatAlertMessage(category, level, detectedOneLine, hotlineName, hotlineNumber, incidentID, timestamp string) (string, [][]InlineKeyboardButton) {
	text := fmt.Sprintf(
		"*SENTRY ALERT*\n\n"+
			"*Severity:* %s\n"+
			"*Category:* %s\n"+
			"*Detected:* %s\n"+
			"*Hotline:* %s at %s\n"+
			"*Incident:* `%s`\n"+
			"*Time:* %s",
		strings.ToUpper(level),
		category,
		detectedOneLine,
		hotlineName,
		hotlineNumber,
		incidentID,
		timestamp,
	)

	buttons := [][]InlineKeyboardButton{
		{
			{Text: "Acknowledge", CallbackData: fmt.Sprintf("ack:%s", incidentID)},
			{Text: "Resolve", CallbackData: fmt.Sprintf("res:%s", incidentID)},
		},
	}

	return text, buttons
}

func ParseCallbackData(data string) (action, incidentID string, ok bool) {
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	action = parts[0]
	incidentID = parts[1]
	if action != "ack" && action != "res" {
		return "", "", false
	}
	return action, incidentID, true
}

func VerifyWebhookSecret(headerSecret, expectedSecret string) bool {
	if headerSecret == "" || expectedSecret == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(headerSecret), []byte(expectedSecret)) == 1
}
