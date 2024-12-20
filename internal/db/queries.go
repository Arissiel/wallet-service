package db

const (
	//создание кошелька
	QueryCreateWallet = `
		INSERT INTO wallets (uuid, balance) 
		VALUES ($1, $2)
	`

	//проверка существует ли кошелек по uuid
	QueryDoesWalletExist = `
		SELECT EXISTS (
			SELECT 1 
			FROM wallets 
			WHERE uuid = $1 AND deleted_at IS NULL
		)
	`

	//получение баланса с блокировкой строки
	QueryGetWalletForUpdate = `
		SELECT balance 
		FROM wallets 
		WHERE uuid = $1 AND deleted_at IS NULL
		FOR UPDATE
	`

	//обновление баланса
	QueryUpdateBalance = `
		UPDATE wallets 
		SET balance = balance + $1, updated_at = NOW() 
		WHERE uuid = $2 AND deleted_at IS NULL
	`

	//получение баланса
	QueryGetBalance = `
		SELECT balance 
		FROM wallets 
		WHERE uuid = $1 AND deleted_at IS NULL
	`

	//снятие средств
	QueryWithdraw = `
		UPDATE wallets 
		SET balance = balance - $1, updated_at = NOW() 
		WHERE uuid = $2 AND deleted_at IS NULL
	`

	//получение ID кошелька по UUID (для транзакций)
	QueryGetWalletID = `
	SELECT wallet_id 
	FROM wallets 
	WHERE uuid = $1 AND deleted_at IS NULL
	`

	//создание записи транзакции
	QueryCreateTransaction = `
		INSERT INTO transactions (wallet_id, operation_type, amount, wallet_status) 
		VALUES ($1, $2, $3, $4)
	`
)
