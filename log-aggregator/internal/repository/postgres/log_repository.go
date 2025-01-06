package postgres

import (
	"database/sql"
	"fmt"
	"time"

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
		INSERT INTO logs (
			id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

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
		return fmt.Errorf("failed to store log: %w", err)
	}

	return nil
}

func (r *LogRepository) StoreBatch(logs []*domain.Log) error {
	if len(logs) == 0 {
		return nil
	}

	// Create a transaction for batch insert
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if we return with error

	// Prepare the statement for batch insert
	stmt, err := tx.Prepare(`
		INSERT INTO logs (
			id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute batch insert
	for _, log := range logs {
		_, err = stmt.Exec(
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
			return fmt.Errorf("failed to store log %s in batch: %w", log.ID, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *LogRepository) FindByID(id string) (*domain.Log, error) {
	query := `
		SELECT 
			id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM logs 
		WHERE id = $1
	`
	log := &domain.Log{}
	err := r.db.QueryRow(query, id).Scan(
		&log.ID,
		&log.Timestamp,
		&log.Host,
		&log.Message,
		&log.Level,
		&log.MetadataStr,
		&log.ProcessCount,
		&log.TotalCPUPercent,
		&log.TotalMemoryUsage,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("log not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find log: %w", err)
	}
	return log, nil
}

func (r *LogRepository) List(limit, offset int) ([]*domain.Log, error) {
	query := `
		SELECT 
			id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM logs 
		ORDER BY timestamp DESC 
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *LogRepository) ListByTimeRange(start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	query := `
		SELECT 
			id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM logs 
		WHERE timestamp BETWEEN $1 AND $2
		ORDER BY timestamp DESC 
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(query, start, end, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by time range: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// Helper function to scan rows into log structs
func (r *LogRepository) scanLogs(rows *sql.Rows) ([]*domain.Log, error) {
	var logs []*domain.Log
	for rows.Next() {
		log := &domain.Log{}
		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.Host,
			&log.Message,
			&log.Level,
			&log.MetadataStr,
			&log.ProcessCount,
			&log.TotalCPUPercent,
			&log.TotalMemoryUsage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating log rows: %w", err)
	}
	return logs, nil
}
