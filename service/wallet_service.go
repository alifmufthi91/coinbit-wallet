package service

import (
	"coinbit-wallet/dto/app"
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/emitter"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/view"
)

type IWalletService interface {
	DepositWallet(request request.WalletDepositRequest) error
	GetWalletDetails(walletId string) (*response.GetWalletDetailsResponse, error)
}

type WalletService struct {
	balanceView        view.IBalanceView
	aboveThresholdView view.IAboveThresholdView
	depositEmitter     emitter.IDepositEmitter
	timestampGen       app.ITimestampGenerator
}

func NewWalletService(bv view.IBalanceView, atv view.IAboveThresholdView, e emitter.IDepositEmitter, tg app.ITimestampGenerator) WalletService {
	return WalletService{
		balanceView:        bv,
		aboveThresholdView: atv,
		depositEmitter:     e,
		timestampGen:       tg,
	}
}

func (ws WalletService) DepositWallet(request request.WalletDepositRequest) error {
	logger.Info("Deposit to wallet, data = %v", request)
	deposit := &model.Deposit{
		WalletId:    request.WalletId,
		Amount:      request.Amount,
		DepositedAt: ws.timestampGen.Generate(),
	}

	if err := ws.depositEmitter.EmitSync(deposit); err != nil {
		return err
	}

	return nil
}

func (ws WalletService) GetWalletDetails(walletId string) (*response.GetWalletDetailsResponse, error) {
	logger.Info("Get wallet details by walletId, walletId = %s", walletId)
	balance, err := ws.balanceView.GetByKey(walletId)
	if err != nil {
		return nil, err
	}

	aboveThreshold, err := ws.aboveThresholdView.GetByKey(walletId)
	if err != nil {
		return nil, err
	}

	return &response.GetWalletDetailsResponse{
		WalletId:       walletId,
		Balance:        balance.GetBalance(),
		AboveThreshold: aboveThreshold.GetStatus(),
	}, nil
}
