package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jjlee2712/coinflow/backend/db"
	"github.com/jjlee2712/coinflow/backend/middleware"
)

type BudgetHandler struct {
	Queries *db.Queries
}

type budgetRequest struct {
	CategoryID int64  `json:"categoryId"`
	Amount     string `json:"amount"`
}

func (h *BudgetHandler) getUserID(r *http.Request) int64 {
	return r.Context().Value(middleware.UserIDKey).(int64)
}

func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	var req budgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body.", http.StatusBadRequest)
		return
	}

	budgets, err := h.Queries.CreateBudget(r.Context(), db.CreateBudgetParams{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(budgets)
}

func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)

	budgets, err := h.Queries.ListBudgetsByUser(r.Context(), userID)

	if err != nil {
		http.Error(w, "Budget Not found!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(budgets)
}

func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	budgetIdStr := r.PathValue("budgetId")
	budgetId, err := strconv.ParseInt(budgetIdStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var req budgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body.", http.StatusBadRequest)
		return
	}

	budgets, err := h.Queries.UpdateBudget(r.Context(), db.UpdateBudgetParams{
		Amount: req.Amount,
		ID:     budgetId,
		UserID: userID,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(budgets)
}

func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	budgetIdStr := r.PathValue("budgetId")
	budgetId, err := strconv.ParseInt(budgetIdStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err := h.Queries.DeleteBudget(r.Context(), db.DeleteBudgetParams{
		ID:     budgetId,
		UserID: userID,
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BudgetHandler) GetBudgetSummary(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)

	budgets, err := h.Queries.GetBudgetSummary(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(budgets)
}
