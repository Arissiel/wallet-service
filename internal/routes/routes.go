package routes

import (
	"errors"
	"fmt"
	"wallet-service/internal/api"
	"wallet-service/internal/db"
	"wallet-service/internal/logger"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, repo db.Repository) error {
	fmt.Printf("Repository: %+v\n", repo)
	if repo == nil {
		err := errors.New("repository is nil")
		logger.Log.Error(err)
		return err
	}

	walletHandlers := api.NewWalletHandler(repo)

	api := router.Group("/api/v1")
	{
		// POST запросы для депозита и снятия
		api.POST("/wallet", walletHandlers.PostWalletOperation)

		// GET запрос для получения баланса
		api.GET("/wallets/:walletUUID", walletHandlers.GetBalance)

		//Для корректной и предсказуемой обработки ошибки, когда не указан walletUUID
		api.GET("/wallets", walletHandlers.GetBalance)
	}
	return nil
}
