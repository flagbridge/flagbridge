package flag

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
)

var keyPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, projectID string, req CreateRequest, userID string) (*Flag, error) {
	if !keyPattern.MatchString(req.Key) {
		return nil, fmt.Errorf("key must be lowercase alphanumeric with hyphens")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	flagType := req.Type
	if flagType == "" {
		flagType = "boolean"
	}

	defaultValue := req.DefaultValue
	if defaultValue == nil {
		defaultValue = json.RawMessage(`false`)
	}

	f := &Flag{
		ProjectID:    projectID,
		Key:          req.Key,
		Name:         req.Name,
		Description:  req.Description,
		Type:         flagType,
		DefaultValue: defaultValue,
		Tags:         req.Tags,
		CreatedBy:    userID,
	}
	if err := s.repo.Create(ctx, f); err != nil {
		return nil, fmt.Errorf("creating flag: %w", err)
	}
	return f, nil
}

func (s *Service) GetByKey(ctx context.Context, projectID, key string) (*Flag, error) {
	return s.repo.GetByKey(ctx, projectID, key)
}

func (s *Service) ListByProject(ctx context.Context, projectID string) ([]Flag, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *Service) Update(ctx context.Context, projectID, key string, req UpdateRequest) (*Flag, error) {
	return s.repo.Update(ctx, projectID, key, req)
}

func (s *Service) Delete(ctx context.Context, projectID, key string) error {
	return s.repo.Delete(ctx, projectID, key)
}

func (s *Service) GetState(ctx context.Context, flagID, envID string) (*FlagState, error) {
	return s.repo.GetState(ctx, flagID, envID)
}

func (s *Service) SetState(ctx context.Context, flagID, envID, userID string, req SetStateRequest) (*FlagState, error) {
	return s.repo.SetState(ctx, flagID, envID, userID, req)
}
