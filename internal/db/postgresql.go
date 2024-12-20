package db

import (
	"database/sql"
	"fmt"
	"wallet-service/internal/logger"
)

func InitDB(dbHost, dbPort, dbUser, dbPassword, dbName string) (*sql.DB, error) {

	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Log.Errorf("Failed to connect to data base: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		logger.Log.Errorf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	logger.Log.Println("Successfully connected to the database")
	return db, nil

}
