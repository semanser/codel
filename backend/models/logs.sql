-- name: CreateLog :one
INSERT INTO logs (
  message, flow_id, type
)
VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetLogsByFlowId :many
SELECT *
FROM logs
WHERE flow_id = $1
ORDER BY created_at DESC;
