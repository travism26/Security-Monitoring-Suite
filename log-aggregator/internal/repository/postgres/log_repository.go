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

	query := `
		INSERT INTO logs (
			id, api_key, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage, organization_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	args := []interface{}{
		log.ID,
		log.APIKey,
		log.UserID,
		log.Timestamp,
		log.Host,
		log.Message,
		log.Level,
		log.MetadataStr,
		log.ProcessCount,
		log.TotalCPUPercent,
		log.TotalMemoryUsage,
		log.OrganizationID, // Always include organization_id, can be empty string
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
	valueArgs := make([]interface{}, 0, len(logs)*12) // 12 parameters per log entry
	for i, log := range logs {
		// Calculate the base position for this log's parameters
		base := i * 12
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8, base+9, base+10, base+11, base+12))
		valueArgs = append(valueArgs,
			log.ID,
			log.APIKey,
			log.UserID,
			log.Timestamp,
			log.Host,
			log.Message,
			log.Level,
			log.MetadataStr,
			log.ProcessCount,
			log.TotalCPUPercent,
			log.TotalMemoryUsage,
			log.OrganizationID, // Always include organization_id, can be empty string
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO logs (
			id, api_key, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage, organization_id
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

func (r *LogRepository) FindByID(userID, id string) (*domain.Log, error) {
	query := `
		SELECT 
			id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM logs 
		WHERE id = $1 AND user_id = $2
	`
	args := []interface{}{id, userID}
	log := &domain.Log{}
	err := r.db.QueryRow(query, args...).Scan(
		&log.ID,
		&log.APIKey,
		&log.OrganizationID,
		&log.UserID,
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

func (r *LogRepository) List(userID string, limit, offset int) ([]*domain.Log, error) {
	// Use CTE to optimize the pagination query
	query := `
		WITH recent_logs AS (
			SELECT 
				id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs
			WHERE user_id = $3
			ORDER BY timestamp DESC
		)
		SELECT 
			id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM recent_logs
		WHERE row_num > $1 AND row_num <= ($1 + $2)
	`
	args := []interface{}{offset, limit, userID}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *LogRepository) ListByTimeRange(userID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	// Use CTE with time range filter for better performance
	query := `
		WITH time_range_logs AS (
			SELECT 
				id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
			ORDER BY timestamp DESC
		)
		SELECT 
			id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM time_range_logs
		WHERE row_num > $4 AND row_num <= ($4 + $5)
	`
	args := []interface{}{userID, start, end, offset, limit}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by time range: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// CountByTimeRange returns the total number of logs within a time range
func (r *LogRepository) CountByTimeRange(userID string, start, end time.Time) (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM logs 
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
	`
	args := []interface{}{userID, start, end}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get log count: %w", err)
	}
	return count, nil
}

// ListByUserID retrieves logs for a specific user
func (r *LogRepository) ListByUserID(userID string, limit, offset int) ([]*domain.Log, error) {
	query := `
		WITH user_logs AS (
			SELECT 
				id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE user_id = $1
			ORDER BY timestamp DESC
		)
		SELECT 
			id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM user_logs
		WHERE row_num > $2 AND row_num <= ($2 + $3)
	`
	args := []interface{}{userID, offset, limit}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by user ID: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

// CountByUserID returns the total number of logs for a specific user
func (r *LogRepository) CountByUserID(userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM logs WHERE user_id = $1`
	args := []interface{}{userID}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count logs by user ID: %w", err)
	}
	return count, nil
}

// ListByUserIDAndTimeRange retrieves logs for a specific user within a time range
func (r *LogRepository) ListByUserIDAndTimeRange(userID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	query := `
		WITH user_time_logs AS (
			SELECT 
				id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE user_id = $1 
				AND timestamp >= $2 
				AND timestamp <= $3
			ORDER BY timestamp DESC
		)
		SELECT 
			id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM user_time_logs
		WHERE row_num > $4 AND row_num <= ($4 + $5)
	`
	args := []interface{}{userID, start, end, offset, limit}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by user ID and time range: %w", err)
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
			&log.APIKey,
			&log.OrganizationID,
			&log.UserID,
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

func (r *LogRepository) ListByHost(userID string, host string, limit, offset int) ([]*domain.Log, error) {
	query := `
		WITH host_logs AS (
			SELECT 
				id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE user_id = $1 AND host = $2
			ORDER BY timestamp DESC
		)
		SELECT 
			id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM host_logs
		WHERE row_num > $3 AND row_num <= ($3 + $4)
	`
	args := []interface{}{userID, host, offset, limit}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logs by host: %w", err)
	}
	defer rows.Close()

	return r.scanLogs(rows)
}

func (r *LogRepository) ListByLevel(userID string, level string, limit, offset int) ([]*domain.Log, error) {
	query := `
		WITH level_logs AS (
			SELECT 
				id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
				process_count, total_cpu_percent, total_memory_usage,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as row_num
			FROM logs 
			WHERE user_id = $1 AND level = $2
			ORDER BY timestamp DESC
		)
		SELECT 
			id, api_key, organization_id, user_id, timestamp, host, message, level, metadata,
			process_count, total_cpu_percent, total_memory_usage
		FROM level_logs
		WHERE row_num > $3 AND row_num <= ($3 + $4)
	`
	args := []interface{}{userID, level, offset, limit}

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
