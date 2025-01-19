package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/travism26/log-aggregator/internal/domain"
)

// APIKeyRepository implements domain.APIKeyRepository
type APIKeyRepository struct {
	db *sql.DB
}

// NewAPIKeyRepository creates a new APIKeyRepository instance
func NewAPIKeyRepository(db *sql.DB) domain.APIKeyRepository {
	return &APIKeyRepository{
		db: db,
	}
}

// Create stores a new API key
func (r *APIKeyRepository) Create(key *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (
			id, organization_id, key_type, key_hash, name,
			created_at, status, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query,
		key.ID,
		key.OrganizationID,
		key.KeyType,
		key.KeyHash,
		key.Name,
		key.CreatedAt,
		key.Status,
		key.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	return nil
}

// GetByHash retrieves an API key by its hash
func (r *APIKeyRepository) GetByHash(hash string) (*domain.APIKey, error) {
	query := `
		SELECT 
			id, organization_id, key_type, key_hash, name,
			created_at, last_used_at, status, expires_at
		FROM api_keys
		WHERE key_hash = $1`

	key := &domain.APIKey{}
	err := r.db.QueryRow(query, hash).Scan(
		&key.ID,
		&key.OrganizationID,
		&key.KeyType,
		&key.KeyHash,
		&key.Name,
		&key.CreatedAt,
		&key.LastUsedAt,
		&key.Status,
		&key.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return key, nil
}

// GetByID retrieves an API key by its ID
func (r *APIKeyRepository) GetByID(id uuid.UUID) (*domain.APIKey, error) {
	query := `
		SELECT 
			id, organization_id, key_type, key_hash, name,
			created_at, last_used_at, status, expires_at
		FROM api_keys
		WHERE id = $1`

	key := &domain.APIKey{}
	err := r.db.QueryRow(query, id).Scan(
		&key.ID,
		&key.OrganizationID,
		&key.KeyType,
		&key.KeyHash,
		&key.Name,
		&key.CreatedAt,
		&key.LastUsedAt,
		&key.Status,
		&key.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return key, nil
}

// ListByOrganization lists all API keys for an organization
func (r *APIKeyRepository) ListByOrganization(orgID uuid.UUID, limit, offset int) ([]*domain.APIKey, error) {
	query := `
		SELECT 
			id, organization_id, key_type, key_hash, name,
			created_at, last_used_at, status, expires_at
		FROM api_keys
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer rows.Close()

	var keys []*domain.APIKey
	for rows.Next() {
		key := &domain.APIKey{}
		err := rows.Scan(
			&key.ID,
			&key.OrganizationID,
			&key.KeyType,
			&key.KeyHash,
			&key.Name,
			&key.CreatedAt,
			&key.LastUsedAt,
			&key.Status,
			&key.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// UpdateLastUsed updates the last_used_at timestamp for an API key
func (r *APIKeyRepository) UpdateLastUsed(id uuid.UUID) error {
	query := `
		UPDATE api_keys
		SET last_used_at = $1
		WHERE id = $2`

	result, err := r.db.Exec(query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update last used timestamp: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// RevokeKey revokes an API key
// Delete permanently removes an API key
func (r *APIKeyRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// Update updates an existing API key
func (r *APIKeyRepository) Update(key *domain.APIKey) error {
	query := `
		UPDATE api_keys
		SET 
			name = $1,
			key_type = $2,
			status = $3,
			expires_at = $4,
			permissions = $5
		WHERE id = $6 AND organization_id = $7`

	result, err := r.db.Exec(query,
		key.Name,
		key.KeyType,
		key.Status,
		key.ExpiresAt,
		key.Permissions,
		key.ID,
		key.OrganizationID,
	)

	if err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

func (r *APIKeyRepository) RevokeKey(id uuid.UUID) error {
	query := `
		UPDATE api_keys
		SET status = $1
		WHERE id = $2`

	result, err := r.db.Exec(query, domain.APIKeyStatusRevoked, id)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}
