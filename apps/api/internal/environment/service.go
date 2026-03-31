package environment

import (
	"context"
	"fmt"
	"regexp"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, projectID string, req CreateRequest) (*Environment, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if !slugPattern.MatchString(req.Slug) {
		return nil, fmt.Errorf("slug must be lowercase alphanumeric with hyphens")
	}

	e := &Environment{
		ProjectID: projectID,
		Name:      req.Name,
		Slug:      req.Slug,
		Color:     req.Color,
		SortOrder: req.SortOrder,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("creating environment: %w", err)
	}
	return e, nil
}

func (s *Service) ListByProject(ctx context.Context, projectID string) ([]Environment, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *Service) GetBySlug(ctx context.Context, projectID, slug string) (*Environment, error) {
	return s.repo.GetBySlug(ctx, projectID, slug)
}
