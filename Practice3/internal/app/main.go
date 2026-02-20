package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Practice3/internal/handler"
	"Practice3/internal/middleware"
	"Practice3/internal/repository"
	"Practice3/internal/repository/_postgres"
	"Practice3/internal/usecase"
	"Practice3/pkg/modules"
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

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "112407",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
