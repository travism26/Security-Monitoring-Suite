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
	fmt.Printf("[DEBUG] Repository: Storing log with ID: %s\n", log.ID)

	var query string
	var args []interface{}

	if log.OrganizationID == "" {
		// Query without organization_id for optional multi-tenancy
		query = `
			INSERT INTO logs (
				id, api_key, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
		args = []interface{}{
			log.ID,
			log.APIKey,
			log.Timestamp,
			log.Host,
			log.Message,
			log.Level,
			log.MetadataStr,
			log.ProcessCount,
			log.TotalCPUPercent,
			log.TotalMemoryUsage,
		}
	} else {
		// Query with organization_id when provided
		query = `
			INSERT INTO logs (
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		args = []interface{}{
			log.ID,
			log.APIKey,
			log.OrganizationID,
			log.Timestamp,
			log.Host,
			log.Message,
			log.Level,
			log.MetadataStr,
			log.ProcessCount,
			log.TotalCPUPercent,
			log.TotalMemoryUsage,
		}
	}

	fmt.Printf("[DEBUG] Repository: Executing query with parameters:\n")
	fmt.Printf("  - ID: %s\n", log.ID)
	fmt.Printf("  - APIKey: %s\n", log.APIKey)
	fmt.Printf("  - OrganizationID: %s\n", log.OrganizationID)
	fmt.Printf("  - Timestamp: %v\n", log.Timestamp)
	fmt.Printf("  - Host: %s\n", log.Host)
	fmt.Printf("  - Level: %s\n", log.Level)
	fmt.Printf("  - ProcessCount: %d\n", log.ProcessCount)
	fmt.Printf("  - TotalCPUPercent: %f\n", log.TotalCPUPercent)
	fmt.Printf("  - TotalMemoryUsage: %d\n", log.TotalMemoryUsage)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		fmt.Printf("[ERROR] Repository: Database error: %v\n", err)
		fmt.Printf("[ERROR] Repository: Failed query parameters: %+v\n", log)
		return fmt.Errorf("failed to store log: %w", err)
	}

	fmt.Printf("[DEBUG] Repository: Successfully stored log with ID: %s\n", log.ID)

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
	valueArgs := make([]interface{}, 0, len(logs)*11) // Increased to 11 to accommodate api_key
	for i, log := range logs {
		if log.OrganizationID == "" {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10))
			valueArgs = append(valueArgs,
				log.ID,
				log.APIKey,
				log.Timestamp,
				log.Host,
				log.Message,
				log.Level,
				log.MetadataStr,
				log.ProcessCount,
				log.TotalCPUPercent,
				log.TotalMemoryUsage,
			)
		} else {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*11+1, i*11+2, i*11+3, i*11+4, i*11+5, i*11+6, i*11+7, i*11+8, i*11+9, i*11+10, i*11+11))
			valueArgs = append(valueArgs,
				log.ID,
				log.APIKey,
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
		}
	}

	var query string
	if len(valueStrings) > 0 {
		if logs[0].OrganizationID == "" {
			query = fmt.Sprintf(`
				INSERT INTO logs (
					id, api_key, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage
				)
				VALUES %s`, strings.Join(valueStrings, ","))
		} else {
			query = fmt.Sprintf(`
				INSERT INTO logs (
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage
				)
				VALUES %s`, strings.Join(valueStrings, ","))
		}
	}

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
	var query string
	var args []interface{}

	if orgID == "" {
		query = `
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM logs 
			WHERE id = $1
		`
		args = []interface{}{id}
	} else {
		query = `
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM logs 
			WHERE organization_id = $1 AND id = $2
		`
		args = []interface{}{orgID, id}
	}
	log := &domain.Log{}
	err := r.db.QueryRow(query, args...).Scan(
		&log.ID,
		&log.APIKey,
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
	var query string
	var args []interface{}

	if orgID == "" {
		query = `
			WITH recent_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM recent_logs
			WHERE row_num > $1 AND row_num <= ($1 + $2)
		`
		args = []interface{}{offset, limit}
	} else {
		query = `
			WITH recent_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs
				WHERE organization_id = $3
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM recent_logs
			WHERE row_num > $1 AND row_num <= ($1 + $2)
		`
		args = []interface{}{offset, limit, orgID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *LogRepository) ListByTimeRange(orgID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	// Use CTE with time range filter for better performance
	var query string
	var args []interface{}

	if orgID == "" {
		query = `
			WITH time_range_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs 
				WHERE timestamp >= $1 AND timestamp <= $2
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM time_range_logs
			WHERE row_num > $3 AND row_num <= ($3 + $4)
		`
		args = []interface{}{start, end, offset, limit}
	} else {
		query = `
			WITH time_range_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs 
				WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM time_range_logs
			WHERE row_num > $4 AND row_num <= ($4 + $5)
		`
		args = []interface{}{orgID, start, end, offset, limit}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by time range: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// CountByTimeRange returns the total number of logs within a time range
func (r *LogRepository) CountByTimeRange(orgID string, start, end time.Time) (int64, error) {
	var query string
	var args []interface{}

	if orgID == "" {
		query = `
			SELECT COUNT(*) 
			FROM logs 
			WHERE timestamp >= $1 AND timestamp <= $2
		`
		args = []interface{}{start, end}
	} else {
		query = `
			SELECT COUNT(*) 
			FROM logs 
			WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
		`
		args = []interface{}{orgID, start, end}
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
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
			&log.APIKey,
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
	var query string
	var args []interface{}

	if orgID == "" {
		query = `
			WITH host_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs 
				WHERE host = $1
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM host_logs
			WHERE row_num > $2 AND row_num <= ($2 + $3)
		`
		args = []interface{}{host, offset, limit}
	} else {
		query = `
			WITH host_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs 
				WHERE organization_id = $1 AND host = $2
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM host_logs
			WHERE row_num > $3 AND row_num <= ($3 + $4)
		`
		args = []interface{}{orgID, host, offset, limit}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by host: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *LogRepository) ListByLevel(orgID string, level string, limit, offset int) ([]*domain.Log, error) {
	var query string
	var args []interface{}

	if orgID == "" {
		query = `
			WITH level_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs 
				WHERE level = $1
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM level_logs
			WHERE row_num > $2 AND row_num <= ($2 + $3)
		`
		args = []interface{}{level, offset, limit}
	} else {
		query = `
			WITH level_logs AS (
				SELECT 
					id, api_key, organization_id, timestamp, host, message, level, metadata,
					process_count, total_cpu_percent, total_memory_usage,
					ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
				FROM logs 
				WHERE organization_id = $1 AND level = $2
				ORDER BY timestamp DESC
			)
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage
			FROM level_logs
			WHERE row_num > $3 AND row_num <= ($3 + $4)
		`
		args = []interface{}{orgID, level, offset, limit}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by level: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// ListByAPIKey retrieves logs for a specific API key with pagination
func (r *LogRepository) ListByAPIKey(apiKey string, limit, offset int) ([]*domain.Log, error) {
	query := `
		WITH api_key_logs AS (
			SELECT 
				id, api_key, organization_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE api_key = $1
			ORDER BY timestamp DESC
		)
		SELECT 
			id, api_key, organization_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM api_key_logs
		WHERE row_num > $2 AND row_num <= ($2 + $3)
	`
	args := []interface{}{apiKey, offset, limit}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by API key: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// CountByAPIKey returns the total number of logs for a specific API key
func (r *LogRepository) CountByAPIKey(apiKey string) (int64, error) {
	query := `SELECT COUNT(*) FROM logs WHERE api_key = $1`
	var count int64
	err := r.db.QueryRow(query, apiKey).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count logs by API key: %w", err)
	}
	return count, nil
}
