package controller

import (
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/emitter"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/util/response_util"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lovoo/goka"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IWalletController interface {
	Deposit(c *gin.Context)
	GetDetails(c *gin.Context)
}

type walletController struct {
	balanceView        *goka.View
	aboveThresholdView *goka.View
}

func NewWalletController(bv *goka.View, atv *goka.View) IWalletController {
	logger.Info("Initializing wallet controller..")
	return walletController{
		balanceView:        bv,
		aboveThresholdView: atv,
	}
}

func (wc walletController) Deposit(c *gin.Context) {
	logger.Info("deposit to wallet request")

	var body request.WalletDepositRequest
	err := c.ShouldBind(&body)
	if err != nil {
		panic(err)
	}

	deposit := &model.Deposit{
		WalletId:    body.WalletId,
		Amount:      body.Amount,
		DepositedAt: timestamppb.Now(),
	}

	if err = emitter.EmitDeposit(deposit); err != nil {
		panic(err)
	}

	response_util.Ok(c, nil, false)
	logger.Info("deposit to wallet success")
}

func (wc walletController) GetDetails(c *gin.Context) {
	logger.Info("Get details wallet request")
	walletId := c.Param("walletId")
	wg := sync.WaitGroup{}

	var aboveThreshold *model.AboveThreshold
	var balance *model.Balance

	wg.Add(2)
	go func() {
		err := util.GetView(wc.aboveThresholdView, walletId, &aboveThreshold)
		if err != nil {
			panic(err)
		}
		if aboveThreshold == nil {
			aboveThreshold = &model.AboveThreshold{}
		}
		wg.Done()
	}()
	go func() {
		err := util.GetView(wc.balanceView, walletId, &balance)
		if err != nil {
			panic(err)
		}
		if balance == nil {
			balance = &model.Balance{}
		}
		wg.Done()
	}()
	wg.Wait()

	isAboveThreshold := false
	if aboveThreshold.StartPeriod != nil {
		if !util.IsWithinTwoMins(aboveThreshold.StartPeriod, timestamppb.Now()) {
			isAboveThreshold = false
		} else {
			isAboveThreshold = aboveThreshold.GetStatus()
		}
	}
	response := response.GetWalletDetailsResponse{
		WalletId:       walletId,
		Balance:        balance.GetBalance(),
		AboveThreshold: isAboveThreshold,
	}

	response_util.Ok(c, response, false)
	logger.Info("get wallet details success")
}
