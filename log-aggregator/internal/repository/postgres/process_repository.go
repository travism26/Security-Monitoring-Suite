package postgres

import (
	"database/sql"
	"fmt"

	"github.com/travism26/log-aggregator/internal/domain"

	_ "github.com/lib/pq"
)

type ProcessRepository struct {
	db *sql.DB
}

func NewProcessRepository(db *sql.DB) *ProcessRepository {
	return &ProcessRepository{db: db}
}

func (r *ProcessRepository) StoreBatch(processes []domain.Process) error {
	query := `
        INSERT INTO process_logs (
            id,
            log_id,
            name,
            pid,
            cpu_percent,
            memory_usage,
            status,
            created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	// Debug logging - Start
	fmt.Printf("=== Debug StoreBatch Method ===\n")
	fmt.Printf("Number of processes to store: %d\n", len(processes))
	if len(processes) > 0 {
		fmt.Printf("Sample process data (first record):\n")
		fmt.Printf("ID: %s\n", processes[0].ID)
		fmt.Printf("LogID: %s\n", processes[0].LogID)
		fmt.Printf("Name: %s\n", processes[0].Name)
		fmt.Printf("PID: %d\n", processes[0].PID)
		fmt.Printf("CPU Percent: %.2f\n", processes[0].CPUPercent)
		fmt.Printf("Memory Usage: %d\n", processes[0].MemoryUsage)
		fmt.Printf("Status: %s\n", processes[0].Status)
		fmt.Printf("Timestamp: %v\n", processes[0].Timestamp)
	}
	fmt.Printf("===========================\n")

	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("=== Error Details ===\n")
		fmt.Printf("Transaction begin error: %v\n", err)
		fmt.Printf("Error Type: %T\n", err)
		fmt.Printf("==================\n")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Store each process as a separate row
	for i, process := range processes {
		_, err = tx.Exec(
			query,
			process.ID,
			process.LogID,
			process.Name,
			process.PID,
			process.CPUPercent,
			process.MemoryUsage,
			process.Status,
			process.Timestamp,
		)
		if err != nil {
			fmt.Printf("=== Error Details ===\n")
			fmt.Printf("Error on process %d/%d\n", i+1, len(processes))
			fmt.Printf("Process ID: %s\n", process.ID)
			fmt.Printf("Process Name: %s\n", process.Name)
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Error Type: %T\n", err)
			fmt.Printf("==================\n")
			return fmt.Errorf("failed to execute query for process %d: %w", process.PID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("=== Error Details ===\n")
		fmt.Printf("Transaction commit error: %v\n", err)
		fmt.Printf("Error Type: %T\n", err)
		fmt.Printf("==================\n")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("=== Success ===\n")
	fmt.Printf("Successfully stored %d processes\n", len(processes))
	fmt.Printf("==============\n")

	return nil
}
