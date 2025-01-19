package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/travism26/log-aggregator/internal/domain"

	_ "github.com/lib/pq"
)

const (
	defaultBatchSize = 1000
)

type LogRepository struct {
	db        *sql.DB
	batchSize int
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{
		db:        db,
		batchSize: defaultBatchSize,
	}
}

func (r *LogRepository) SetBatchSize(size int) {
	if size > 0 {
		r.batchSize = size
	}
}

func (r *LogRepository) Store(log *domain.Log) error {
	query := `
		INSERT INTO logs (
			id, organization_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		log.ID,
		log.OrganizationID,
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

	// Process in batches to avoid memory issues with large sets
	for i := 0; i < len(logs); i += r.batchSize {
		end := i + r.batchSize
		if end > len(logs) {
			end = len(logs)
		}
		if err := r.storeBatchChunk(logs[i:end]); err != nil {
			return fmt.Errorf("failed to store batch chunk %d-%d: %w", i, end, err)
		}
	}

	return nil
}

func (r *LogRepository) storeBatchChunk(logs []*domain.Log) error {
	// Create a transaction for batch insert
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if we return with error

	// Build the bulk insert query
	valueStrings := make([]string, 0, len(logs))
	valueArgs := make([]interface{}, 0, len(logs)*9)
	for i, log := range logs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9))
		valueArgs = append(valueArgs,
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
	}

	query := fmt.Sprintf(`
		INSERT INTO logs (
			id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		)
		VALUES %s`, strings.Join(valueStrings, ","))

	// Execute the bulk insert
	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *LogRepository) FindByID(orgID, id string) (*domain.Log, error) {
	query := `
		SELECT 
			id, organization_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM logs 
		WHERE organization_id = $1 AND id = $2
	`
	log := &domain.Log{}
	err := r.db.QueryRow(query, orgID, id).Scan(
		&log.ID,
		&log.OrganizationID,
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

func (r *LogRepository) List(orgID string, limit, offset int) ([]*domain.Log, error) {
	// Use CTE to optimize the pagination query
	query := `
		WITH recent_logs AS (
			SELECT 
				id, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs
			WHERE organization_id = $3
			ORDER BY timestamp DESC
		)
		SELECT 
			id, organization_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM recent_logs
		WHERE row_num > $2 AND row_num <= ($2 + $1)
	`
	rows, err := r.db.Query(query, limit, offset, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *LogRepository) ListByTimeRange(orgID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	// Use CTE with time range filter for better performance
	query := `
		WITH time_range_logs AS (
			SELECT 
				id, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
			ORDER BY timestamp DESC
		)
		SELECT 
			id, organization_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM time_range_logs
		WHERE row_num > $4 AND row_num <= ($4 + $3)
	`
	rows, err := r.db.Query(query, orgID, start, end, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by time range: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// CountByTimeRange returns the total number of logs within a time range
func (r *LogRepository) CountByTimeRange(orgID string, start, end time.Time) (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM logs 
		WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
	`
	var count int64
	err := r.db.QueryRow(query, orgID, start, end).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get log count: %w", err)
	}
	return count, nil
}

// Helper function to scan rows into log structs
func (r *LogRepository) scanLogs(rows *sql.Rows) ([]*domain.Log, error) {
	var logs []*domain.Log
	for rows.Next() {
		log := &domain.Log{}
		err := rows.Scan(
			&log.ID,
			&log.OrganizationID,
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

func (r *LogRepository) ListByHost(orgID string, host string, limit, offset int) ([]*domain.Log, error) {
	query := `
		WITH host_logs AS (
			SELECT 
				id, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE organization_id = $1 AND host = $2
			ORDER BY timestamp DESC
		)
		SELECT 
			id, organization_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM host_logs
		WHERE row_num > $4 AND row_num <= ($4 + $3)
	`
	rows, err := r.db.Query(query, orgID, host, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by host: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *LogRepository) ListByLevel(orgID string, level string, limit, offset int) ([]*domain.Log, error) {
	query := `
		WITH level_logs AS (
			SELECT 
				id, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE organization_id = $1 AND level = $2
			ORDER BY timestamp DESC
		)
		SELECT 
			id, organization_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM level_logs
		WHERE row_num > $4 AND row_num <= ($4 + $3)
	`
	rows, err := r.db.Query(query, orgID, level, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by level: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}
