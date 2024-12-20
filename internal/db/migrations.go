package db

import (
	"database/sql"
	"wallet-service/internal/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(db *sql.DB) error {
	// Экземпляр драйвера для работы с базой
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Fatalf("Failed to create a driver for migrations: %v", err)
		return err
	}

	// Создаем объект миграций
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/internal/db/migrations",
		"wallet_db",
		driver,
	)
	if err != nil {
		logger.Log.Fatalf("Failed to create migrations: %v", err)
		return err
	}

	// Выполняем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	logger.Log.Println("Migrations applied successfully")
	return nil
}
