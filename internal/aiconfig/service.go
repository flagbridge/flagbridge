package aiconfig

import (
	"context"
	"fmt"

	"github.com/flagbridge/flagbridge/internal/crypto"
)

type Service struct {
	repo          *Repository
	encryptionKey []byte
}

func NewService(repo *Repository, encryptionKeyRaw string) *Service {
	return &Service{
		repo:          repo,
		encryptionKey: crypto.DeriveKey(encryptionKeyRaw),
	}
}

func (s *Service) Create(ctx context.Context, projectID string, req CreateRequest) (*Provider, error) {
	if err := validateProvider(req.Provider); err != nil {
		return nil, err
	}

	modelID := req.ModelID
	if modelID == "" {
		modelID = defaultModel(req.Provider)
	}

	var encryptedKey string
	if req.APIKey != "" {
		var err error
		encryptedKey, err = crypto.Encrypt(req.APIKey, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("encrypting API key: %w", err)
		}
	}

	maxTokens := 4096
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}

	temperature := 0.7
	if req.Temperature != nil {
		temperature = *req.Temperature
	}

	return s.repo.Upsert(ctx, projectID, req.Provider, modelID, encryptedKey, req.BaseURL, maxTokens, temperature)
}

func (s *Service) Get(ctx context.Context, projectID string) (*Provider, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

func (s *Service) Update(ctx context.Context, projectID string, req UpdateRequest) (*Provider, error) {
	existing, err := s.repo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("AI provider config not found")
	}

	provider := existing.ProviderName
	if req.Provider != nil {
		if err := validateProvider(*req.Provider); err != nil {
			return nil, err
		}
		provider = *req.Provider
	}

	modelID := existing.ModelID
	if req.ModelID != nil {
		modelID = *req.ModelID
	}

	var encryptedKey string
	if req.APIKey != nil && *req.APIKey != "" {
		encryptedKey, err = crypto.Encrypt(*req.APIKey, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("encrypting API key: %w", err)
		}
	}

	maxTokens := existing.MaxTokens
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}

	temperature := existing.Temperature
	if req.Temperature != nil {
		temperature = *req.Temperature
	}

	baseURL := existing.BaseURL
	if req.BaseURL != nil {
		baseURL = req.BaseURL
	}

	return s.repo.Upsert(ctx, projectID, provider, modelID, encryptedKey, baseURL, maxTokens, temperature)
}

func (s *Service) Delete(ctx context.Context, projectID string) error {
	return s.repo.Delete(ctx, projectID)
}

func (s *Service) DecryptAPIKey(ctx context.Context, projectID string) (string, error) {
	encryptedKey, err := s.repo.GetEncryptedKey(ctx, projectID)
	if err != nil {
		return "", err
	}
	return crypto.Decrypt(encryptedKey, s.encryptionKey)
}

func (s *Service) CheckAndIncrementUsage(ctx context.Context, projectID string) error {
	usage, limit, err := s.repo.IncrementUsage(ctx, projectID)
	if err != nil {
		return err
	}
	if usage > limit {
		return fmt.Errorf("monthly AI usage limit exceeded (%d/%d)", usage, limit)
	}
	return nil
}

func validateProvider(provider string) error {
	if !ValidProviders[provider] {
		return fmt.Errorf("invalid provider %q: must be one of anthropic, openai, ollama", provider)
	}
	return nil
}

func defaultModel(provider string) string {
	switch provider {
	case "anthropic":
		return "claude-sonnet-4-20250514"
	case "openai":
		return "gpt-4o"
	case "ollama":
		return "llama3"
	default:
		return ""
	}
}
