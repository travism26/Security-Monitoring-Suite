package exporter

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
	"github.com/travism26/system-monitoring-agent/internal/config"
)

type HTTPExporter struct {
	apiEndpoint string
	client      *http.Client
	enabled     bool
	retryQueue  chan MetricBatch
	storage     MetricStorage
	config      *config.Config
	headers     map[string]string
}

type MetricBatch struct {
	Data      types.MetricPayload // Changed to use types.MetricPayload instead of map[string]interface{}
	Timestamp time.Time
	Attempts  int
}

func NewHTTPExporter(cfg *config.Config, storage MetricStorage) (*HTTPExporter, error) {
	enabled := cfg.Tenant.Endpoints.Metrics != ""
	if !enabled {
		log.Println("HTTP exporter disabled: no endpoint configured")
	}

	// Initialize headers
	headers := map[string]string{
		"Content-Type":         "application/json",
		"X-Tenant-Environment": cfg.Tenant.Environment,
		"X-Tenant-Type":        cfg.Tenant.Type,
	}

	// Add optional tenant headers if provided
	if cfg.Tenant.ID != "" {
		headers[cfg.HTTP.Headers.TenantID] = cfg.Tenant.ID
	}
	if cfg.Tenant.APIKey != "" {
		headers[cfg.HTTP.Headers.APIKey] = cfg.Tenant.APIKey
	}

	exporter := &HTTPExporter{
		apiEndpoint: cfg.Tenant.Endpoints.Metrics,
		enabled:     enabled,
		retryQueue:  make(chan MetricBatch, 1000),
		storage:     storage,
		config:      cfg,
		headers:     headers,
		client:      createHTTPClient(cfg),
	}

	go exporter.retryWorker()
	return exporter, nil
}

// Changed interface to use types.MetricPayload
// before: func (h *HTTPExporter) Export(data map[string]interface{}) error {
func (h *HTTPExporter) Export(data types.MetricPayload) error {
	if !h.enabled {
		return nil
	}

	batch := MetricBatch{
		Data:      data,
		Timestamp: time.Now(),
		Attempts:  0,
	}

	if err := h.sendBatch(batch); err != nil {
		if h.storage != nil {
			if err := h.storage.Store(batch); err != nil {
				log.Printf("Failed to store metrics: %v", err)
			}
		}
		select {
		case h.retryQueue <- batch:
		default:
			log.Println("Retry queue full, metric will be dropped")
		}
		return err
	}

	return nil
}

func (h *HTTPExporter) sendBatch(batch MetricBatch) error {
	jsonData, err := json.Marshal(batch.Data)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal metrics data: %v", err)
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	log.Printf("[DEBUG] Preparing to send metrics to endpoint: %s (payload size: %d bytes)", h.apiEndpoint, len(jsonData))

	req, err := http.NewRequest("POST", h.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range h.headers {
		req.Header.Set(key, value)
	}
	log.Printf("[DEBUG] Request headers: %v", req.Header)

	// Do sends an HTTP request and returns an HTTP response
	log.Printf("[DEBUG] Sending HTTP request to %s", h.apiEndpoint)
	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Failed to send metrics: %v", err)
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for error cases
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[ERROR] Server returned non-success status %d. Response body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	log.Printf("[DEBUG] Successfully sent metrics batch. Status: %d", resp.StatusCode)

	if h.storage != nil {
		h.storage.Remove(batch)
	}
	return nil
}

func (h *HTTPExporter) retryWorker() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case batch := <-h.retryQueue:
			if batch.Attempts >= 5 {
				log.Printf("Dropping metric batch after %d attempts", batch.Attempts)
				continue
			}

			batch.Attempts++
			if err := h.sendBatch(batch); err != nil {
				time.Sleep(time.Second * time.Duration(1<<batch.Attempts))
				h.retryQueue <- batch
			}

		case <-ticker.C:
			if h.storage != nil {
				batches, err := h.storage.LoadUnsent()
				if err != nil {
					log.Printf("Failed to load stored metrics: %v", err)
					continue
				}

				for _, batch := range batches {
					h.retryQueue <- batch
				}
			}
		}
	}
}

// createHTTPClient creates an HTTP client with TLS configuration
func createHTTPClient(cfg *config.Config) *http.Client {
	log.Printf("[DEBUG] Creating HTTP client with timeout: %d seconds", cfg.HTTP.Timeout)

	// Create TLS config
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: !cfg.Security.ValidateSSL,
	}

	log.Printf("[DEBUG] TLS Configuration - MinVersion: TLS1.2, ValidateSSL: %v", cfg.Security.ValidateSSL)

	// If TLS cert/key are configured, load them
	if cfg.Security.TLS.CertFile != "" && cfg.Security.TLS.KeyFile != "" {
		log.Printf("[DEBUG] Loading TLS certificate from: %s and key from: %s",
			cfg.Security.TLS.CertFile, cfg.Security.TLS.KeyFile)

		cert, err := tls.LoadX509KeyPair(
			filepath.Clean(cfg.Security.TLS.CertFile),
			filepath.Clean(cfg.Security.TLS.KeyFile),
		)
		if err != nil {
			log.Printf("[ERROR] Failed to load TLS cert/key: %v", err)
		} else {
			log.Printf("[DEBUG] Successfully loaded TLS certificate")
			// Add the certificate to the config
			tlsConfig.Certificates = []tls.Certificate{cert}

			// Create cert pool and add our certificate
			certPool := x509.NewCertPool()
			certData, err := ioutil.ReadFile(filepath.Clean(cfg.Security.TLS.CertFile))
			if err != nil {
				log.Printf("Warning: Failed to read certificate file: %v", err)
			} else {
				if ok := certPool.AppendCertsFromPEM(certData); !ok {
					log.Printf("Warning: Failed to append certificate to pool")
				} else {
					tlsConfig.RootCAs = certPool
				}
			}
		}
	}

	return &http.Client{
		Timeout: time.Duration(cfg.HTTP.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    tlsConfig,
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  false,
		},
	}
}

func (h *HTTPExporter) Close() error {
	h.client.CloseIdleConnections()
	return nil
}
