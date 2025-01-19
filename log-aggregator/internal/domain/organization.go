package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OrganizationStatus represents the possible states of an organization
type OrganizationStatus string

const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusInactive  OrganizationStatus = "inactive"
	OrganizationStatusSuspended OrganizationStatus = "suspended"
)

// Organization represents a tenant in the multi-tenant system
type Organization struct {
	ID        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Status    OrganizationStatus `json:"status"`
	Settings  json.RawMessage    `json:"settings"`
}

// OrganizationRepository defines the interface for organization data operations
type OrganizationRepository interface {
	Create(org *Organization) error
	Update(org *Organization) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*Organization, error)
	List(offset, limit int) ([]*Organization, error)
	GetByName(name string) (*Organization, error)
	UpdateStatus(id uuid.UUID, status OrganizationStatus) error
}

// NewOrganization creates a new Organization instance with default values
func NewOrganization(name string) *Organization {
	now := time.Now().UTC()
	return &Organization{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    OrganizationStatusActive,
		Settings:  json.RawMessage("{}"),
	}
}

// IsActive checks if the organization is in active status
func (o *Organization) IsActive() bool {
	return o.Status == OrganizationStatusActive
}

// IsSuspended checks if the organization is suspended
func (o *Organization) IsSuspended() bool {
	return o.Status == OrganizationStatusSuspended
}

// UpdateSettings updates the organization's settings
func (o *Organization) UpdateSettings(settings map[string]interface{}) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	o.Settings = json.RawMessage(data)
	o.UpdatedAt = time.Now().UTC()
	return nil
}
