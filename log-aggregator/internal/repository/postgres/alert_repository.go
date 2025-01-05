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
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	_, err := r.db.Exec(
		query,
		alert.ID,
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

func (r *AlertRepository) FindByID(id string) (*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE id = $1`

	alert := &domain.Alert{}
	err := r.db.QueryRow(query, id).Scan(
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

func (r *AlertRepository) List(limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
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

func (r *AlertRepository) FindByStatus(status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	return r.findAlertsByQuery(query, status, limit, offset)
}

func (r *AlertRepository) FindBySeverity(severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error) {
	query := `
		SELECT 
			id, title, description, severity, status,
			source, created_at, updated_at, resolved_at
		FROM alerts
		WHERE severity = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	return r.findAlertsByQuery(query, severity, limit, offset)
}

func (r *AlertRepository) findAlertsByQuery(query string, param interface{}, limit, offset int) ([]*domain.Alert, error) {
	rows, err := r.db.Query(query, param, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find alerts: %w", err)
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
