-- name: GetAllRunningContainers :many
SELECT * FROM containers WHERE status = 'running';

-- name: CreateContainer :one
INSERT INTO containers (
  name, image, status
)
VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateContainerStatus :one
UPDATE containers
SET status = $1
WHERE id = $2
RETURNING *;

-- name: UpdateContainerLocalId :one
UPDATE containers
SET local_id = $1
WHERE id = $2
RETURNING *;
