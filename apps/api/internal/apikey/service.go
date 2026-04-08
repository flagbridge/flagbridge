package apikey

import (
	"context"
	"fmt"

	"github.com/flagbridge/flagbridge/apps/api/internal/auth"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest, userID string) (*CreateResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.ProjectID == "" {
		return nil, fmt.Errorf("project_id is required")
	}

	scope := req.Scope
	if scope == "" {
		scope = "eval"
	}

	validScopes := map[string]bool{"eval": true, "test": true, "mgmt": true, "full": true}
	if !validScopes[scope] {
		return nil, fmt.Errorf("invalid scope %q: must be one of eval, test, mgmt, full", scope)
	}

	fullKey, hash, prefix, err := auth.GenerateAPIKey(scope)
	if err != nil {
		return nil, fmt.Errorf("generating API key: %w", err)
	}

	ak, err := s.repo.Create(ctx, req.Name, hash, prefix, scope, req.ProjectID, req.EnvironmentID, userID)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		APIKey:  *ak,
		FullKey: fullKey,
	}, nil
}

func (s *Service) List(ctx context.Context) ([]APIKey, error) {
	return s.repo.List(ctx)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
