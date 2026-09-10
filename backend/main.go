package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jjlee2712/coinflow/backend/config"
	"github.com/jjlee2712/coinflow/backend/handler"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}