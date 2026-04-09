package project

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

func (s *Service) Create(ctx context.Context, req CreateRequest, userID string) (*Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if !slugPattern.MatchString(req.Slug) {
		return nil, fmt.Errorf("slug must be lowercase alphanumeric with hyphens")
	}

	p := &Project{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CreatedBy:   userID,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}
	return p, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Project, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *Service) Update(ctx context.Context, slug string, req UpdateRequest) (*Project, error) {
	if req.Name == nil && req.Description == nil {
		return nil, fmt.Errorf("at least one field must be provided")
	}
	return s.repo.Update(ctx, slug, req)
}

func (s *Service) Delete(ctx context.Context, slug string) error {
	return s.repo.Delete(ctx, slug)
}

func (s *Service) List(ctx context.Context) ([]Project, error) {
	return s.repo.List(ctx)
}
