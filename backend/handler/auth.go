package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jjlee2712/coinflow/backend/db"
	"github.com/jjlee2712/coinflow/backend/service"
)

type AuthHandler struct {
	Queries *db.Queries
	JWTSecret string
}

type registerRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}


func(h *AuthHandler) Register(w http.ResponseWriter, r *http.Request){
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); 
	err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	hashed, err := service.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusBadRequest)
		return
	}

	user, err := h.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Email: req.Email,
		Password: hashed,
	})

	if err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func(h *AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req)
	err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	user, err := h.Queries.GetUserByEmail(r.Context(), req.Email);
	if err != nil {
		http.Error(w, "Invalid Email or Credentials", http.StatusUnauthorized)
    return
	}

	if err := service.CheckPassowrd(req.Password, user.Password)
	err != nil {
		http.Error(w, "Invalid Email or Credentials", http.StatusUnauthorized)
    return
	}

  token, err := service.CreateToken(user.ID, h.JWTSecret);
  if err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

	w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(loginResponse{Token: token})
}
