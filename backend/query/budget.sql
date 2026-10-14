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
BD.category_id,
SUM(BD.amount) as budgetAmount,
COALESCE(SUM(EP.amount), 0) as totalSpent,
SUM(BD.amount) - SUM(EP.amount) as remaining,
(SUM(EP.amount) / SUM(BD.amount)) * 100 as percentage,
DATE_TRUNC('month', NOW())
FROM budgets BD
LEFT JOIN expenses EP ON BD.category_id = EP.category_id
    AND EP.user_id = BD.user_id
    AND DATE_TRUNC('month', EP.date) = DATE_TRUNC('month', NOW())
WHERE BD.user_id = $1
GROUP BY BD.id, BD.category_id, BD.amount;