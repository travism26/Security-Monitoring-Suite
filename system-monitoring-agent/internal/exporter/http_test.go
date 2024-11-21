package exporter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPExporter(t *testing.T) {
	// Test case 1: with endpoint
	endpoint := "http://example.com"
	exporter := NewHTTPExporter(endpoint)

	if !exporter.enabled {
		t.Error("Expected exporter to be enabled with endpoint")
	}
	if exporter.apiEndpoint != endpoint {
		t.Errorf("Expected endpoint %s, got %s", endpoint, exporter.apiEndpoint)
	}

	// Test case 2: without endpoint (disabled)
	disabledExporter := NewHTTPExporter("")
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
	exporter := NewHTTPExporter(server.URL)

	// Prepare test data
	testData := map[string]interface{}{
		"cpu":    75.5,
		"memory": float64(2048),
	}
	t.Log("📤 Sending test data:", testData)

	// Test export
	err := exporter.Export(testData)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	t.Log("✅ Export completed without errors")

	// Verify the exported data
	t.Log("🔍 Verifying received data matches sent data...")
	for key, expected := range testData {
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

	exporter := NewHTTPExporter(server.URL)
	testData := map[string]interface{}{"test": "data"}

	// Test export with error response
	err := exporter.Export(testData)
	if err != nil { // Note: Our implementation logs errors but returns nil
		t.Errorf("Expected nil error on HTTP failure, got: %v", err)
	}
}

func TestHTTPExporter_Close(t *testing.T) {
	exporter := NewHTTPExporter("http://example.com")
	err := exporter.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
