package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
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
	testData := map[string]interface{}{
		"cpu":    50.5,
		"memory": float64(1024),
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
	var exportedData map[string]interface{}
	err = json.Unmarshal(content[:len(content)-1], &exportedData) // Remove trailing newline
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify the exported data
	if exportedData["cpu"] != testData["cpu"] {
		t.Errorf("Expected cpu value %v, got %v", testData["cpu"], exportedData["cpu"])
	}
	if exportedData["memory"] != testData["memory"] {
		t.Errorf("Expected memory value %v (type: %T), got %v (type: %T)",
			testData["memory"], testData["memory"],
			exportedData["memory"], exportedData["memory"])
	}
}
