package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/ports"
)

type Client struct {
	enabled  bool
	botToken string
	chatID   string
	http     *http.Client
}

func NewClient(enabled bool, botToken, chatID string, timeout time.Duration) *Client {
	return &Client{
		enabled:  enabled,
		botToken: strings.TrimSpace(botToken),
		chatID:   strings.TrimSpace(chatID),
		http:     &http.Client{Timeout: timeout},
	}
}

func (c *Client) Send(ctx context.Context, msg ports.TelegramMessage) error {
	if !c.enabled {
		return nil
	}
	if c.botToken == "" || c.chatID == "" {
		return nil
	}

	if strings.TrimSpace(msg.PhotoURL) != "" {
		if err := c.sendPhoto(ctx, msg); err == nil {
			return nil
		}
		// fallback a mensaje de texto
	}

	return c.sendMessage(ctx, msg)
}

func (c *Client) sendPhoto(ctx context.Context, msg ports.TelegramMessage) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", c.botToken)
	payload := map[string]any{
		"chat_id": c.chatID,
		"photo":   msg.PhotoURL,
		"caption": truncate(msg.Text, 1024),
	}
	if msg.ParseMode != "" {
		payload["parse_mode"] = msg.ParseMode
	}
	return c.postJSON(ctx, endpoint, payload)
}

func (c *Client) sendMessage(ctx context.Context, msg ports.TelegramMessage) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken)
	payload := map[string]any{
		"chat_id":                  c.chatID,
		"text":                     truncate(msg.Text, 4096),
		"disable_web_page_preview": msg.DisablePreview,
	}
	if msg.ParseMode != "" {
		payload["parse_mode"] = msg.ParseMode
	}
	return c.postJSON(ctx, endpoint, payload)
}

func (c *Client) postJSON(ctx context.Context, endpoint string, payload map[string]any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	return fmt.Errorf("telegram http error status=%d body=%s", resp.StatusCode, string(respBody))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
