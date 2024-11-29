package exporter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/travism26/shared-monitoring-libs/types"
)

type MetricStorage interface {
	Store(batch MetricBatch) error
	Remove(batch MetricBatch) error
	LoadUnsent() ([]MetricBatch, error)
	Close() error
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewMetricStorage(storageDir string) (*SQLiteStorage, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	dbPath := fmt.Sprintf("%s/metrics.db", storageDir)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create metrics table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			attempts INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics table: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Store(batch MetricBatch) error {
	data, err := json.Marshal(batch.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal metric data: %w", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO metrics (data, timestamp, attempts) VALUES (?, ?, ?)",
		string(data),
		batch.Timestamp,
		batch.Attempts,
	)
	if err != nil {
		return fmt.Errorf("failed to store metric: %w", err)
	}

	return nil
}

func (s *SQLiteStorage) Remove(batch MetricBatch) error {
	data, err := json.Marshal(batch.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal metric data: %w", err)
	}

	result, err := s.db.Exec(
		"DELETE FROM metrics WHERE data = ? AND timestamp = ?",
		string(data),
		batch.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to remove metric: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("no matching metric found to remove")
	}

	return nil
}

func (s *SQLiteStorage) LoadUnsent() ([]MetricBatch, error) {
	rows, err := s.db.Query(`
		SELECT data, timestamp, attempts 
		FROM metrics 
		WHERE attempts < 5 
		ORDER BY timestamp ASC
		LIMIT 100
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var batches []MetricBatch
	for rows.Next() {
		var (
			dataStr   string
			timestamp time.Time
			attempts  int
		)

		if err := rows.Scan(&dataStr, &timestamp, &attempts); err != nil {
			return nil, fmt.Errorf("failed to scan metric row: %w", err)
		}

		var data types.MetricPayload
		if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metric data: %w", err)
		}

		batches = append(batches, MetricBatch{
			Data:      data,
			Timestamp: timestamp,
			Attempts:  attempts,
		})
	}

	return batches, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
