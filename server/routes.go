package server

import (
	"coinbit-wallet/controller"
	"coinbit-wallet/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lovoo/goka"
)

func NewRouter(bv *goka.View, atv *goka.View) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())

	router.Use(middleware.ErrorHandlingMiddleware())

	v1 := router.Group("api/v1")
	{
		wallet := v1.Group("wallet")
		walletController := controller.NewWalletController(bv, atv)
		{
			wallet.POST("/deposit", walletController.Deposit)
			wallet.GET("/details/:walletId", walletController.GetDetails)
		}
	}

	return router
}
