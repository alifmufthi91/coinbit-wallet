package service

import (
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/emitter"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/view"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type IWalletService interface {
	DepositWallet(request request.WalletDepositRequest) error
	GetWalletDetails(walletId string) (*response.GetWalletDetailsResponse, error)
}

type WalletService struct {
	balanceView        *view.BalanceView
	aboveThresholdView *view.AboveThresholdView
}

func NewWalletService(bv *view.BalanceView, atv *view.AboveThresholdView) WalletService {
	return WalletService{
		balanceView:        bv,
		aboveThresholdView: atv,
	}
}

func (ws WalletService) DepositWallet(request request.WalletDepositRequest) error {
	logger.Info("Deposit to wallet, data = %v", request)
	deposit := &model.Deposit{
		WalletId:    request.WalletId,
		Amount:      request.Amount,
		DepositedAt: timestamppb.Now(),
	}

	if err := emitter.EmitDeposit(deposit); err != nil {
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
