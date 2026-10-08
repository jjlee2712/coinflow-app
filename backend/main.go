package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"

	"github.com/jjlee2712/coinflow/backend/config"
	"github.com/jjlee2712/coinflow/backend/db"
	"github.com/jjlee2712/coinflow/backend/handler"
	"github.com/jjlee2712/coinflow/backend/middleware"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func runMigrations(db *sql.DB) error {
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func main() {
	cfg := config.Load()

	database, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err := runMigrations(database); err != nil {
		log.Fatal(err)
	}

	queries := db.New(database)

	authHandler := &handler.AuthHandler{
		Queries: queries,
		JWTSecret: cfg.JWTSecret,
	}

	categoryHandler := &handler.CategoryHandler{
		Queries: queries,
	}

	expensesHandler := &handler.ExpensesHandler{
		Queries: queries,
	}

	auth := middleware.AuthMiddleware(cfg.JWTSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.Handle("GET /categories", auth(http.HandlerFunc(categoryHandler.List)))
	mux.Handle("POST /categories", auth(http.HandlerFunc(categoryHandler.Create)))
	mux.Handle("PUT /categories/{categoryId}", auth(http.HandlerFunc(categoryHandler.Update)))
	mux.Handle("DELETE /categories/{categoryId}", auth(http.HandlerFunc(categoryHandler.Delete)))

	mux.Handle("POST /expenses", auth(http.HandlerFunc(expensesHandler.Create)))
	mux.Handle("GET /expenses", auth(http.HandlerFunc(expensesHandler.List)))
	mux.Handle("GET /expenses/{expensesId}", auth(http.HandlerFunc(expensesHandler.GetExpensesByID)))
	mux.Handle("PUT /expenses/{expensesId}", auth(http.HandlerFunc(expensesHandler.Update)))
	mux.Handle("DELETE /expenses/{expensesId}", auth(http.HandlerFunc(expensesHandler.Delete)))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
