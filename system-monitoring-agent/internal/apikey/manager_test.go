package apikey

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name       string
		initialKey string
		cfg        Config
		wantErr    bool
	}{
		{
			name:       "Valid initialization",
			initialKey: "valid-key-12345",
			cfg: Config{
				ValidationInterval: time.Minute,
				MaxKeyAge:          24 * time.Hour,
				EncryptKeys:        true,
			},
			wantErr: false,
		},
		{
			name:       "Empty key",
			initialKey: "",
			cfg: Config{
				ValidationInterval: time.Minute,
				MaxKeyAge:          24 * time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.initialKey, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && manager == nil {
				t.Error("NewManager() returned nil manager without error")
			}
		})
	}
}

func TestManager_UpdateKey(t *testing.T) {
	cfg := Config{
		ValidationInterval: time.Minute,
		MaxKeyAge:          24 * time.Hour,
		EncryptKeys:        true,
	}
	manager, _ := NewManager("initial-key", cfg)

	tests := []struct {
		name    string
		newKey  string
		wantErr bool
	}{
		{
			name:    "Valid key update",
			newKey:  "new-valid-key-12345",
			wantErr: false,
		},
		{
			name:    "Empty key update",
			newKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.UpdateKey(tt.newKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && manager.GetKey() != tt.newKey {
				t.Errorf("UpdateKey() failed to update key, got = %v, want %v", manager.GetKey(), tt.newKey)
			}
		})
	}
}

func TestManager_ValidateKey(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func() *Manager
		wantErr    bool
		wantStatus bool
	}{
		{
			name: "Valid key",
			setupFunc: func() *Manager {
				m, _ := NewManager("valid-key", Config{
					MaxKeyAge: 24 * time.Hour,
				})
				return m
			},
			wantErr:    false,
			wantStatus: true,
		},
		{
			name: "Expired key",
			setupFunc: func() *Manager {
				m, _ := NewManager("expired-key", Config{
					MaxKeyAge: -1 * time.Hour, // Force expiration
				})
				return m
			},
			wantErr:    true,
			wantStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.setupFunc()
			err := manager.ValidateKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			status := manager.GetStatus()
			if status.IsValid != tt.wantStatus {
				t.Errorf("ValidateKey() key status = %v, want %v", status.IsValid, tt.wantStatus)
			}
		})
	}
}

func TestManager_Encryption(t *testing.T) {
	cfg := Config{
		EncryptKeys: true,
		MaxKeyAge:   24 * time.Hour,
	}
	manager, err := NewManager("test-key", cfg)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test key encryption
	encrypted, err := manager.encryptKey("secret-key")
	if err != nil {
		t.Errorf("encryptKey() error = %v", err)
		return
	}

	// Test key decryption
	decrypted, err := manager.decryptKey(encrypted)
	if err != nil {
		t.Errorf("decryptKey() error = %v", err)
		return
	}

	if decrypted != "secret-key" {
		t.Errorf("Key encryption/decryption failed, got = %v, want = %v", decrypted, "secret-key")
	}
}
