package exporter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/travism26/shared-monitoring-libs/types"
	"github.com/travism26/system-monitoring-agent/internal/config"
)

type Exporter interface {
	Export(data types.MetricPayload) error
	Close() error
}

func TestNewHTTPExporter(t *testing.T) {
	// Test case 1: with endpoint
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Endpoint: "http://example.com",
		},
	}
	exporter, err := NewHTTPExporter(cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create exporter: %v", err)
	}

	if !exporter.enabled {
		t.Error("Expected exporter to be enabled with endpoint")
	}
	if exporter.apiEndpoint != cfg.HTTP.Endpoint {
		t.Errorf("Expected endpoint %s, got %s", cfg.HTTP.Endpoint, exporter.apiEndpoint)
	}

	// Test case 2: without endpoint (disabled)
	disabledCfg := &config.Config{
		HTTP: config.HTTPConfig{
			Endpoint: "",
		},
	}
	disabledExporter, _ := NewHTTPExporter(disabledCfg, nil)
	if disabledExporter.enabled {
		t.Error("Expected exporter to be disabled without endpoint")
	}
}

func TestHTTPExporter_Export(t *testing.T) {
	// Create a test server
	var receivedData map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("🔵 Test server received a request")
		t.Log("📝 Request Method:", r.Method)
		t.Log("📝 Content-Type:", r.Header.Get("Content-Type"))

		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Decode and log the received data
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&receivedData); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		t.Log("📥 Received data:", receivedData)

		// Respond with success
		w.WriteHeader(http.StatusOK)
		t.Log("✅ Responding with HTTP 200 OK")
	}))
	defer server.Close()

	// Log the test server URL
	t.Log("🌐 Test server started at:", server.URL)

	// Create exporter with test server URL
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Endpoint: server.URL,
		},
	}
	exporter, err := NewHTTPExporter(cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create exporter: %v", err)
	}

	// Prepare test data
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
		Processes: struct {
			TotalCount       int                 `json:"total_count"`
			TotalCPUPercent  float64             `json:"total_cpu_percent"`
			TotalMemoryUsage uint64              `json:"total_memory_usage"`
			List             []types.ProcessInfo `json:"list"`
		}{
			TotalCount:       1,
			TotalCPUPercent:  10.5,
			TotalMemoryUsage: 1024,
			List:             []types.ProcessInfo{},
		},
		Metadata: struct {
			CollectionDuration string   `json:"collection_duration"`
			CollectorCount     int      `json:"collector_count"`
			Errors             []string `json:"errors,omitempty"`
		}{
			CollectionDuration: "100ms",
			CollectorCount:     5,
		},
	}
	t.Log("📤 Sending test data:", testData)

	// Convert struct to map for comparison
	testDataMap, _ := json.Marshal(testData)
	var testDataAsMap map[string]interface{}
	json.Unmarshal(testDataMap, &testDataAsMap)

	// Test export
	err = exporter.Export(testData)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	t.Log("✅ Export completed without errors")

	// Verify the exported data
	t.Log("🔍 Verifying received data matches sent data...")
	for key, expected := range testDataAsMap {
		received, ok := receivedData[key]
		if !ok {
			t.Errorf("❌ Missing key in exported data: %s", key)
			continue
		}

		t.Logf("Checking %s: expected=%v (%T), received=%v (%T)",
			key, expected, expected, received, received)

		if received != expected {
			t.Errorf("❌ For key %s: expected %v, got %v", key, expected, received)
		} else {
			t.Logf("✅ Key %s matches expected value", key)
		}
	}
}

func TestHTTPExporter_ExportError(t *testing.T) {
	// Test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Endpoint: server.URL,
		},
	}
	exporter, err := NewHTTPExporter(cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create exporter: %v", err)
	}

	// Create test data as MetricPayload
	testData := types.MetricPayload{
		Timestamp: "2025-02-07T06:09:25Z",
		TenantID:  "test-tenant",
		Metrics: map[string]interface{}{
			"cpu_usage":    0,
			"memory_usage": types.MemoryMetrics{Used: 0},
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

	// Test export with error response
	err = exporter.Export(testData)
	if err != nil { // Note: Our implementation logs errors but returns nil
		t.Errorf("Expected nil error on HTTP failure, got: %v", err)
	}
}

func TestHTTPExporter_Close(t *testing.T) {
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Endpoint: "http://example.com",
		},
	}
	exporter, err := NewHTTPExporter(cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create exporter: %v", err)
	}
	err = exporter.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
