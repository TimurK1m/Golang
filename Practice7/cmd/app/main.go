package main

import (
	v1 "Secure/internal/controller/htts/v1"
	"Secure/internal/entity"
	"Secure/internal/usecase"
	"Secure/internal/usecase/repo"
	"Secure/pkg/postgres"
	"log"

	postgresDriver "gorm.io/driver/postgres"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
		// 🔌 Подключение к БД
		dsn := "host=localhost user=postgres password=112407 dbname=testdb port=5432 sslmode=disable"

		 // 🔥 ВКЛЮЧИЛ

		

		db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})
		if err != nil {
			log.Fatal("failed to connect database:", err)
		}
		err = db.AutoMigrate(&entity.User{})
		if err != nil {
			log.Fatal("migration failed:", err)
		}
		pg := &postgres.Postgres{Conn: db}

		// 🧱 Слои
		userRepo := repo.NewUserRepo(pg)
		userUseCase := usecase.NewUserUseCase(userRepo)

		// 🚀 Gin
		r := gin.Default()

	//r.Use(utils.RateLimiter()) // 🔥 ВКЛЮЧИЛ
		api := r.Group("/api")
		v1.NewUserRoutes(api, userUseCase, db.Logger)

		// ▶️ Запуск
		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}