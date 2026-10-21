package service

import "github.com/jjlee2712/coinflow/backend/db"

type BudgetAlert struct {
	CategoryID int64   `json:"categoryId"`
	Percentage float32 `json:"percentage"`
	Severity   string  `json:"severity"`
}

func CheckBudgetAlerts(summaries []db.GetBudgetSummaryRow) []BudgetAlert {
	var alerts []BudgetAlert

	for _, c := range summaries {
		var severity string

		if c.Percentage >= 100 {
			severity = "exceeded"
		} else if c.Percentage >= 80 {
			severity = "warning"
		} else {
			continue
		}

		alerts = append(alerts, BudgetAlert{
			CategoryID: c.CategoryID,
			Percentage: float32(c.Percentage),
			Severity:   severity,
		})
	}
	return alerts
}
