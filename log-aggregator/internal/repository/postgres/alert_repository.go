package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/travism26/log-aggregator/internal/domain"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{
		db: db,
	}
}

func (r *AlertRepository) Store(alert *domain.Alert) error {
	query := `
		INSERT INTO alerts (
			id, organization_id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	_, err := r.db.Exec(
		query,
		alert.ID,
		alert.OrganizationID,
		alert.Title,
		alert.Description,
		alert.Severity,
		alert.Status,
		alert.Source,
		alert.CreatedAt,
		alert.UpdatedAt,
		alert.ResolvedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to store alert: %w", err)
	}

	// Store related logs if any
	if len(alert.RelatedLogs) > 0 {
		for _, logID := range alert.RelatedLogs {
			err = r.storeAlertLogRelation(alert.ID, logID)
			if err != nil {
				return fmt.Errorf("failed to store alert-log relation: %w", err)
			}
		}
	}

	// Store metadata if any
	if len(alert.Metadata) > 0 {
		err = r.storeAlertMetadata(alert.ID, alert.Metadata)
		if err != nil {
			return fmt.Errorf("failed to store alert metadata: %w", err)
		}
	}

	return nil
}

func (r *AlertRepository) storeAlertLogRelation(alertID, logID string) error {
	query := `
		INSERT INTO alert_logs (alert_id, log_id)
		VALUES ($1, $2)`

	_, err := r.db.Exec(query, alertID, logID)
	return err
}

func (r *AlertRepository) storeAlertMetadata(alertID string, metadata map[string]interface{}) error {
	query := `
		INSERT INTO alert_metadata (alert_id, key, value)
		VALUES ($1, $2, $3)`

	for key, value := range metadata {
		_, err := r.db.Exec(query, alertID, key, fmt.Sprintf("%v", value))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AlertRepository) FindByID(orgID, id string) (*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE organization_id = $1 AND id = $2`

	alert := &domain.Alert{}
	err := r.db.QueryRow(query, orgID, id).Scan(
		&alert.ID,
		&alert.Title,
		&alert.Description,
		&alert.Severity,
		&alert.Status,
		&alert.Source,
		&alert.CreatedAt,
		&alert.UpdatedAt,
		&alert.ResolvedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("alert not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find alert: %w", err)
	}

	// Load related logs
	alert.RelatedLogs, err = r.getRelatedLogs(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get related logs: %w", err)
	}

	// Load metadata
	alert.Metadata, err = r.getAlertMetadata(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert metadata: %w", err)
	}

	return alert, nil
}

func (r *AlertRepository) getRelatedLogs(alertID string) ([]string, error) {
	query := `
		SELECT log_id
		FROM alert_logs
		WHERE alert_id = $1`

	rows, err := r.db.Query(query, alertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logIDs []string
	for rows.Next() {
		var logID string
		if err := rows.Scan(&logID); err != nil {
			return nil, err
		}
		logIDs = append(logIDs, logID)
	}

	return logIDs, nil
}

func (r *AlertRepository) getAlertMetadata(alertID string) (map[string]interface{}, error) {
	query := `
		SELECT key, value
		FROM alert_metadata
		WHERE alert_id = $1`

	rows, err := r.db.Query(query, alertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metadata := make(map[string]interface{})
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		metadata[key] = value
	}

	return metadata, nil
}

func (r *AlertRepository) List(orgID string, limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}
	defer rows.Close()

	var alerts []*domain.Alert
	for rows.Next() {
		alert := &domain.Alert{}
		err := rows.Scan(
			&alert.ID,
			&alert.Title,
			&alert.Description,
			&alert.Severity,
			&alert.Status,
			&alert.Source,
			&alert.CreatedAt,
			&alert.UpdatedAt,
			&alert.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		// Load related logs
		alert.RelatedLogs, err = r.getRelatedLogs(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get related logs: %w", err)
		}

		// Load metadata
		alert.Metadata, err = r.getAlertMetadata(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get alert metadata: %w", err)
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *AlertRepository) Update(alert *domain.Alert) error {
	query := `
		UPDATE alerts
		SET 
			title = $1,
			description = $2,
			severity = $3,
			status = $4,
			source = $5,
			updated_at = $6,
			resolved_at = $7
		WHERE id = $8`

	result, err := r.db.Exec(
		query,
		alert.Title,
		alert.Description,
		alert.Severity,
		alert.Status,
		alert.Source,
		time.Now(),
		alert.ResolvedAt,
		alert.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("alert not found: %s", alert.ID)
	}

	return nil
}

func (r *AlertRepository) FindByStatus(orgID string, status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE organization_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, orgID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find alerts by status: %w", err)
	}
	defer rows.Close()

	return r.scanAlerts(rows)
}

func (r *AlertRepository) FindBySeverity(orgID string, severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE organization_id = $1 AND severity = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, orgID, severity, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find alerts by severity: %w", err)
	}
	defer rows.Close()

	return r.scanAlerts(rows)
}

func (r *AlertRepository) CountBySeverity(orgID string, severity domain.AlertSeverity) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM alerts
		WHERE organization_id = $1 AND severity = $2`

	var count int64
	err := r.db.QueryRow(query, orgID, severity).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count alerts by severity: %w", err)
	}
	return count, nil
}

func (r *AlertRepository) CountByStatus(orgID string, status domain.AlertStatus) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM alerts
		WHERE organization_id = $1 AND status = $2`

	var count int64
	err := r.db.QueryRow(query, orgID, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count alerts by status: %w", err)
	}
	return count, nil
}

func (r *AlertRepository) CountByTimeRange(orgID string, start, end time.Time) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM alerts
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3`

	var count int64
	err := r.db.QueryRow(query, orgID, start, end).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count alerts by time range: %w", err)
	}
	return count, nil
}

func (r *AlertRepository) CountBySource(orgID string, source string) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM alerts
		WHERE organization_id = $1 AND source = $2`

	var count int64
	err := r.db.QueryRow(query, orgID, source).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count alerts by source: %w", err)
	}
	return count, nil
}

func (r *AlertRepository) Delete(orgID, id string) error {
	// Start a transaction since we need to delete from multiple tables
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if we return with error

	// Delete alert metadata
	_, err = tx.Exec(`DELETE FROM alert_metadata WHERE alert_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert metadata: %w", err)
	}

	// Delete alert-log relations
	_, err = tx.Exec(`DELETE FROM alert_logs WHERE alert_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert-log relations: %w", err)
	}

	// Delete the alert
	result, err := tx.Exec(`DELETE FROM alerts WHERE organization_id = $1 AND id = $2`, orgID, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("alert not found: %s", id)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *AlertRepository) FindBySource(orgID string, source string, limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE organization_id = $1 AND source = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, orgID, source, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find alerts by source: %w", err)
	}
	defer rows.Close()

	var alerts []*domain.Alert
	for rows.Next() {
		alert := &domain.Alert{}
		err := rows.Scan(
			&alert.ID,
			&alert.Title,
			&alert.Description,
			&alert.Severity,
			&alert.Status,
			&alert.Source,
			&alert.CreatedAt,
			&alert.UpdatedAt,
			&alert.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		// Load related logs
		alert.RelatedLogs, err = r.getRelatedLogs(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get related logs: %w", err)
		}

		// Load metadata
		alert.Metadata, err = r.getAlertMetadata(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get alert metadata: %w", err)
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *AlertRepository) ListByTimeRange(orgID string, start, end time.Time, limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5`

	rows, err := r.db.Query(query, orgID, start, end, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by time range: %w", err)
	}
	defer rows.Close()

	var alerts []*domain.Alert
	for rows.Next() {
		alert := &domain.Alert{}
		err := rows.Scan(
			&alert.ID,
			&alert.Title,
			&alert.Description,
			&alert.Severity,
			&alert.Status,
			&alert.Source,
			&alert.CreatedAt,
			&alert.UpdatedAt,
			&alert.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		// Load related logs
		alert.RelatedLogs, err = r.getRelatedLogs(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get related logs: %w", err)
		}

		// Load metadata
		alert.Metadata, err = r.getAlertMetadata(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get alert metadata: %w", err)
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// Helper function to scan alert rows into alert structs
func (r *AlertRepository) scanAlerts(rows *sql.Rows) ([]*domain.Alert, error) {
	var alerts []*domain.Alert
	for rows.Next() {
		alert := &domain.Alert{}
		err := rows.Scan(
			&alert.ID,
			&alert.Title,
			&alert.Description,
			&alert.Severity,
			&alert.Status,
			&alert.Source,
			&alert.CreatedAt,
			&alert.UpdatedAt,
			&alert.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		// Load related logs
		alert.RelatedLogs, err = r.getRelatedLogs(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get related logs: %w", err)
		}

		// Load metadata
		alert.Metadata, err = r.getAlertMetadata(alert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get alert metadata: %w", err)
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}
