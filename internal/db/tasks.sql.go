// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: tasks.sql

package db

import (
	"context"
)

const CreateTask = `-- name: CreateTask :one
INSERT INTO tasks (description, schedule, selector, command, scope, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, description, schedule, selector, command, scope, created_at, updated_at
`

// CreateTask
//
//	INSERT INTO tasks (description, schedule, selector, command, scope, created_at, updated_at)
//	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
//	RETURNING id, description, schedule, selector, command, scope, created_at, updated_at
func (q *Queries) CreateTask(ctx context.Context, description string, schedule string, selector string, command string, scope string) (Task, error) {
	row := q.db.QueryRowContext(ctx, CreateTask,
		description,
		schedule,
		selector,
		command,
		scope,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.Schedule,
		&i.Selector,
		&i.Command,
		&i.Scope,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const DeleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?
`

// DeleteTask
//
//	DELETE FROM tasks
//	WHERE id = ?
func (q *Queries) DeleteTask(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, DeleteTask, id)
	return err
}

const GetTasks = `-- name: GetTasks :many
SELECT id, description, schedule, selector, command, scope, created_at, updated_at FROM tasks
`

// GetTasks
//
//	SELECT id, description, schedule, selector, command, scope, created_at, updated_at FROM tasks
func (q *Queries) GetTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, GetTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Description,
			&i.Schedule,
			&i.Selector,
			&i.Command,
			&i.Scope,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
