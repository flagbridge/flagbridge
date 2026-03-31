package targeting

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetRules(ctx context.Context, flagID, envID string) ([]Rule, error) {
	return s.repo.GetByFlagAndEnv(ctx, flagID, envID)
}

func (s *Service) SetRules(ctx context.Context, flagID, envID string, inputs []RuleInput) ([]Rule, error) {
	return s.repo.SetRules(ctx, flagID, envID, inputs)
}
