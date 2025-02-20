package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/travism26/log-aggregator/internal/domain"
)

func BenchmarkLogRepository_UserIsolation(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping performance test in short mode")
	}

	db := setupTestDB(b)
	defer db.Close()

	repo := NewLogRepository(db)
	now := time.Now()

	// Setup test data
	numUsers := 10
	logsPerUser := 1000
	users := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		users[i] = fmt.Sprintf("user%d", i)
	}

	// Create test logs for each user
	b.Run("Setup_InsertTestData", func(b *testing.B) {
		for _, userID := range users {
			for j := 0; j < logsPerUser; j++ {
				log := &domain.Log{
					ID:        fmt.Sprintf("log-%s-%d", userID, j),
					UserID:    userID,
					APIKey:    fmt.Sprintf("key-%s", userID),
					Host:      fmt.Sprintf("host-%d", j%5),
					Message:   fmt.Sprintf("test message %d", j),
					Level:     []string{"INFO", "WARN", "ERROR"}[j%3],
					Timestamp: now.Add(time.Duration(j) * time.Minute),
				}
				if err := repo.Store(log); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	// Benchmark List operation with user isolation
	b.Run("List_WithUserIsolation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := users[i%numUsers]
			if _, err := repo.List(userID, 100, 0); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark time range queries with user isolation
	b.Run("ListByTimeRange_WithUserIsolation", func(b *testing.B) {
		start := now
		end := now.Add(30 * time.Minute)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := users[i%numUsers]
			if _, err := repo.ListByTimeRange(userID, start, end, 100, 0); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark host filtering with user isolation
	b.Run("ListByHost_WithUserIsolation", func(b *testing.B) {
		host := "host-1"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := users[i%numUsers]
			if _, err := repo.ListByHost(userID, host, 100, 0); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark level filtering with user isolation
	b.Run("ListByLevel_WithUserIsolation", func(b *testing.B) {
		level := "ERROR"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := users[i%numUsers]
			if _, err := repo.ListByLevel(userID, level, 100, 0); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark count queries with user isolation
	b.Run("CountByTimeRange_WithUserIsolation", func(b *testing.B) {
		start := now
		end := now.Add(30 * time.Minute)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := users[i%numUsers]
			if _, err := repo.CountByTimeRange(userID, start, end); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark FindByID with user isolation
	b.Run("FindByID_WithUserIsolation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := users[i%numUsers]
			logID := fmt.Sprintf("log-%s-%d", userID, i%logsPerUser)
			if _, err := repo.FindByID(userID, logID); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark concurrent access from multiple users
	b.Run("ConcurrentAccess_MultipleUsers", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				userID := users[time.Now().Nanosecond()%numUsers]
				if _, err := repo.List(userID, 100, 0); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
