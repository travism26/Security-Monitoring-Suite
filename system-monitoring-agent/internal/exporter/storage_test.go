package exporter

import (
	"os"
	"testing"
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

	// Convert to MetricBatch before passing to storage
	metricBatch := MetricBatch{
		// convert testData fields to MetricBatch format
	}

	// Test Store
	if err := storage.Store(metricBatch); err != nil {
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
	if err := storage.Remove(metricBatch); err != nil {
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
