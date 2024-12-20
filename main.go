package main

import (
	"wallet-service/config"
	"wallet-service/internal/db"
	"wallet-service/internal/logger"
	"wallet-service/internal/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Инициализация логгера
	logger.InitLogger()

	logger.Log.Infof("Logger initialized")

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatalf("Failed to load config: %v", err)
	}

	//Инициализация базы данных
	dataBase, err := db.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dataBase.Close()

	// Выполнение миграций
	if err := db.RunMigrations(dataBase); err != nil {
		logger.Log.Fatalf("Failed to apply migrations: %v", err)
	}

	//экземпляр репозитория
	repo := db.NewPostgresRepository(dataBase)

	//инициализация маршрутов
	router := gin.Default()
	if err := routes.SetupRoutes(router, repo); err != nil {
		logger.Log.Fatalf("Failed to set up routes: %v", err)
	}

	//запуск сервера
	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.Log.Fatalf("Failed to start server: %v", err)
	}

}
