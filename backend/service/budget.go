package service

import (
	"strconv"

	"github.com/jjlee2712/coinflow/backend/db"
)

type BudgetAlert struct {
	CategoryID   int64   `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Percentage   float32 `json:"percentage"`
	Severity     string  `json:"severity"`
}

func CheckBudgetAlerts(summaries []db.GetBudgetSummaryRow) []BudgetAlert {
	var alerts []BudgetAlert

	for _, c := range summaries {
		pct, err := strconv.ParseFloat(c.Percentage, 64)
		if err != nil {
			continue
		}

		var severity string
		if pct >= 100 {
			severity = "exceeded"
		} else if pct >= 80 {
			severity = "warning"
		} else {
			continue
		}

		alerts = append(alerts, BudgetAlert{
			CategoryID:   c.CategoryID,
			CategoryName: c.Categoryname,
			Percentage:   float32(pct),
			Severity:     severity,
		})
	}
	return alerts
}
