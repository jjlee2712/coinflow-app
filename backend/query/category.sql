-- name: CreateCategory :one
INSERT INTO categories(user_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: ListCategoriesByUser :many
SELECT * FROM categories
WHERE user_id = $1;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 AND user_id = $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $1
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;