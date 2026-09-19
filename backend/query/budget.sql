-- name: CreateBudget :one
INSERT INTO budgets (user_id, category_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListBudgetsByUser :many
SELECT * FROM budgets
WHERE user_id = $1;

-- name: GetBudget :one
SELECT * FROM budgets
WHERE id = $1 AND user_id = $2;

-- name: UpdateBudget :one
UPDATE budgets
SET amount = $1
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budgets
WHERE id = $1 AND user_id = $2;