-- name: CreateExpense :one
INSERT INTO expenses (user_id, category_id, amount, note, date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListExpensesByUser :many
SELECT * FROM expenses
WHERE user_id = $1
ORDER BY date DESC;

-- name: GetExpense :one
SELECT * FROM expenses
WHERE id = $1 AND user_id = $2;

-- name: UpdateExpense :one
UPDATE expenses
SET category_id = $1, amount = $2, note = $3, date = $4
WHERE id = $5 AND user_id = $6
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1 AND user_id = $2;