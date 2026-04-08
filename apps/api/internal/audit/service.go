package audit

import (
	"context"
	"log/slog"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Log appends an audit entry asynchronously. Never fails the caller.
func (s *Service) Log(ctx context.Context, input LogInput) {
	go func() {
		if err := s.repo.Append(context.Background(), input); err != nil {
			slog.Error("failed to write audit log", "error", err, "action", input.Action, "entity", input.EntityType)
		}
	}()
}

// LogSync appends an audit entry synchronously.
func (s *Service) LogSync(ctx context.Context, input LogInput) error {
	return s.repo.Append(ctx, input)
}

func (s *Service) List(ctx context.Context, params ListParams) ([]Entry, int, error) {
	return s.repo.List(ctx, params)
}
