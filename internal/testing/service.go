package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSession(ctx context.Context, req CreateSessionRequest, userID string) (*Session, error) {
	if req.ProjectID == "" {
		return nil, fmt.Errorf("project_id is required")
	}

	ttl := req.TTLSecs
	if ttl <= 0 {
		ttl = 3600 // 1 hour default
	}
	if ttl > 86400 {
		ttl = 86400 // max 24 hours
	}

	session := &Session{
		ProjectID: req.ProjectID,
		Label:     req.Label,
		Overrides: make(map[string]Override),
		CreatedBy: userID,
		ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Second),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}
	return session, nil
}

func (s *Service) GetSession(ctx context.Context, id string) (*Session, error) {
	return s.repo.GetSession(ctx, id)
}

func (s *Service) DeleteSession(ctx context.Context, id string) error {
	return s.repo.DeleteSession(ctx, id)
}

func (s *Service) SetOverride(ctx context.Context, sessionID string, req SetOverrideRequest) (*Session, error) {
	if req.FlagKey == "" {
		return nil, fmt.Errorf("flag_key is required")
	}

	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	session.Overrides[req.FlagKey] = Override(req)

	if err := s.repo.UpdateOverrides(ctx, sessionID, session.Overrides); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) SetOverridesBatch(ctx context.Context, sessionID string, req SetOverridesBatchRequest) (*Session, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	for _, o := range req.Overrides {
		if o.FlagKey == "" {
			continue
		}
		session.Overrides[o.FlagKey] = Override(o)
	}

	if err := s.repo.UpdateOverrides(ctx, sessionID, session.Overrides); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) DeleteOverride(ctx context.Context, sessionID, flagKey string) (*Session, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	delete(session.Overrides, flagKey)

	if err := s.repo.UpdateOverrides(ctx, sessionID, session.Overrides); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) ListSessions(ctx context.Context, projectID string) ([]Session, error) {
	return s.repo.ListSessions(ctx, projectID)
}

// ResolveOverride looks up the override value for a flag in a session.
// Returns nil if no override is set for this flag.
func (s *Service) ResolveOverride(ctx context.Context, sessionID, flagKey string) (*json.RawMessage, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, nil // session not found/expired → no override, don't error
	}

	override, ok := session.Overrides[flagKey]
	if !ok {
		return nil, nil
	}
	return &override.Value, nil
}
