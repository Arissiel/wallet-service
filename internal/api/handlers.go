package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"wallet-service/internal/db"
	"wallet-service/internal/logger"
)

type WalletHandlersInterface interface {
	PostWalletOperation(c *gin.Context)
	GetBalance(c *gin.Context)
}

type WalletHandlers struct {
	Repo db.Repository
}

func NewWalletHandler(repo db.Repository) *WalletHandlers {
	return &WalletHandlers{Repo: repo}
}

func (h *WalletHandlers) PostWalletOperation(c *gin.Context) {
	logger.Log.Debugf("Entering handler PostWalletOperation")
	defer logger.Log.Debugf("Exiting handler PostWalletOperation")
	//структура запроса
	var req struct {
		WalletUUID    string `json:"walletId" binding:"required,uuid"`
		OperationType string `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
		Amount        int64  `json:"amount" binding:"required,gt=0"`
	}

	//привязываем JSON запрос к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	logger.Log.Infof("Processing operation %s for wallet %s with amount %d", req.OperationType, req.WalletUUID, req.Amount)

	switch req.OperationType {
	case "DEPOSIT":
		//пополнение кошелька
		err := h.Repo.DepositMoney(req.WalletUUID, req.Amount)
		if err != nil {
			logger.Log.Errorf("Failed to deposit money for wallet %s: %v", req.WalletUUID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deposit money"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})

	case "WITHDRAW":
		//Вывод средств
		err := h.Repo.WithdrawMoney(req.WalletUUID, req.Amount)
		if err != nil {
			//Обработка ошибок в зависимости от их типа
			if errors.Is(err, db.ErrInsufficientFunds) || errors.Is(err, db.ErrWalletNotFound) {
				logger.Log.Warnf("Withdraw failed for wallet %s: %v", req.WalletUUID, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				logger.Log.Errorf("Failed to withdraw money for wallet %s: %v", req.WalletUUID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Withdraw successful"})

	default:
		logger.Log.Warnf("Invalid operation type: %s", req.OperationType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
	}
}

func (h *WalletHandlers) GetBalance(c *gin.Context) {
	//получение UUID кошелька из параметра запроса
	walletUUID := c.Param("walletUUID")
	if walletUUID == "" {
		logger.Log.Warn("Missing walletUUID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing walletUUID parameter"})
		return
	}

	logger.Log.Infof("Fetching balance for wallet %s", walletUUID)

	//получение баланса из репозитория
	balance, err := h.Repo.GetBalance(walletUUID)
	if err != nil {
		if errors.Is(err, db.ErrWalletNotFound) {
			logger.Log.Warnf("Wallet %s not found", walletUUID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		} else {
			logger.Log.Errorf("Failed to fetch balance for wallet %s: %v", walletUUID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
		}
		return
	}

	logger.Log.Infof("Successfully retrieved balance for wallet %s: %d", walletUUID, balance)
	c.JSON(http.StatusOK, gin.H{"walletId": walletUUID, "balance": balance})
}
