package main

import (
	"Practice4/internal/repository"
	"Practice4/internal/repository/_postgres"
	"Practice4/pkg/modules"
	"context"
	"fmt"
	"time"
)

func main() {
	Run()
}

func Run(){
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
dbConfig := initPostgreConfig()
_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
repositories := repository.NewRepositories(_postgre)
users, err := repositories.GetUsers()
if err != nil {
fmt.Printf("Error fetching users: %v\n", err)
return
}
fmt.Printf("Users: %+v\n", users)
}
func initPostgreConfig() *modules.PostgreConfig {
return &modules.PostgreConfig{
Host: "localhost",
Port: "5432",
Username: "postgres",
Password: "112407",
DBName: "mydb",
SSLMode: "disable",
ExecTimeout: 5 * time.Second,
}
}