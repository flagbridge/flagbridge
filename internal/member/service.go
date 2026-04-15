package member

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

func (s *Service) List(ctx context.Context, projectID string) ([]Member, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *Service) Add(ctx context.Context, projectID, email, role string) (*Member, error) {
	if !ValidRoles[role] {
		return nil, fmt.Errorf("invalid role: %s (valid: admin, engineer, product, viewer)", role)
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	userID, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return s.repo.Add(ctx, projectID, userID, role)
}

func (s *Service) UpdateRole(ctx context.Context, memberID, role string) (*Member, error) {
	if !ValidRoles[role] {
		return nil, fmt.Errorf("invalid role: %s (valid: admin, engineer, product, viewer)", role)
	}
	return s.repo.UpdateRole(ctx, memberID, role)
}

func (s *Service) Remove(ctx context.Context, memberID string) error {
	return s.repo.Remove(ctx, memberID)
}

// AddCreatorAsAdmin is called when a project is created to auto-add the creator.
func (s *Service) AddCreatorAsAdmin(ctx context.Context, projectID, userID string) error {
	_, err := s.repo.Add(ctx, projectID, userID, "admin")
	return err
}
