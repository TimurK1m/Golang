package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"Practice4/internal/handler"
	"Practice4/internal/middleware"
	"Practice4/internal/repository"
	"Practice4/internal/repository/_postgres"
	"Practice4/internal/usecase"
	"Practice4/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// DB config
	dbConfig := initPostgreConfig()

	// connect to postgres
	db := _postgres.NewPGXDialect(ctx, dbConfig)

	// repositories
	repos := repository.NewRepositories(db)

	// usecase
	userUsecase := usecase.NewUserUsecase(repos.UserRepository)

	// handler
	userHandler := handler.NewUserHandler(userUsecase)

	// mux
	mux := http.NewServeMux()

	// middleware chain
	handlerWithMiddleware := middleware.Logger(
		middleware.APIKey(userHandler),
	)

	mux.Handle("/users", handlerWithMiddleware)
	mux.Handle("/users/", handlerWithMiddleware)

	// healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("HEALTH ROUTE REGISTERED")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8000"
    }

	log.Printf("Server started on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "112407"),
		DBName:   getEnv("DB_NAME", "mydb"),
		SSLMode:  "disable",
		ExecTimeout: 5 * time.Second,
	}
}