package _postgres

import (
	"Practice4/pkg/modules"
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)
type Dialect struct {
	DB *sqlx.DB
}
func NewPGXDialect(ctx context.Context, cfg *modules.PostgreConfig) *Dialect {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var db *sqlx.DB
	var err error

	
	for i := 0; i < 10; i++ {
		fmt.Println("⏳ Waiting for database...")

		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Println("✅ Connected to database!")
				break
			}
		}

		fmt.Println("❌ Database not ready, retrying...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("❌ Could not connect to database after retries: " + err.Error())
	}

	AutoMigrate(cfg)

	return &Dialect{
		DB: db,
	}
}
func AutoMigrate(cfg *modules.PostgreConfig) {
	sourceURL := "file://database/migrations"
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,cfg.SSLMode)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}