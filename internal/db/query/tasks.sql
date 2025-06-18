-- name: GetTasks :many
SELECT * FROM tasks;

-- name: CreateTask :one
INSERT INTO tasks (description, schedule, selector, command, scope, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?;
