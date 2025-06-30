package db

import (
	"database/sql"
	"errors"
	"os"

	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite"
)

var Q *Queries

func init() {
	_ = os.MkdirAll("data", 0755)

	dsn := "file:data/data.db"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	_, _ = db.Exec(`PRAGMA journal_mode = WAL;`)
	_, _ = db.Exec(`PRAGMA foreign_keys = ON;`)

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
		err := db.Close()
		q.db = nil
		return err
	}
	if db, ok := q.db.(*sql.Tx); ok {
		err := db.Rollback()
		q.db = nil
		return err
	}
	return errors.New("db is neither *sql.DB nor *sql.Tx")
}

func (q *Queries) Ping() error {
	if q.db == nil {
		return errors.New("db is nil")
	}
	if db, ok := q.db.(*sql.DB); ok {
		return db.Ping()
	}
	if db, ok := q.db.(*sql.Tx); ok {
		return db.Commit()
	}
	return nil
}

func (q *Queries) Begin() (*Queries, error) {
	if q.db == nil {
		return nil, errors.New("db is nil")
	}
	if db, ok := q.db.(*sql.DB); ok {
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		return &Queries{
			db: tx,
		}, nil
	}
	return nil, nil
}
