-- name: CreateFlow :one
INSERT INTO flows (
  name, status, container_id, model, model_provider
)
VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: ReadAllFlows :many
SELECT
  f.*,
  c.name AS container_name
FROM flows f
LEFT JOIN containers c ON f.container_id = c.id
ORDER BY f.created_at DESC;

-- name: ReadFlow :one
SELECT
  f.*,
  c.name AS container_name,
  c.image AS container_image,
  c.status AS container_status,
  c.local_id AS container_local_id
FROM flows f
LEFT JOIN containers c ON f.container_id = c.id
WHERE f.id = ?;

-- name: UpdateFlowStatus :one
UPDATE flows
SET status = ?
WHERE id = ?
RETURNING *;

-- name: UpdateFlowName :one
UPDATE flows
SET name = ?
WHERE id = ?
RETURNING *;

-- name: UpdateFlowContainer :one
UPDATE flows
SET container_id = ?
WHERE id = ?
RETURNING *;
