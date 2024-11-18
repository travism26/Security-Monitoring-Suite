package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPExporter struct {
	apiEndpoint string
	client      *http.Client
	enabled     bool
}

func NewHTTPExporter(endpoint string) *HTTPExporter {
	enabled := endpoint != ""
	if !enabled {
		log.Println("HTTP exporter disabled: no endpoint configured")
	}

	return &HTTPExporter{
		apiEndpoint: endpoint,
		enabled:     enabled,
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
				DisableKeepAlives:  false,
				// MaxRetries:         3,
				// RetryWaitMin:       1 * time.Second,
				// RetryWaitMax:       5 * time.Second,
			},
		},
	}
}

func (h *HTTPExporter) Export(data map[string]interface{}) error {
	if !h.enabled {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := h.client.Post(h.apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Warning: Failed to send metrics to HTTP endpoint: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Warning: HTTP endpoint returned status %d", resp.StatusCode)
		return nil
	}

	return nil
}

func (h *HTTPExporter) Close() error {
	h.client.CloseIdleConnections()
	return nil
}
