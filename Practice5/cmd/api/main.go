package main

import (
	"database/sql"
	"log"
	"net/http"

	"Practice5/internal/handlers"
	"Practice5/internal/repository"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "user=postgres password=112407 dbname=practice5 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)

	handler := &handlers.Handler{
		Repo: repo,
	}

	http.HandleFunc("/users", handler.GetUsersHandler)
	http.HandleFunc("/users/common-friends", handler.GetCommonFriendsHandler)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}