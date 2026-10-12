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

-- name: GetBudgetSummary :many
SELECT 
EP.category_id,
SUM(BD.amount) as budgetAmount,
SUM(EP.amount) as totalSpent,
SUM(BD.amount) - SUM(EP.amount) as remaining,
(SUM(EP.amount) / SUM(BD.amount)) * 100 as percentage,
DATE_TRUNC('month', NOW())
FROM users U
LEFT JOIN expenses EP ON U.id = EP.user_id
LEFT JOIN budgets BD ON U.id = BD.user_id AND EP.category_id = BD.category_id
WHERE U.id = $1
GROUP BY EP.category_id;