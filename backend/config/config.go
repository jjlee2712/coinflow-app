package config

import "os"

type Config struct {
	Port string
	DBUrl string
}

func Load() Config {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == ""{
		dbUrl = "postgres://coinflow:coinflow@localhost:5432/coinflow?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{Port: port, DBUrl: dbUrl}
}