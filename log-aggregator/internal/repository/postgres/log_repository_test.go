package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/travism26/log-aggregator/internal/domain"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	return db, mock
}

func TestLogRepository_UserIsolation(t *testing.T) {
	t.Run("Store enforces user isolation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		log := &domain.Log{
			ID:        "test-log",
			UserID:    "user1",
			APIKey:    "key123",
			Host:      "test-host",
			Message:   "test message",
			Level:     "INFO",
			Timestamp: time.Now(),
		}

		// Expect the INSERT query to include user_id
		mock.ExpectExec("INSERT INTO logs").
			WithArgs(
				log.ID,
				log.APIKey,
				log.UserID, // Verify user_id is included
				log.Timestamp,
				log.Host,
				log.Message,
				log.Level,
				sqlmock.AnyArg(), // metadata
				sqlmock.AnyArg(), // process_count
				sqlmock.AnyArg(), // total_cpu_percent
				sqlmock.AnyArg(), // total_memory_usage
				sqlmock.AnyArg(), // organization_id
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Store(log)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("FindByID enforces user isolation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		userID := "user1"
		logID := "test-log"

		rows := sqlmock.NewRows([]string{
			"id", "api_key", "organization_id", "timestamp", "host",
			"message", "level", "metadata", "process_count",
			"total_cpu_percent", "total_memory_usage",
		}).AddRow(
			logID, "key123", "org1", time.Now(), "test-host",
			"test message", "INFO", "{}", 10,
			50.0, 1024,
		)

		// Expect query to filter by both id AND user_id
		mock.ExpectQuery("SELECT .+ FROM logs WHERE id = \\$1 AND user_id = \\$2").
			WithArgs(logID, userID).
			WillReturnRows(rows)

		log, err := repo.FindByID(userID, logID)
		assert.NoError(t, err)
		assert.NotNil(t, log)
		assert.Equal(t, logID, log.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("List enforces user isolation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		userID := "user1"
		limit := 10
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "api_key", "organization_id", "timestamp", "host",
			"message", "level", "metadata", "process_count",
			"total_cpu_percent", "total_memory_usage",
		}).AddRow(
			"log1", "key123", "org1", time.Now(), "host1",
			"message1", "INFO", "{}", 10, 50.0, 1024,
		).AddRow(
			"log2", "key123", "org1", time.Now(), "host2",
			"message2", "ERROR", "{}", 15, 60.0, 2048,
		)

		// Expect query to filter by user_id in CTE
		mock.ExpectQuery("WITH recent_logs AS \\(SELECT .+ FROM logs WHERE user_id = \\$1").
			WithArgs(userID, limit, offset).
			WillReturnRows(rows)

		logs, err := repo.List(userID, limit, offset)
		assert.NoError(t, err)
		assert.Len(t, logs, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListByTimeRange enforces user isolation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		userID := "user1"
		start := time.Now().Add(-1 * time.Hour)
		end := time.Now()
		limit := 10
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "api_key", "organization_id", "timestamp", "host",
			"message", "level", "metadata", "process_count",
			"total_cpu_percent", "total_memory_usage",
		}).AddRow(
			"log1", "key123", "org1", time.Now(), "host1",
			"message1", "INFO", "{}", 10, 50.0, 1024,
		)

		// Expect query to filter by user_id AND time range
		mock.ExpectQuery("SELECT .+ FROM logs WHERE user_id = \\$1 AND timestamp >= \\$2 AND timestamp <= \\$3").
			WithArgs(userID, start, end, limit, offset).
			WillReturnRows(rows)

		logs, err := repo.ListByTimeRange(userID, start, end, limit, offset)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountByTimeRange enforces user isolation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		userID := "user1"
		start := time.Now().Add(-1 * time.Hour)
		end := time.Now()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

		// Expect query to filter by user_id AND time range
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM logs WHERE user_id = \\$1 AND timestamp >= \\$2 AND timestamp <= \\$3").
			WithArgs(userID, start, end).
			WillReturnRows(rows)

		count, err := repo.CountByTimeRange(userID, start, end)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListByHost enforces user isolation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		userID := "user1"
		host := "test-host"
		limit := 10
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "api_key", "organization_id", "timestamp", "host",
			"message", "level", "metadata", "process_count",
			"total_cpu_percent", "total_memory_usage",
		}).AddRow(
			"log1", "key123", "org1", time.Now(), host,
			"message1", "INFO", "{}", 10, 50.0, 1024,
		)

		// Expect query to filter by user_id AND host
		mock.ExpectQuery("SELECT .+ FROM logs WHERE user_id = \\$1 AND host = \\$2").
			WithArgs(userID, host, limit, offset).
			WillReturnRows(rows)

		logs, err := repo.ListByHost(userID, host, limit, offset)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, host, logs[0].Host)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListByLevel enforces user isolation", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		userID := "user1"
		level := "ERROR"
		limit := 10
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "api_key", "organization_id", "timestamp", "host",
			"message", "level", "metadata", "process_count",
			"total_cpu_percent", "total_memory_usage",
		}).AddRow(
			"log1", "key123", "org1", time.Now(), "host1",
			"error message", level, "{}", 10, 50.0, 1024,
		)

		// Expect query to filter by user_id AND level
		mock.ExpectQuery("SELECT .+ FROM logs WHERE user_id = \\$1 AND level = \\$2").
			WithArgs(userID, level, limit, offset).
			WillReturnRows(rows)

		logs, err := repo.ListByLevel(userID, level, limit, offset)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, level, logs[0].Level)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid user ID returns no data", func(t *testing.T) {
		db, mock := setupMockDB(t)
		defer db.Close()

		repo := NewLogRepository(db)
		invalidUserID := "invalid-user"
		logID := "test-log"

		// Expect query to return no rows for invalid user
		mock.ExpectQuery("SELECT .+ FROM logs WHERE id = \\$1 AND user_id = \\$2").
			WithArgs(logID, invalidUserID).
			WillReturnRows(sqlmock.NewRows([]string{}))

		log, err := repo.FindByID(invalidUserID, logID)
		assert.Error(t, err)
		assert.Nil(t, log)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
