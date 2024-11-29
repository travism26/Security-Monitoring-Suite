package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/travism26/shared-monitoring-libs/types"
)

func TestNewFileExporter(t *testing.T) {
	// Test the constructor
	path := "test.json"
	exporter := NewFileExporter(path)

	if exporter.outputFilePath != path {
		t.Errorf("Expected outputFilePath to be %s, got %s", path, exporter.outputFilePath)
	}
}

func TestFileExporter_Export(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "exporter_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after test

	// Create test file path
	testFile := filepath.Join(tmpDir, "test.json")

	// Initialize exporter
	exporter := NewFileExporter(testFile)

	// Test data
	testData := types.MetricPayload{
		Data: types.MetricData{
			Metrics: map[string]interface{}{
				"cpu_usage":    75.5,
				"memory_usage": 2048,
			},
		},
	}

	// Test export
	err = exporter.Export(testData)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Read the file content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Parse the JSON line
	var exportedData types.MetricPayload
	err = json.Unmarshal(content[:len(content)-1], &exportedData) // Remove trailing newline
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify the exported data
	if exportedData.Data.Metrics["cpu_usage"] != testData.Data.Metrics["cpu_usage"] {
		t.Errorf("Expected cpu value %v, got %v", testData.Data.Metrics["cpu_usage"], exportedData.Data.Metrics["cpu_usage"])
	}
	if exportedData.Data.Metrics["memory_usage"] != testData.Data.Metrics["memory_usage"] {
		t.Errorf("Expected memory value %v (type: %T), got %v (type: %T)",
			testData.Data.Metrics["memory_usage"], testData.Data.Metrics["memory_usage"],
			exportedData.Data.Metrics["memory_usage"], exportedData.Data.Metrics["memory_usage"])
	}
}
