package postgres

import (
	"database/sql"

	"github.com/travism26/log-aggregator/internal/domain"
)

type ProcessRepository struct {
	db *sql.DB
}

func NewProcessRepository(db *sql.DB) *ProcessRepository {
	return &ProcessRepository{db: db}
}

func (r *ProcessRepository) StoreBatch(processes []domain.Process) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO process_logs (id, log_id, name, pid, cpu_percent, memory_usage, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, process := range processes {
		_, err = stmt.Exec(
			process.ID,
			process.LogID,
			process.Name,
			process.PID,
			process.CPUPercent,
			process.MemoryUsage,
			process.Status,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
