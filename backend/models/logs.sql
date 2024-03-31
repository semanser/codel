-- name: CreateLog :one
INSERT INTO logs (
  message, flow_id, type
)
VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: GetLogsByFlowId :many
SELECT *
FROM logs
WHERE flow_id = ?
ORDER BY created_at ASC;
