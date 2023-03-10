package controller

import (
	"coinbit-wallet/dto/request"
	"coinbit-wallet/service"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/util/response_util"

	"github.com/gin-gonic/gin"
)

type IWalletController interface {
	Deposit(c *gin.Context)
	GetDetails(c *gin.Context)
}

type walletController struct {
	walletService service.IWalletService
}

func NewWalletController(ws service.IWalletService) IWalletController {
	logger.Info("Initializing wallet controller..")
	return walletController{
		walletService: ws,
	}
}

func (wc walletController) Deposit(c *gin.Context) {
	logger.Info("deposit to wallet request")

	var body request.WalletDepositRequest
	err := c.ShouldBind(&body)
	if err != nil {
		panic(err)
	}

	err = wc.walletService.DepositWallet(body)
	if err != nil {
		panic(err)
	}
	response_util.Ok(c, nil, false)
	logger.Info("deposit to wallet success")
}

func (wc walletController) GetDetails(c *gin.Context) {
	logger.Info("Get details wallet request")
	walletId := c.Param("walletId")

	response, err := wc.walletService.GetWalletDetails(walletId)
	if err != nil {
		panic(err)
	}

	response_util.Ok(c, response, false)
	logger.Info("get wallet details success")
}
