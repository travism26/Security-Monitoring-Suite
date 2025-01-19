package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/travism26/log-aggregator/internal/domain"
)

// APIKeyService handles API key operations
type APIKeyService struct {
	repo domain.APIKeyRepository
}

// NewAPIKeyService creates a new APIKeyService instance
func NewAPIKeyService(repo domain.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		repo: repo,
	}
}

// ValidateKey validates an API key and returns the associated organization ID
func (s *APIKeyService) ValidateKey(keyHash string) (*domain.APIKey, error) {
	apiKey, err := s.repo.GetByHash(keyHash)
	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}

	if apiKey.Status != domain.APIKeyStatusActive {
		return nil, fmt.Errorf("API key is not active")
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("API key has expired")
	}

	// Update last used timestamp
	err = s.repo.UpdateLastUsed(apiKey.ID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update last used timestamp for API key %s: %v\n", apiKey.ID, err)
	}

	return apiKey, nil
}

// GenerateKey generates a new API key
func (s *APIKeyService) GenerateKey(orgID uuid.UUID, name string, keyType domain.APIKeyType) (string, error) {
	// Generate a random UUID for the key
	keyString := uuid.New().String()

	// Hash the key for storage
	hash := sha256.Sum256([]byte(keyString))
	keyHash := hex.EncodeToString(hash[:])

	apiKey := &domain.APIKey{
		ID:             uuid.New(),
		OrganizationID: orgID,
		KeyType:        keyType,
		KeyHash:        keyHash,
		Name:           name,
		CreatedAt:      time.Now().UTC(),
		Status:         domain.APIKeyStatusActive,
	}

	err := s.repo.Create(apiKey)
	if err != nil {
		return "", fmt.Errorf("failed to store API key: %w", err)
	}

	return keyString, nil
}

// RevokeKey revokes an API key
func (s *APIKeyService) RevokeKey(id uuid.UUID) error {
	return s.repo.RevokeKey(id)
}

// ListKeys lists all API keys for an organization
func (s *APIKeyService) ListKeys(orgID uuid.UUID, limit, offset int) ([]*domain.APIKey, error) {
	return s.repo.ListByOrganization(orgID, limit, offset)
}

// GetKey gets an API key by ID
func (s *APIKeyService) GetKey(id uuid.UUID) (*domain.APIKey, error) {
	return s.repo.GetByID(id)
}

// HashKey hashes an API key string
func (s *APIKeyService) HashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
