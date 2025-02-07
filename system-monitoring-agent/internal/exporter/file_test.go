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
		Timestamp: "2025-02-07T06:09:25Z",
		TenantID:  "test-tenant",
		Metrics: map[string]interface{}{
			"cpu_usage":    75.5,
			"memory_usage": 2048,
		},
		Host: struct {
			OS        string `json:"os"`
			Arch      string `json:"arch"`
			Hostname  string `json:"hostname"`
			CPUCores  int    `json:"cpu_cores"`
			GoVersion string `json:"go_version"`
		}{
			OS:        "linux",
			Arch:      "amd64",
			Hostname:  "test-host",
			CPUCores:  4,
			GoVersion: "go1.21",
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
	if exportedData.Metrics["cpu_usage"] != testData.Metrics["cpu_usage"] {
		t.Errorf("Expected cpu value %v, got %v", testData.Metrics["cpu_usage"], exportedData.Metrics["cpu_usage"])
	}
	if exportedData.Metrics["memory_usage"] != testData.Metrics["memory_usage"] {
		t.Errorf("Expected memory value %v (type: %T), got %v (type: %T)",
			testData.Metrics["memory_usage"], testData.Metrics["memory_usage"],
			exportedData.Metrics["memory_usage"], exportedData.Metrics["memory_usage"])
	}
	if exportedData.TenantID != testData.TenantID {
		t.Errorf("Expected tenant ID %v, got %v", testData.TenantID, exportedData.TenantID)
	}
	if exportedData.Host.Hostname != testData.Host.Hostname {
		t.Errorf("Expected hostname %v, got %v", testData.Host.Hostname, exportedData.Host.Hostname)
	}
}
