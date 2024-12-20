package db

import (
	"database/sql"
)

type Repository interface {
	DepositMoney(walletUUID string, amount int64) error
	WithdrawMoney(walletUUID string, amount int64) error
	GetBalance(walletUUID string) (int64, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
