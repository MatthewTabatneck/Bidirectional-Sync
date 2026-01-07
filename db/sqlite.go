package db

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

type Store struct {
	db *sql.DB
}

// NewStore initializes the SQLite database and creates the table if it doesn't exist.
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the files table
	// We use the 'path' as the primary key because it's unique to each file
	query := `
	CREATE TABLE IF NOT EXISTS files (
		path TEXT PRIMARY KEY,
		hash TEXT NOT NULL,
		size INTEGER NOT NULL,
		mod_time DATETIME NOT NULL,
		is_deleted INTEGER DEFAULT 0
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// UpsertFile inserts a new file record or updates an existing one if the path matches.
// Using "ON CONFLICT" makes this an "Upsert" operation.
func (s *Store) UpsertFile(path string, hash string, size int64, modTime time.Time) error {
	query := `
    INSERT INTO files (path, hash, size, mod_time, is_deleted)
    VALUES (?, ?, ?, ?, 0)
    ON CONFLICT(path) DO UPDATE SET
        hash = excluded.hash,
        size = excluded.size,
        mod_time = excluded.mod_time,
        is_deleted = 0;`

	_, err := s.db.Exec(query, path, hash, size, modTime)
	if err != nil {
		return err
	}
	return nil
}

// Close safely shuts down the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
