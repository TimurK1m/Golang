package app

import (
	"Practice4/internal/repository/_postgres"
	"Practice4/pkg/modules"
	"context"
	"fmt"
	"time"
)
func Run(){
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
dbConfig := initPostgreConfig()
_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
fmt.Println(_postgre)
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