package postgres

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/travism26/log-aggregator/internal/domain"
)

type testingTB interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
}

func setupTestDB(tb testingTB) *sql.DB {
	tb.Helper()
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/log_aggregator_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		tb.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		tb.Fatal(err)
	}

	// Clean up database before tests
	if _, err := db.Exec("TRUNCATE TABLE logs"); err != nil {
		tb.Fatal(err)
	}

	return db
}

func TestLogRepository_UserIsolation_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	defer db.Close()

	repo := NewLogRepository(db)

	// Test data
	user1ID := "user1"
	user2ID := "user2"
	now := time.Now()

	// Create logs for user1
	user1Logs := []*domain.Log{
		{
			ID:        "log1-user1",
			UserID:    user1ID,
			APIKey:    "key1",
			Host:      "host1",
			Message:   "user1 message 1",
			Level:     "INFO",
			Timestamp: now,
		},
		{
			ID:        "log2-user1",
			UserID:    user1ID,
			APIKey:    "key1",
			Host:      "host2",
			Message:   "user1 message 2",
			Level:     "ERROR",
			Timestamp: now.Add(time.Minute),
		},
	}

	// Create logs for user2
	user2Logs := []*domain.Log{
		{
			ID:        "log1-user2",
			UserID:    user2ID,
			APIKey:    "key2",
			Host:      "host1",
			Message:   "user2 message 1",
			Level:     "INFO",
			Timestamp: now,
		},
		{
			ID:        "log2-user2",
			UserID:    user2ID,
			APIKey:    "key2",
			Host:      "host2",
			Message:   "user2 message 2",
			Level:     "ERROR",
			Timestamp: now.Add(time.Minute),
		},
	}

	t.Run("Store and retrieve logs with user isolation", func(t *testing.T) {
		// Store logs for both users
		for _, log := range user1Logs {
			err := repo.Store(log)
			require.NoError(t, err)
		}
		for _, log := range user2Logs {
			err := repo.Store(log)
			require.NoError(t, err)
		}

		// Test user1 can only see their logs
		user1Results, err := repo.List(user1ID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, user1Results, 2)
		for _, log := range user1Results {
			assert.Equal(t, user1ID, log.UserID)
		}

		// Test user2 can only see their logs
		user2Results, err := repo.List(user2ID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, user2Results, 2)
		for _, log := range user2Results {
			assert.Equal(t, user2ID, log.UserID)
		}
	})

	t.Run("Time range queries respect user isolation", func(t *testing.T) {
		start := now.Add(-time.Hour)
		end := now.Add(time.Hour)

		// Get logs for time range for each user
		user1Results, err := repo.ListByTimeRange(user1ID, start, end, 10, 0)
		require.NoError(t, err)
		assert.Len(t, user1Results, 2)
		for _, log := range user1Results {
			assert.Equal(t, user1ID, log.UserID)
		}

		user2Results, err := repo.ListByTimeRange(user2ID, start, end, 10, 0)
		require.NoError(t, err)
		assert.Len(t, user2Results, 2)
		for _, log := range user2Results {
			assert.Equal(t, user2ID, log.UserID)
		}
	})

	t.Run("Host filtering respects user isolation", func(t *testing.T) {
		// Test host filtering for each user
		user1Results, err := repo.ListByHost(user1ID, "host1", 10, 0)
		require.NoError(t, err)
		assert.Len(t, user1Results, 1)
		assert.Equal(t, user1ID, user1Results[0].UserID)
		assert.Equal(t, "host1", user1Results[0].Host)

		user2Results, err := repo.ListByHost(user2ID, "host1", 10, 0)
		require.NoError(t, err)
		assert.Len(t, user2Results, 1)
		assert.Equal(t, user2ID, user2Results[0].UserID)
		assert.Equal(t, "host1", user2Results[0].Host)
	})

	t.Run("Level filtering respects user isolation", func(t *testing.T) {
		// Test level filtering for each user
		user1Results, err := repo.ListByLevel(user1ID, "ERROR", 10, 0)
		require.NoError(t, err)
		assert.Len(t, user1Results, 1)
		assert.Equal(t, user1ID, user1Results[0].UserID)
		assert.Equal(t, "ERROR", user1Results[0].Level)

		user2Results, err := repo.ListByLevel(user2ID, "ERROR", 10, 0)
		require.NoError(t, err)
		assert.Len(t, user2Results, 1)
		assert.Equal(t, user2ID, user2Results[0].UserID)
		assert.Equal(t, "ERROR", user2Results[0].Level)
	})

	t.Run("Count queries respect user isolation", func(t *testing.T) {
		start := now.Add(-time.Hour)
		end := now.Add(time.Hour)

		// Test count for each user
		user1Count, err := repo.CountByTimeRange(user1ID, start, end)
		require.NoError(t, err)
		assert.Equal(t, int64(2), user1Count)

		user2Count, err := repo.CountByTimeRange(user2ID, start, end)
		require.NoError(t, err)
		assert.Equal(t, int64(2), user2Count)
	})

	t.Run("FindByID respects user isolation", func(t *testing.T) {
		// User1 can access their own log
		log1, err := repo.FindByID(user1ID, "log1-user1")
		require.NoError(t, err)
		assert.NotNil(t, log1)
		assert.Equal(t, user1ID, log1.UserID)

		// User2 cannot access user1's log
		_, err = repo.FindByID(user2ID, "log1-user1")
		assert.Error(t, err)

		// User2 can access their own log
		log2, err := repo.FindByID(user2ID, "log1-user2")
		require.NoError(t, err)
		assert.NotNil(t, log2)
		assert.Equal(t, user2ID, log2.UserID)

		// User1 cannot access user2's log
		_, err = repo.FindByID(user1ID, "log1-user2")
		assert.Error(t, err)
	})
}
