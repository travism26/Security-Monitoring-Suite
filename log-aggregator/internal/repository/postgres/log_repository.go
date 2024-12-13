package postgres

import (
	"database/sql"

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
		INSERT INTO logs (id, timestamp, host, message, level, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, log.ID, log.Timestamp, log.Host,
		log.Message, log.Level, log.Metadata)
	return err
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
