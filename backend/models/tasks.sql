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
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ReadTasksByFlowId :many
SELECT * FROM tasks
WHERE flow_id = $1;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $1
WHERE id = $2
RETURNING *;

-- name: UpdateTaskResults :one
UPDATE tasks
SET results = $1
WHERE id = $2
RETURNING *;

-- name: UpdateTaskToolCallId :one
UPDATE tasks
SET tool_call_id = $1
WHERE id = $2
RETURNING *;
