package postgres

import (
	"database/sql"
	"fmt"

	"github.com/travism26/log-aggregator/internal/domain"

	_ "github.com/lib/pq"
)

type LogRepository struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{
		db: db,
	}
}

func (r *LogRepository) Store(log *domain.Log) error {
	query := `
		INSERT INTO logs (id, timestamp, host, message, level, metadata, process_count, total_cpu_percent, total_memory_usage)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	// Debug logging
	fmt.Printf("=== Debug Store Method ===\n")
	fmt.Printf("ID: %s\n", log.ID)
	fmt.Printf("Timestamp: %v\n", log.Timestamp)
	fmt.Printf("Host: %s\n", log.Host)
	fmt.Printf("Message: %s\n", log.Message)
	fmt.Printf("Level: %s\n", log.Level)
	fmt.Printf("MetadataStr (type: %T): %s\n", log.MetadataStr, log.MetadataStr)
	fmt.Printf("Raw Metadata (type: %T): %+v\n", log.Metadata, log.Metadata)
	fmt.Printf("=====================\n")

	_, err := r.db.Exec(query,
		log.ID,
		log.Timestamp,
		log.Host,
		log.Message,
		log.Level,
		log.MetadataStr,
		log.ProcessCount,
		log.TotalCPUPercent,
		log.TotalMemoryUsage,
	)

	if err != nil {
		fmt.Printf("=== Error Details ===\n")
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error Type: %T\n", err)
		fmt.Printf("==================\n")
		return fmt.Errorf("failed to store log: %w", err)
	}

	fmt.Printf("=== Success ===\n")
	fmt.Printf("Successfully stored log with ID: %s\n", log.ID)
	fmt.Printf("==============\n")

	return nil
}

func (r *LogRepository) FindByID(id string) (*domain.Log, error) {
	query := `
		SELECT id, timestamp, host, message, level, metadata
		FROM logs WHERE id = $1
	`
	log := &domain.Log{}
	err := r.db.QueryRow(query, id).Scan(
		&log.ID, &log.Timestamp, &log.Host,
		&log.Message, &log.Level, &log.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (r *LogRepository) List(limit, offset int) ([]*domain.Log, error) {
	query := `
		SELECT id, timestamp, host, message, level, metadata
		FROM logs ORDER BY timestamp DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.Log
	for rows.Next() {
		log := &domain.Log{}
		err := rows.Scan(
			&log.ID, &log.Timestamp, &log.Host,
			&log.Message, &log.Level, &log.Metadata,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *LogRepository) StoreLog(log *domain.Log) error {
	query := `
		INSERT INTO logs (timestamp, host, message, level, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		log.Timestamp,
		log.Host,
		log.Message,
		log.Level,
		log.MetadataStr,
	)

	if err != nil {
		return fmt.Errorf("failed to store log: %w", err)
	}

	return nil
}
