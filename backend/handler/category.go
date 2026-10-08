package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jjlee2712/coinflow/backend/db"
	"github.com/jjlee2712/coinflow/backend/middleware"
)

type CategoryHandler struct {
	Queries *db.Queries
}

type categoryRequest struct {
	Name string `json:"name"`
}

func (h *CategoryHandler) getUserID(r *http.Request) int64{
	return r.Context().Value(middleware.UserIDKey).(int64)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request){
	userID := h.getUserID(r)

	categories, err := h.Queries.ListCategoriesByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) Create (w http.ResponseWriter, r *http.Request){
	userID := h.getUserID(r)
	var req	categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req)
	err != nil {
		http.Error(w, "Invalid Request Body.", http.StatusBadRequest)
		return		
	}

	category, err := h.Queries.CreateCategory(r.Context(), db.CreateCategoryParams{
		UserID: userID,
		Name: req.Name,
	})
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request){
	userID := h.getUserID(r)
	idStr := r.PathValue("categoryId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req)
	err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	category, err := h.Queries.UpdateCategory(r.Context(), db.UpdateCategoryParams{
		Name: req.Name,
		ID: id,
		UserID: userID,
	})

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request){
	userID := h.getUserID(r)
	categoryIdStr := r.PathValue("categoryId")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	if err := h.Queries.DeleteCategory(r.Context(), db.DeleteCategoryParams{
		ID: categoryId,
		UserID: userID,
	})
	err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}