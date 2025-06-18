package db

import (
	"database/sql"
	"errors"

	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite"
)

var Q *Queries

func init() {
	dsn := "file:data.db?_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	Q = New(db)
}

func (q *Queries) Close() error {
	if q.db == nil {
		return nil
	}
	if db, ok := q.db.(*sql.DB); ok {
		return db.Close()
	}
	return nil
}

func (q *Queries) Ping() error {
	if q.db == nil {
		return errors.New("db is nil")
	}
	if db, ok := q.db.(*sql.DB); ok {
		return db.Ping()
	}
	return nil
}
