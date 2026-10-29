package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/jjlee2712/coinflow/backend/config"
	"github.com/jjlee2712/coinflow/backend/db"
	"github.com/jjlee2712/coinflow/backend/handler"
	"github.com/jjlee2712/coinflow/backend/hub"
	"github.com/jjlee2712/coinflow/backend/middleware"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	cfg := config.Load()

	database, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Printf("skipping integration tests: cannot open DB: %v\n", err)
		os.Exit(0)
	}
	if err := database.Ping(); err != nil {
		fmt.Printf("skipping integration tests: cannot reach DB: %v\n", err)
		os.Exit(0)
	}
	if err := runMigrations(database); err != nil {
		fmt.Printf("skipping integration tests: migrations failed: %v\n", err)
		os.Exit(1)
	}

	queries := db.New(database)
	h := hub.NewHub()

	authHandler := &handler.AuthHandler{Queries: queries, JWTSecret: cfg.JWTSecret}
	categoryHandler := &handler.CategoryHandler{Queries: queries}
	expensesHandler := &handler.ExpensesHandler{Queries: queries, Hub: h}
	budgetHandler := &handler.BudgetHandler{Queries: queries}
	auth := middleware.AuthMiddleware(cfg.JWTSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.Handle("GET /categories", auth(http.HandlerFunc(categoryHandler.List)))
	mux.Handle("POST /expenses", auth(http.HandlerFunc(expensesHandler.Create)))
	mux.Handle("GET /expenses", auth(http.HandlerFunc(expensesHandler.List)))
	mux.Handle("POST /budgets", auth(http.HandlerFunc(budgetHandler.Create)))
	mux.Handle("GET /budgets", auth(http.HandlerFunc(budgetHandler.List)))
	mux.Handle("GET /budgets/summary", auth(http.HandlerFunc(budgetHandler.GetBudgetSummary)))

	testServer = httptest.NewServer(corsMiddleware(mux))
	defer testServer.Close()

	os.Exit(m.Run())
}

func doRequest(t *testing.T, method, path string, body any, token string) *http.Response {
	t.Helper()
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
	}
	req, err := http.NewRequest(method, testServer.URL+path, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

// registers a unique user and returns their JWT
func loginAsTestUser(t *testing.T) string {
	t.Helper()
	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	resp := doRequest(t, "POST", "/auth/register", map[string]string{
		"email": email, "password": "testpassword123",
	}, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d", resp.StatusCode)
	}

	resp = doRequest(t, "POST", "/auth/login", map[string]string{
		"email": email, "password": "testpassword123",
	}, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: expected 200, got %d", resp.StatusCode)
	}
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["token"]
}

// --- Auth tests ---

func TestHealth(t *testing.T) {
	resp := doRequest(t, "GET", "/health", nil, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRegister_Success(t *testing.T) {
	email := fmt.Sprintf("reg-%d@example.com", time.Now().UnixNano())
	resp := doRequest(t, "POST", "/auth/register", map[string]string{
		"email": email, "password": "password123",
	}, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	email := fmt.Sprintf("dup-%d@example.com", time.Now().UnixNano())
	body := map[string]string{"email": email, "password": "password123"}
	doRequest(t, "POST", "/auth/register", body, "").Body.Close()

	resp := doRequest(t, "POST", "/auth/register", body, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 on duplicate, got %d", resp.StatusCode)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	resp := doRequest(t, "POST", "/auth/login", map[string]string{
		"email": "nobody@example.com", "password": "wrong",
	}, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

// --- Expenses tests ---

func TestExpenses_RequiresAuth(t *testing.T) {
	resp := doRequest(t, "GET", "/expenses", nil, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", resp.StatusCode)
	}
}

func TestCreateAndListExpenses(t *testing.T) {
	token := loginAsTestUser(t)

	// fetch default categories created on register
	resp := doRequest(t, "GET", "/categories", nil, token)
	var categories []map[string]any
	json.NewDecoder(resp.Body).Decode(&categories)
	resp.Body.Close()
	if len(categories) == 0 {
		t.Fatal("expected default categories to exist after register")
	}
	categoryID := categories[0]["ID"]

	// create expense
	resp = doRequest(t, "POST", "/expenses", map[string]any{
		"categoryId": categoryID,
		"amount":     "50.00",
		"date":       "2026-10-01",
		"note":       "Test expense",
	}, token)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create expense: expected 201, got %d", resp.StatusCode)
	}

	// list expenses
	resp = doRequest(t, "GET", "/expenses", nil, token)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list expenses: expected 200, got %d", resp.StatusCode)
	}
	var expenses []map[string]any
	json.NewDecoder(resp.Body).Decode(&expenses)
	if len(expenses) == 0 {
		t.Error("expected at least 1 expense after creating one")
	}
}

// --- Budget tests ---

func TestBudgets_RequiresAuth(t *testing.T) {
	resp := doRequest(t, "GET", "/budgets", nil, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", resp.StatusCode)
	}
}

func TestCreateBudgetAndGetSummary(t *testing.T) {
	token := loginAsTestUser(t)

	// get category
	resp := doRequest(t, "GET", "/categories", nil, token)
	var categories []map[string]any
	json.NewDecoder(resp.Body).Decode(&categories)
	resp.Body.Close()
	categoryID := categories[0]["ID"]

	// create budget
	resp = doRequest(t, "POST", "/budgets", map[string]any{
		"categoryId": categoryID,
		"amount":     "100.00",
	}, token)
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create budget: expected 201, got %d", resp.StatusCode)
	}

	// get summary
	resp = doRequest(t, "GET", "/budgets/summary", nil, token)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("budget summary: expected 200, got %d", resp.StatusCode)
	}
	var summary []map[string]any
	json.NewDecoder(resp.Body).Decode(&summary)
	if len(summary) == 0 {
		t.Error("expected at least 1 row in budget summary")
	}
}

func TestBudgetSummary_EmptyWithoutBudgets(t *testing.T) {
	token := loginAsTestUser(t)

	resp := doRequest(t, "GET", "/budgets/summary", nil, token)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var summary []map[string]any
	json.NewDecoder(resp.Body).Decode(&summary)
	if len(summary) != 0 {
		t.Errorf("expected empty summary for new user, got %d rows", len(summary))
	}
}
