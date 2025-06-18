package db

import (
	"database/sql"
	"embed"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
)

//go:embed migrations/*.sql
var migrations embed.FS

func migrate(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	for name, content := range getAllMigrations {
		if slices.Contains(appliedMigrations, name) {
			continue
		}
		log.Info("Applying migration", "name", name)
		if _, err := tx.Exec(string(content)); err != nil {
			return err
		}
		if _, err := tx.Exec(`
			INSERT INTO migrations (name) VALUES (?);
		`, name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func getAllMigrations(yield func(name string, content []byte) bool) {
	_ = fs.WalkDir(migrations, ".", func(filename string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := migrations.ReadFile(filename)
		if err != nil {
			return err
		}
		name := filepath.Base(filename[:len(filename)-4])
		if !yield(name, content) {
			return fs.SkipAll
		}
		return nil
	})
}

func getAppliedMigrations(db *sql.DB) ([]string, error) {
	var migrations []string
	rows, err := db.Query(`
		SELECT name FROM migrations ORDER BY created_at ASC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations = append(migrations, name)
	}

	return migrations, nil
}
