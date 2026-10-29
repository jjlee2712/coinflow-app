package service

import (
	"testing"

	"github.com/jjlee2712/coinflow/backend/db"
)

func makeRow(categoryID int64, name, pct string) db.GetBudgetSummaryRow {
	return db.GetBudgetSummaryRow{
		CategoryID:   categoryID,
		Categoryname: name,
		Percentage:   pct,
	}
}

func TestCheckBudgetAlerts_NoAlerts(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{
		makeRow(1, "Food", "50"),
		makeRow(2, "Transport", "79.9"),
	})
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestCheckBudgetAlerts_ExactlyAt80_IsWarning(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{
		makeRow(1, "Food", "80"),
	})
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != "warning" {
		t.Errorf("expected 'warning' at 80%%, got %q", alerts[0].Severity)
	}
	if alerts[0].CategoryName != "Food" {
		t.Errorf("expected category name 'Food', got %q", alerts[0].CategoryName)
	}
}

func TestCheckBudgetAlerts_Between80And100_IsWarning(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{
		makeRow(1, "Transport", "92.5"),
	})
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != "warning" {
		t.Errorf("expected 'warning', got %q", alerts[0].Severity)
	}
}

func TestCheckBudgetAlerts_ExactlyAt100_IsExceeded(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{
		makeRow(1, "Food", "100"),
	})
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != "exceeded" {
		t.Errorf("expected 'exceeded' at 100%%, got %q", alerts[0].Severity)
	}
}

func TestCheckBudgetAlerts_Above100_IsExceeded(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{
		makeRow(1, "Health", "150"),
	})
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != "exceeded" {
		t.Errorf("expected 'exceeded', got %q", alerts[0].Severity)
	}
}

func TestCheckBudgetAlerts_Mixed(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{
		makeRow(1, "Food", "50"),       // no alert
		makeRow(2, "Transport", "85"),  // warning
		makeRow(3, "Health", "110"),    // exceeded
	})
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestCheckBudgetAlerts_InvalidPercentage_IsSkipped(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{
		makeRow(1, "Food", "not-a-number"),
	})
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts for invalid percentage, got %d", len(alerts))
	}
}

func TestCheckBudgetAlerts_Empty(t *testing.T) {
	alerts := CheckBudgetAlerts([]db.GetBudgetSummaryRow{})
	if alerts != nil {
		t.Errorf("expected nil for empty input, got %v", alerts)
	}
}
