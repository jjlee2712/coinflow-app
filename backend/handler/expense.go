package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jjlee2712/coinflow/backend/db"
	"github.com/jjlee2712/coinflow/backend/hub"
	"github.com/jjlee2712/coinflow/backend/middleware"
)

type ExpensesHandler struct {
	Queries *db.Queries
	Hub     *hub.Hub
}

type expensesRequest struct {
	CategoryID int64  `json:"categoryId"`
	Amount     string `json:"amount"`
	Note       string `json:"note"`
	Date       string `json:"date"`
}

type wsMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func (h *ExpensesHandler) getUserID(r *http.Request) int64 {
	return r.Context().Value(middleware.UserIDKey).(int64)
}

func (h *ExpensesHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserID(r)
	var req expensesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body.", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid Date Format, YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	expenses, err := h.Queries.CreateExpense(r.Context(), db.CreateExpenseParams{
		UserID:     userId,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Note:       sql.NullString{String: req.Note, Valid: req.Note != ""},
		Date:       date,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	go func() {
		budget, _ := h.Queries.GetBudgetSummary(context.Background(), userId)

		jsonBytes, err := json.Marshal(wsMessage{
			Type: "summary_update",
			Data: budget,
		})
		if err == nil {
			h.Hub.Broadcast(userId, jsonBytes)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpensesHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)

	categoryIdStr := r.URL.Query().Get("categoryId")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var expenses []db.Expense
	var err error

	if categoryIdStr != "" || fromStr != "" || toStr != "" {
		var categoryId int64
		if categoryIdStr != "" {
			categoryId, err = strconv.ParseInt(categoryIdStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid category_id", http.StatusBadRequest)
				return
			}
		}

		var from, to time.Time
		if fromStr != "" {
			from, err = time.Parse("2006-01-02", fromStr)
			if err != nil {
				http.Error(w, "Invalid from date, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}
		if toStr != "" {
			to, err = time.Parse("2006-01-02", toStr)
			if err != nil {
				http.Error(w, "Invalid to date, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}

		expenses, err = h.Queries.ListExpensesByFilter(r.Context(), db.ListExpensesByFilterParams{
			UserID:  userID,
			Column2: categoryId,
			Column3: from,
			Column4: to,
		})
	} else {
		expenses, err = h.Queries.ListExpensesByUser(r.Context(), userID)
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpensesHandler) GetExpensesByID(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	expensesIdStr := r.PathValue("expensesId")
	expensesId, err := strconv.ParseInt(expensesIdStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	expenses, err := h.Queries.GetExpense(r.Context(), db.GetExpenseParams{
		UserID: userID,
		ID:     expensesId,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpensesHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	expensesIdStr := r.PathValue("expensesId")
	expensesId, err := strconv.ParseInt(expensesIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var req expensesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid Date Format, YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	expenses, err := h.Queries.UpdateExpense(r.Context(), db.UpdateExpenseParams{
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Note:       sql.NullString{String: req.Note, Valid: req.Note != ""},
		Date:       date,
		ID:         expensesId,
		UserID:     userID,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpensesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	expensesIdStr := r.PathValue("expensesId")
	expensesId, err := strconv.ParseInt(expensesIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err := h.Queries.DeleteExpense(r.Context(), db.DeleteExpenseParams{
		ID:     expensesId,
		UserID: userID,
	}); err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	go func() {
		budget, _ := h.Queries.GetBudgetSummary(context.Background(), userID)

		jsonBytes, err := json.Marshal(wsMessage{
			Type: "summary_update",
			Data: budget,
		})
		if err == nil {
			h.Hub.Broadcast(userID, jsonBytes)
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}
