package server

import (
	"coinbit-wallet/controller"
	"coinbit-wallet/dto/app"
	"coinbit-wallet/emitter"
	"coinbit-wallet/middleware"
	"coinbit-wallet/service"
	"coinbit-wallet/view"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(bv view.IBalanceView, atv view.IAboveThresholdView, de emitter.IDepositEmitter) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.Use(middleware.ErrorHandlingMiddleware())

	timestampGen := app.NewTimeStampGenerator()
	walletService := service.NewWalletService(bv, atv, de, timestampGen)

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
