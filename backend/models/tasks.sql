-- name: CreateTask :one
INSERT INTO tasks (
  type,
  status,
  args,
  results,
  flow_id,
  message,
  tool_call_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: ReadTasksByFlowId :many
SELECT * FROM tasks
WHERE flow_id = ?
ORDER BY created_at ASC;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = ?
WHERE id = ?
RETURNING *;

-- name: UpdateTaskResults :one
UPDATE tasks
SET results = ?
WHERE id = ?
RETURNING *;

-- name: UpdateTaskToolCallId :one
UPDATE tasks
SET tool_call_id = ?
WHERE id = ?
RETURNING *;
