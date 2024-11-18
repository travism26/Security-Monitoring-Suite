package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPExporter struct {
	apiEndpoint string
	client      *http.Client
}

func NewHTTPExporter(endpoint string) *HTTPExporter {
	return &HTTPExporter{
		apiEndpoint: endpoint,
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

func (h *HTTPExporter) Export(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := h.client.Post(h.apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned error status: %d", resp.StatusCode)
	}

	return nil
}

func (h *HTTPExporter) Close() error {
	// Clean up any resources if needed
	return nil
}
