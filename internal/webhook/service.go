package webhook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
)

type Service struct {
	repo       *Repository
	dispatcher *Dispatcher
}

func NewService(repo *Repository, dispatcher *Dispatcher) *Service {
	return &Service{repo: repo, dispatcher: dispatcher}
}

func (s *Service) Create(ctx context.Context, projectID string, req CreateRequest, userID string) (*Webhook, error) {
	if req.URL == "" {
		return nil, fmt.Errorf("url is required")
	}
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if len(req.Events) == 0 {
		return nil, fmt.Errorf("at least one event is required")
	}

	// Validate event types
	validEvents := make(map[string]bool, len(AllEvents))
	for _, e := range AllEvents {
		validEvents[e] = true
	}
	for _, e := range req.Events {
		if !validEvents[e] {
			return nil, fmt.Errorf("invalid event type: %s", e)
		}
	}

	secret := req.Secret
	if secret == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return nil, fmt.Errorf("generating secret: %w", err)
		}
		secret = "whsec_" + hex.EncodeToString(b)
	}

	wh := &Webhook{
		ProjectID: projectID,
		URL:       req.URL,
		Secret:    secret,
		Events:    req.Events,
		Active:    true,
		CreatedBy: userID,
	}

	if err := s.repo.Create(ctx, wh); err != nil {
		return nil, fmt.Errorf("creating webhook: %w", err)
	}
	return wh, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Webhook, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByProject(ctx context.Context, projectID string) ([]Webhook, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Webhook, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) ListDeliveryLogs(ctx context.Context, webhookID string, limit int) ([]DeliveryLog, error) {
	return s.repo.ListDeliveryLogs(ctx, webhookID, limit)
}

// TestWebhook sends a test event to a webhook to verify it's working.
func (s *Service) TestWebhook(ctx context.Context, id string) error {
	wh, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	s.dispatcher.Dispatch(wh.ProjectID, "webhook.test", map[string]string{
		"webhook_id": wh.ID,
		"message":    "This is a test webhook delivery from FlagBridge",
	})
	return nil
}

// Dispatcher returns the dispatcher for use by other packages to fire events.
func (s *Service) Dispatcher() *Dispatcher {
	return s.dispatcher
}
