package db

import (
	"database/sql"
	"errors"
	"fmt"
	"wallet-service/internal/logger"
)

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func (r *PostgresRepository) DepositMoney(walletUUID string, amount int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Errorf("Failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var committed bool

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logger.Log.Errorf("Transaction panicked: %v", p)
		} else if err != nil && !committed {
			tx.Rollback()
			logger.Log.Errorf("Transaction rolled back: %v", err)
		}
	}()

	// Блокируем строку кошелька
	var balance int64
	var walletID int

	logger.Log.Debugf("Executing query: %s with params: %v", QueryGetWalletForUpdate, walletUUID)
	err = tx.QueryRow(QueryGetWalletForUpdate, walletUUID).Scan(&balance)
	if err == sql.ErrNoRows {
		// Если кошелька нет, создаем новый
		logger.Log.Infof("Wallet with UUID %s not found. Creating a new wallet.", walletUUID)

		if _, err = tx.Exec(QueryCreateWallet, walletUUID, amount); err != nil {
			logger.Log.Errorf("Failed to create wallet with UUID %s: %v", walletUUID, err)
			return fmt.Errorf("failed to create wallet: %w", err)
		}

		logger.Log.Debugf("Executing query: %s with params: %v", QueryGetWalletID, walletUUID)
		if err = tx.QueryRow(QueryGetWalletID, walletUUID).Scan(&walletID); err != nil {
			logger.Log.Errorf("Failed to retrieve wallet ID for UUID %s: %v", walletUUID, err)
			return fmt.Errorf("failed to get wallet ID: %w", err)
		}

		if _, err = tx.Exec(QueryCreateTransaction, walletID, "DEPOSIT", amount, "ACTIVE"); err != nil {
			logger.Log.Errorf("Failed to create transaction for wallet UUID %s: %v", walletUUID, err)
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		logger.Log.Infof("Wallet with UUID %s created successfully with initial deposit of %d.", walletUUID, amount)
	} else if err != nil {
		logger.Log.Errorf("Failed to lock wallet with UUID %s for update: %v", walletUUID, err)
		return fmt.Errorf("failed to lock wallet for update: %w", err)
	} else {
		// Обновляем баланс, если кошелек существует
		logger.Log.Infof("Wallet with UUID %s exists. Depositing amount: %d.", walletUUID, amount)
		if _, err = tx.Exec(QueryUpdateBalance, amount, walletUUID); err != nil {
			logger.Log.Errorf("Failed to deposit money to wallet UUID %s: %v", walletUUID, err)
			return fmt.Errorf("failed to deposit money: %w", err)
		}

		// Получаем wallet_id для создания транзакции
		logger.Log.Debugf("Executing query: %s with params: %v", QueryGetWalletID, walletUUID)
		if err = tx.QueryRow(QueryGetWalletID, walletUUID).Scan(&walletID); err != nil {
			logger.Log.Errorf("Failed to retrieve wallet ID for UUID %s: %v", walletUUID, err)
			return fmt.Errorf("failed to get wallet ID: %w", err)
		}

		//Создаем транзакцию
		if _, err = tx.Exec(QueryCreateTransaction, walletID, "DEPOSIT", amount, "ACTIVE"); err != nil {
			logger.Log.Errorf("Failed to create transaction for wallet UUID %s: %v", walletUUID, err)
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		logger.Log.Infof("Deposit of %d to wallet UUID %s completed successfully.", amount, walletUUID)
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Errorf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	logger.Log.Info("Transaction committed successfully")
	return nil
}

func (r *PostgresRepository) WithdrawMoney(walletUUID string, amount int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Errorf("Failed to start transaction: %v", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var committed bool

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			logger.Log.Errorf("Transaction panicked: %v", p)
		} else if !committed {
			_ = tx.Rollback()
			logger.Log.Errorf("Transaction rolled back: %v", err)
		}
	}()

	var balance int64
	var walletID int

	logger.Log.Debugf("Executing query: %s with params: %v", QueryGetWalletForUpdate, walletUUID)
	err = tx.QueryRow(QueryGetWalletForUpdate, walletUUID).Scan(&balance)
	if err == sql.ErrNoRows {
		logger.Log.Error(ErrWalletNotFound)
		return ErrWalletNotFound
	} else if err != nil {
		logger.Log.Errorf("Failed to lock wallet with UUID %s for update: %v", walletUUID, err)
		return fmt.Errorf("failed to lock wallet for update: %w", err)
	}

	if balance < amount {
		logger.Log.Error(ErrInsufficientFunds)
		return ErrInsufficientFunds
	}

	if _, err = tx.Exec(QueryWithdraw, amount, walletUUID); err != nil {
		logger.Log.Errorf("Failed to withdraw money from wallet UUID %s: %v", walletUUID, err)
		return fmt.Errorf("failed to withdraw money: %w", err)
	}

	logger.Log.Debugf("Executing query: %s with params: %v", QueryGetWalletID, walletUUID)
	if err = tx.QueryRow(QueryGetWalletID, walletUUID).Scan(&walletID); err != nil {
		logger.Log.Errorf("Failed to get wallet ID for UUID %s: %v", walletUUID, err)
		return fmt.Errorf("failed to get wallet ID: %w", err)
	}

	if _, err = tx.Exec(QueryCreateTransaction, walletID, "WITHDRAW", amount, "ACTIVE"); err != nil {
		logger.Log.Errorf("Failed to create transaction for wallet UUID %s: %v", walletUUID, err)
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Errorf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	logger.Log.Info("Transaction committed successfully")
	return nil
}

func (r *PostgresRepository) GetBalance(walletUUID string) (int64, error) {
	var balance int64

	logger.Log.Infof("Fetching balance for wallet UUID: %s", walletUUID)

	// Выполняем запрос для получения баланса
	if err := r.db.QueryRow(QueryGetBalance, walletUUID).Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Warnf("Wallet with UUID %s not found.", walletUUID)
			return 0, ErrWalletNotFound
		}
		logger.Log.Errorf("Failed to fetch balance for wallet UUID %s: %v", walletUUID, err)
		return 0, fmt.Errorf("failed to get wallet balance: %w", err)
	}

	logger.Log.Infof("Successfully retrieved balance for wallet UUID %s: %d", walletUUID, balance)
	return balance, nil

}
