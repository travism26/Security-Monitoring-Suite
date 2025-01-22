// Package apikey handles secure API key management and validation
package apikey

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	ErrKeyExpired      = errors.New("API key has expired")
	ErrKeyInvalid      = errors.New("invalid API key")
	ErrKeyRotationFail = errors.New("failed to rotate API key")
	ErrEncryptionFail  = errors.New("failed to encrypt API key")
)

// KeyStatus represents the current state of an API key
type KeyStatus struct {
	IsValid         bool
	ExpiresAt       time.Time
	LastUsed        time.Time
	RotatedAt       time.Time
	ValidationError error
}

// Manager handles API key operations and lifecycle
type Manager struct {
	mu sync.RWMutex

	// Current active key
	currentKey string

	// Encrypted backup key for rotation
	backupKey []byte

	// Key metadata
	status KeyStatus

	// Encryption key for securing stored keys
	encryptionKey []byte

	// Configuration
	config Config
}

// Config holds API key manager configuration
type Config struct {
	// Validation endpoint for checking key status
	ValidationEndpoint string

	// How often to check key validity
	ValidationInterval time.Duration

	// Maximum key age before rotation is recommended
	MaxKeyAge time.Duration

	// Whether to encrypt stored keys
	EncryptKeys bool
}

// NewManager creates a new API key manager
func NewManager(initialKey string, cfg Config) (*Manager, error) {
	if initialKey == "" {
		return nil, ErrKeyInvalid
	}

	// Generate encryption key if encryption is enabled
	var encKey []byte
	if cfg.EncryptKeys {
		encKey = make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, encKey); err != nil {
			return nil, fmt.Errorf("failed to generate encryption key: %w", err)
		}
	}

	m := &Manager{
		currentKey: initialKey,
		status: KeyStatus{
			IsValid:   true, // Assume initially valid
			ExpiresAt: time.Now().Add(cfg.MaxKeyAge),
			LastUsed:  time.Now(),
			RotatedAt: time.Now(),
		},
		encryptionKey: encKey,
		config:        cfg,
	}

	// Start validation routine if interval is set
	if cfg.ValidationInterval > 0 {
		go m.startValidationRoutine()
	}

	return m, nil
}

// GetKey returns the current API key
func (m *Manager) GetKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentKey
}

// GetStatus returns the current key status
func (m *Manager) GetStatus() KeyStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

// UpdateKey sets a new API key
func (m *Manager) UpdateKey(newKey string) error {
	if newKey == "" {
		return ErrKeyInvalid
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Encrypt and store current key as backup before updating
	if m.config.EncryptKeys {
		encrypted, err := m.encryptKey(m.currentKey)
		if err != nil {
			return fmt.Errorf("failed to backup current key: %w", err)
		}
		m.backupKey = encrypted
	}

	m.currentKey = newKey
	m.status.RotatedAt = time.Now()
	m.status.LastUsed = time.Now()
	m.status.ExpiresAt = time.Now().Add(m.config.MaxKeyAge)
	m.status.IsValid = true
	m.status.ValidationError = nil

	return nil
}

// RotateKey generates and updates to a new API key
func (m *Manager) RotateKey() error {
	// Implementation would integrate with your key management system
	// to generate and validate a new key
	return ErrKeyRotationFail
}

// ValidateKey checks if the current key is valid
func (m *Manager) ValidateKey() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check expiration
	if time.Now().After(m.status.ExpiresAt) {
		m.status.IsValid = false
		m.status.ValidationError = ErrKeyExpired
		return ErrKeyExpired
	}

	// Implementation would call validation endpoint
	// For now, just update last used timestamp
	m.status.LastUsed = time.Now()
	return nil
}

// startValidationRoutine periodically validates the API key
func (m *Manager) startValidationRoutine() {
	ticker := time.NewTicker(m.config.ValidationInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.ValidateKey(); err != nil {
			// Log validation error but continue running
			// The error is stored in status.ValidationError
			continue
		}
	}
}

// encryptKey encrypts an API key for secure storage
func (m *Manager) encryptKey(key string) ([]byte, error) {
	if !m.config.EncryptKeys {
		return []byte(key), nil
	}

	block, err := aes.NewCipher(m.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create cipher", ErrEncryptionFail)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create GCM", ErrEncryptionFail)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("%w: failed to generate nonce", ErrEncryptionFail)
	}

	return gcm.Seal(nonce, nonce, []byte(key), nil), nil
}

// decryptKey decrypts a stored API key
func (m *Manager) decryptKey(encrypted []byte) (string, error) {
	if !m.config.EncryptKeys {
		return string(encrypted), nil
	}

	block, err := aes.NewCipher(m.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("%w: failed to create cipher", ErrEncryptionFail)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("%w: failed to create GCM", ErrEncryptionFail)
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", fmt.Errorf("%w: invalid encrypted key", ErrEncryptionFail)
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("%w: failed to decrypt", ErrEncryptionFail)
	}

	return string(plaintext), nil
}
