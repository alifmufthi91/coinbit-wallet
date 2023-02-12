package server

import (
	"coinbit-wallet/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())

	v1 := router.Group("api/v1")
	{
		wallet := v1.Group("wallet")
		walletController := controller.NewWalletController()
		{
			wallet.POST("/deposit", walletController.Deposit)
			wallet.GET("/details/:walletId", walletController.GetDetails)
		}
	}

	return router
}
