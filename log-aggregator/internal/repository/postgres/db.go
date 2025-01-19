package postgres

import (
	"database/sql"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

// SetDB sets the database connection for the package
func SetDB(database *sql.DB) {
	once.Do(func() {
		db = database
	})
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
