package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/travism26/shared-monitoring-libs/types"
)

type HTTPExporter struct {
	apiEndpoint string
	client      *http.Client
	enabled     bool
	retryQueue  chan MetricBatch
	storage     MetricStorage
}

type MetricBatch struct {
	Data      types.MetricPayload // Changed to use types.MetricPayload instead of map[string]interface{}
	Timestamp time.Time
	Attempts  int
}

func NewHTTPExporter(endpoint string, storage MetricStorage) (*HTTPExporter, error) {
	enabled := endpoint != ""
	if !enabled {
		log.Println("HTTP exporter disabled: no endpoint configured")
	}

	exporter := &HTTPExporter{
		apiEndpoint: endpoint,
		enabled:     enabled,
		retryQueue:  make(chan MetricBatch, 1000),
		storage:     storage,
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
				DisableKeepAlives:  false,
			},
		},
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
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := h.client.Post(h.apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

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

func (h *HTTPExporter) Close() error {
	h.client.CloseIdleConnections()
	return nil
}
