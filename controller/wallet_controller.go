package controller

import (
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/emitter"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util/logger"
	responseUtil "coinbit-wallet/util/response"
	"coinbit-wallet/view"
	"encoding/json"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IWalletController interface {
	Deposit(c *gin.Context)
	GetDetails(c *gin.Context)
}

type walletController struct {
}

func NewWalletController() IWalletController {
	logger.Info("Initializing wallet controller..")
	return walletController{}
}

func (wc walletController) Deposit(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)
	logger.Info("deposit to wallet request")

	var body request.WalletDepositRequest
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	v := validator.New()
	err = v.Struct(body)
	if err != nil {
		panic(err)
	}

	deposit := &model.Deposit{
		WalletId:    body.WalletId,
		Amount:      body.Amount,
		DepositedAt: timestamppb.Now(),
	}

	go emitter.EmitDeposit(deposit)

	responseUtil.Success(c, nil, false)
	logger.Info("deposit to wallet success")
}

func (wc walletController) GetDetails(c *gin.Context) {
	defer responseUtil.ErrorHandling(c)
	logger.Info("Get details wallet request")
	walletId := c.Param("walletId")
	wg := sync.WaitGroup{}

	var aboveThresholdMap *model.AboveThresholdMap
	var balanceMap *model.BalanceMap

	wg.Add(2)
	go func() {
		aboveThresholdMap = view.GetAboveThresholdView()
		wg.Done()
	}()
	go func() {
		balanceMap = view.GetBalanceView()
		wg.Done()
	}()

	var balance float32
	var aboveThreshold bool
	wg.Wait()
	if val, ok := aboveThresholdMap.Items[walletId]; ok {
		if timestamppb.Now().Seconds > val.StartPeriod.Seconds+120 {
			aboveThreshold = false
		} else {
			aboveThreshold = val.GetStatus()
		}
	}
	if val, ok := balanceMap.Items[walletId]; ok {
		balance = val.GetBalance()
	}
	response := response.GetWalletDetailsResponse{
		Balance:        balance,
		AboveThreshold: aboveThreshold,
	}

	responseUtil.Success(c, response, false)
	logger.Info("get wallet details success")
}
