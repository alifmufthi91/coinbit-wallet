package server

import (
	"coinbit-wallet/controller"
	"coinbit-wallet/middleware"
	"coinbit-wallet/service"
	"coinbit-wallet/view"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(bv *view.BalanceView, atv *view.AboveThresholdView) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.Use(middleware.ErrorHandlingMiddleware())

	walletService := service.NewWalletService(bv, atv)

	v1 := router.Group("api/v1")
	{
		wallet := v1.Group("wallet")
		walletController := controller.NewWalletController(walletService)
		{
			wallet.POST("/deposit", walletController.Deposit)
			wallet.GET("/details/:walletId", walletController.GetDetails)
		}
	}

	return router
}
