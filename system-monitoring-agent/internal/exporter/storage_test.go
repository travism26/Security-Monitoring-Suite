package exporter

import (
	"os"
	"testing"
	"time"
)

func TestSQLiteStorage(t *testing.T) {
	// Create temporary directory for test database
	tmpDir, err := os.MkdirTemp("", "metric_storage_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize storage
	storage, err := NewMetricStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Test storing and loading metrics
	testBatch := MetricBatch{
		Data: map[string]interface{}{
			"cpu":    75.5,
			"memory": 1024,
		},
		Timestamp: time.Now(),
		Attempts:  0,
	}

	// Test Store
	if err := storage.Store(testBatch); err != nil {
		t.Errorf("Failed to store metric: %v", err)
	}

	// Test LoadUnsent
	batches, err := storage.LoadUnsent()
	if err != nil {
		t.Errorf("Failed to load metrics: %v", err)
	}
	if len(batches) != 1 {
		t.Errorf("Expected 1 batch, got %d", len(batches))
	}

	// Test Remove
	if err := storage.Remove(testBatch); err != nil {
		t.Errorf("Failed to remove metric: %v", err)
	}

	// Verify removal
	batches, err = storage.LoadUnsent()
	if err != nil {
		t.Errorf("Failed to load metrics after removal: %v", err)
	}
	if len(batches) != 0 {
		t.Errorf("Expected 0 batches after removal, got %d", len(batches))
	}
}
