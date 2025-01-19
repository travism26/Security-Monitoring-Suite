package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// APIKeyType represents the type of API key
type APIKeyType string

const (
	APIKeyTypeAgent    APIKeyType = "agent"
	APIKeyTypeCustomer APIKeyType = "customer"
)

// APIKeyStatus represents the possible states of an API key
type APIKeyStatus string

const (
	APIKeyStatusActive  APIKeyStatus = "active"
	APIKeyStatusRevoked APIKeyStatus = "revoked"
)

// APIKey represents an API key in the system
type APIKey struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	KeyType        APIKeyType      `json:"key_type"`
	KeyHash        string          `json:"key_hash"`
	Name           string          `json:"name"`
	CreatedAt      time.Time       `json:"created_at"`
	ExpiresAt      *time.Time      `json:"expires_at,omitempty"`
	LastUsedAt     *time.Time      `json:"last_used_at,omitempty"`
	Status         APIKeyStatus    `json:"status"`
	Permissions    json.RawMessage `json:"permissions"`
}

// APIKeyRepository defines the interface for API key data operations
type APIKeyRepository interface {
	Create(key *APIKey) error
	Update(key *APIKey) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*APIKey, error)
	GetByHash(hash string) (*APIKey, error)
	ListByOrganization(orgID uuid.UUID, offset, limit int) ([]*APIKey, error)
	UpdateLastUsed(id uuid.UUID) error
	RevokeKey(id uuid.UUID) error
}

// NewAPIKey creates a new APIKey instance with default values
func NewAPIKey(orgID uuid.UUID, name string, keyType APIKeyType, hash string) *APIKey {
	now := time.Now().UTC()
	return &APIKey{
		ID:             uuid.New(),
		OrganizationID: orgID,
		KeyType:        keyType,
		KeyHash:        hash,
		Name:           name,
		CreatedAt:      now,
		Status:         APIKeyStatusActive,
		Permissions:    json.RawMessage("{}"),
	}
}

// IsActive checks if the API key is active
func (k *APIKey) IsActive() bool {
	return k.Status == APIKeyStatusActive &&
		(k.ExpiresAt == nil || k.ExpiresAt.After(time.Now().UTC()))
}

// IsExpired checks if the API key has expired
func (k *APIKey) IsExpired() bool {
	return k.ExpiresAt != nil && k.ExpiresAt.Before(time.Now().UTC())
}

// SetExpiration sets the expiration time for the API key
func (k *APIKey) SetExpiration(duration time.Duration) {
	expiresAt := time.Now().UTC().Add(duration)
	k.ExpiresAt = &expiresAt
}

// UpdatePermissions updates the API key's permissions
func (k *APIKey) UpdatePermissions(permissions map[string]interface{}) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	k.Permissions = json.RawMessage(data)
	return nil
}

// Revoke marks the API key as revoked
func (k *APIKey) Revoke() {
	k.Status = APIKeyStatusRevoked
}
