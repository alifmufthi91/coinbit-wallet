package service_test

import (
	"coinbit-wallet/dto/app"
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/emitter"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/service"
	"coinbit-wallet/view"
	"errors"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletServiceSuite struct {
	suite.Suite
	walletService      service.IWalletService
	balanceView        *view.MockBalanceView
	aboveThresholdView *view.MockAboveThresholdView
	depositEmitter     *emitter.MockDepositEmitter
	timestamp          *timestamppb.Timestamp
}

func TestWalletServiceSuite(t *testing.T) {
	suite.Run(t, new(WalletServiceSuite))
}

func (s *WalletServiceSuite) SetupSuite() {
	s.balanceView = new(view.MockBalanceView)
	s.aboveThresholdView = new(view.MockAboveThresholdView)
	s.depositEmitter = new(emitter.MockDepositEmitter)
	s.timestamp = timestamppb.Now()
	timestampGen := app.MockTimestampGenerator{
		Timestamp: s.timestamp,
	}
	s.walletService = service.NewWalletService(s.balanceView, s.aboveThresholdView, s.depositEmitter, timestampGen)
}

func (s *WalletServiceSuite) AfterTest(_, _ string) {
	s.balanceView.AssertExpectations(s.T())
	s.aboveThresholdView.AssertExpectations(s.T())
	s.depositEmitter.AssertExpectations(s.T())
}

func (s *WalletServiceSuite) TestWalletService_DepositWallet() {

	walletDepositRequest := request.WalletDepositRequest{
		WalletId: "111-222",
		Amount:   2000,
	}

	deposit := model.Deposit{
		WalletId:    walletDepositRequest.WalletId,
		Amount:      walletDepositRequest.Amount,
		DepositedAt: s.timestamp,
	}

	s.depositEmitter.On("EmitSync", &deposit).Return(nil).Once()
	err := s.walletService.DepositWallet(walletDepositRequest)
	require.NoError(s.T(), err)

	expectedError := errors.New("something unexpected happen on EmitSync")
	s.depositEmitter.On("EmitSync", &deposit).Return(expectedError).Once()
	err = s.walletService.DepositWallet(walletDepositRequest)
	require.Error(s.T(), err)
}

func (s *WalletServiceSuite) TestWalletService_GetWalletDetails() {

	walletId := "111-222"

	balance := model.Balance{
		WalletId: walletId,
		Balance:  2000,
	}

	aboveThreshold := model.AboveThreshold{
		WalletId:            walletId,
		AmountWithinTwoMins: 2000,
		Status:              false,
		StartPeriod:         timestamppb.Now(),
	}

	expectedWalletDetails := response.GetWalletDetailsResponse{
		WalletId:       walletId,
		Balance:        balance.Balance,
		AboveThreshold: aboveThreshold.Status,
	}

	s.balanceView.On("GetByKey", walletId).Return(&balance, nil).Once()
	s.aboveThresholdView.On("GetByKey", walletId).Return(&aboveThreshold, nil).Once()
	walletDetails, err := s.walletService.GetWalletDetails(walletId)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(expectedWalletDetails, *walletDetails))

	expectedError := errors.New("unable to get balance from view")
	var emptyBalance *model.Balance
	s.balanceView.On("GetByKey", walletId).Return(emptyBalance, expectedError).Once()
	walletDetails, err = s.walletService.GetWalletDetails(walletId)

	require.Error(s.T(), err)
	require.Equal(s.T(), expectedError.Error(), err.Error())
	require.Nil(s.T(), walletDetails)

	expectedError2 := errors.New("unable to get aboveThreshold from view")
	var emptyAboveThreshold *model.AboveThreshold
	s.balanceView.On("GetByKey", walletId).Return(&balance, nil).Once()
	s.aboveThresholdView.On("GetByKey", walletId).Return(emptyAboveThreshold, expectedError2).Once()
	walletDetails, err = s.walletService.GetWalletDetails(walletId)

	require.Error(s.T(), err)
	require.Equal(s.T(), expectedError2.Error(), err.Error())
	require.Nil(s.T(), walletDetails)
}
