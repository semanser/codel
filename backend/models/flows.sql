-- name: InsertFlow :one
INSERT INTO flows (
  name, status, container_id
)
VALUES (
  $1, $2, $3
)
RETURNING *;
