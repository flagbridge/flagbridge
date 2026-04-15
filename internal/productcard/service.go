package productcard

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, flagID string) (*ProductCard, error) {
	return s.repo.GetByFlagID(ctx, flagID)
}

func (s *Service) Upsert(ctx context.Context, flagID, projectID string, req UpsertRequest) (*ProductCard, error) {
	if req.Status != nil {
		if !ValidStatuses[*req.Status] {
			return nil, fmt.Errorf("invalid status: %s (valid: planning, active, rolled_out, archived)", *req.Status)
		}
	}
	return s.repo.Upsert(ctx, flagID, projectID, req)
}

func (s *Service) Delete(ctx context.Context, flagID string) error {
	return s.repo.Delete(ctx, flagID)
}
