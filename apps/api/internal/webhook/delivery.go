package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Dispatcher handles asynchronous webhook delivery with retries.
type Dispatcher struct {
	repo   *Repository
	client *http.Client
}

func NewDispatcher(repo *Repository) *Dispatcher {
	return &Dispatcher{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Dispatch sends a webhook event to all matching webhooks asynchronously.
func (d *Dispatcher) Dispatch(projectID, eventType string, payload any) {
	go func() {
		ctx := context.Background()

		webhooks, err := d.repo.FindByEvent(ctx, projectID, eventType)
		if err != nil {
			slog.Error("failed to find webhooks", "error", err, "event", eventType)
			return
		}

		for _, wh := range webhooks {
			go d.deliver(ctx, wh, eventType, payload)
		}
	}()
}

func (d *Dispatcher) deliver(ctx context.Context, wh Webhook, eventType string, payload any) {
	body, err := json.Marshal(map[string]any{
		"event":      eventType,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"webhook_id": wh.ID,
		"data":       payload,
	})
	if err != nil {
		slog.Error("failed to marshal webhook payload", "error", err)
		return
	}

	maxRetries := 5
	var lastStatusCode int
	var lastResponse string

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, wh.URL, bytes.NewReader(body))
		if err != nil {
			slog.Error("failed to create webhook request", "error", err, "url", wh.URL)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "FlagBridge-Webhook/1.0")
		req.Header.Set("X-FlagBridge-Event", eventType)
		req.Header.Set("X-FlagBridge-Delivery", wh.ID)

		// HMAC-SHA256 signature
		if wh.Secret != "" {
			sig := signPayload(wh.Secret, body)
			req.Header.Set("X-FlagBridge-Signature", "sha256="+sig)
		}

		resp, err := d.client.Do(req)
		if err != nil {
			lastStatusCode = 0
			lastResponse = err.Error()
			slog.Warn("webhook delivery failed", "attempt", attempt, "url", wh.URL, "error", err)
			backoff(attempt)
			continue
		}

		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()

		lastStatusCode = resp.StatusCode
		lastResponse = string(respBody)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success — log and return
			d.logDelivery(ctx, wh.ID, eventType, string(body), lastStatusCode, lastResponse, attempt, true)
			return
		}

		slog.Warn("webhook delivery non-2xx", "attempt", attempt, "url", wh.URL, "status", resp.StatusCode)
		backoff(attempt)
	}

	// All retries failed
	d.logDelivery(ctx, wh.ID, eventType, string(body), lastStatusCode, lastResponse, maxRetries, false)
	slog.Error("webhook delivery exhausted retries", "webhook_id", wh.ID, "url", wh.URL, "event", eventType)
}

func (d *Dispatcher) logDelivery(ctx context.Context, webhookID, eventType, payload string, statusCode int, response string, attempts int, success bool) {
	log := &DeliveryLog{
		WebhookID:  webhookID,
		EventType:  eventType,
		Payload:    payload,
		StatusCode: statusCode,
		Response:   response,
		Attempts:   attempts,
		Success:    success,
	}
	if err := d.repo.LogDelivery(ctx, log); err != nil {
		slog.Error("failed to log webhook delivery", "error", err)
	}
}

func signPayload(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func backoff(attempt int) {
	// Exponential backoff: 1s, 2s, 4s, 8s, 16s
	d := time.Duration(1<<uint(attempt-1)) * time.Second
	if d > 16*time.Second {
		d = 16 * time.Second
	}
	time.Sleep(d)
}

// VerifySignature checks the HMAC-SHA256 signature of an incoming webhook payload.
// Useful for documentation/SDK examples.
func VerifySignature(secret string, body []byte, signature string) bool {
	expected := signPayload(secret, body)
	return fmt.Sprintf("sha256=%s", expected) == signature
}
