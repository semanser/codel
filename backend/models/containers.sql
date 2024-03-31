-- name: GetAllRunningContainers :many
SELECT * FROM containers WHERE status = 'running';

-- name: CreateContainer :one
INSERT INTO containers (
  name, image, status
)
VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateContainerStatus :one
UPDATE containers
SET status = ?
WHERE id = ?
RETURNING *;

-- name: UpdateContainerLocalId :one
UPDATE containers
SET local_id = ?
WHERE id = ?
RETURNING *;
