-- Кошельки
CREATE TABLE wallets (
    wallet_id SERIAL PRIMARY KEY,                          -- Автоинкрементируемый ID 
    uuid UUID NOT NULL UNIQUE,                             -- Приходит с запросом
    balance BIGINT NOT NULL DEFAULT 0,                     -- Баланс кошелька
    deleted_at TIMESTAMP NULL,                             -- Мягкое удаление (NULL - активный кошелёк)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),           -- Дата создания кошелька
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()            -- Дата последнего обновления
);

--Транзакции
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,                                 -- Автоинкрементируемый ID 
    wallet_id INT NOT NULL,                                -- Связь с кошельком    
    operation_type VARCHAR(10) NOT NULL,                   -- Тип операции DEPOSIT или WITHDRAW
    amount BIGINT NOT NULL CHECK (amount > 0),             -- Сумма операции (должна быть > 0)
    wallet_status VARCHAR(10) NOT NULL DEFAULT 'ACTIVE',   -- Статус кошелька во время операции    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),           -- Дата создания транзакции

 -- Внешний ключ на wallets с запретом удаления, если есть связанные транзакции 
    CONSTRAINT fk_wallet
        FOREIGN KEY (wallet_id)
        REFERENCES wallets(wallet_id)
        ON DELETE NO ACTION
);

--Индексы для ускорения аналитических запросов

-- Индекс для ускорения поиска активных кошельков
CREATE INDEX idx_wallets_deleted_at ON wallets (deleted_at);

-- Индекс для ускорения поиска транзакций по кошельку
CREATE INDEX idx_transactions_wallet_id ON transactions (wallet_id);

-- Индекс для ускорения поиска транзакций по дате создания
CREATE INDEX idx_transactions_created_at ON transactions (created_at);
