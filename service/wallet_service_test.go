package service_test

import (
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/mocks"
	"coinbit-wallet/service"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletServiceSuite struct {
	suite.Suite
	walletService      service.IWalletService
	balanceView        *mocks.MockBalanceView
	aboveThresholdView *mocks.MockAboveThresholdView
	depositEmitter     *mocks.MockDepositEmitter
	timestamp          time.Time
}

func TestWalletServiceSuite(t *testing.T) {
	suite.Run(t, new(WalletServiceSuite))
}

func (s *WalletServiceSuite) SetupSuite() {
	s.balanceView = new(mocks.MockBalanceView)
	s.aboveThresholdView = new(mocks.MockAboveThresholdView)
	s.depositEmitter = new(mocks.MockDepositEmitter)
	s.timestamp = time.Now()
	timestampGen := mocks.MockTimestampGenerator{
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

	testCases := []struct {
		name            string
		depositRequest  request.WalletDepositRequest
		expectedDeposit *model.Deposit
		expectedError   error
	}{
		{
			name: "Success deposit request",
			depositRequest: request.WalletDepositRequest{
				WalletId: "111-222",
				Amount:   2000,
			},
			expectedDeposit: &model.Deposit{
				WalletId:    "111-222",
				Amount:      2000,
				DepositedAt: timestamppb.New(s.timestamp),
			},
			expectedError: nil,
		},
		{
			name: "Failed deposit request during emitSync",
			depositRequest: request.WalletDepositRequest{
				WalletId: "111-222",
				Amount:   2000,
			},
			expectedDeposit: nil,
			expectedError:   errors.New("something unexpected happen on EmitSync"),
		},
	}

	for _, tc := range testCases {
		deposit := model.Deposit{
			WalletId:    tc.depositRequest.WalletId,
			Amount:      tc.depositRequest.Amount,
			DepositedAt: timestamppb.New(s.timestamp),
		}

		s.depositEmitter.On("EmitSync", &deposit).Return(tc.expectedError).Once()
		err := s.walletService.DepositWallet(tc.depositRequest)

		if tc.expectedError != nil {
			require.Error(s.T(), err, tc.name)
		} else {
			require.NoError(s.T(), err, tc.name)
		}
	}
}

func (s *WalletServiceSuite) TestWalletService_GetWalletDetails() {

	testCases := []struct {
		name                  string
		walletId              string
		balance               *model.Balance
		aboveThreshold        *model.AboveThreshold
		balanceError          error
		aboveThresholdError   error
		expectedWalletDetails *response.GetWalletDetailsResponse
	}{
		{
			name:     "Success to get wallet details",
			walletId: "111-222",
			balance: &model.Balance{
				WalletId: "111-222",
				Balance:  2000,
			},
			aboveThreshold: &model.AboveThreshold{
				WalletId:            "111-222",
				AmountWithinTwoMins: 2000,
				Status:              false,
				StartPeriod:         timestamppb.New(s.timestamp),
			},
			balanceError:        nil,
			aboveThresholdError: nil,
			expectedWalletDetails: &response.GetWalletDetailsResponse{
				WalletId:       "111-222",
				Balance:        2000,
				AboveThreshold: false,
			},
		},
		{
			name:                  "Failed to get balance from view",
			walletId:              "111-222",
			balance:               nil,
			aboveThreshold:        nil,
			balanceError:          errors.New("unable to get balance from view"),
			aboveThresholdError:   nil,
			expectedWalletDetails: nil,
		},
		{
			name:                  "Failed to get above threshold from view",
			walletId:              "111-222",
			balance:               nil,
			aboveThreshold:        nil,
			balanceError:          nil,
			aboveThresholdError:   errors.New("unable to get aboveThreshold from view"),
			expectedWalletDetails: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.balanceView.On("GetByKey", tc.walletId).Return(tc.balance, tc.balanceError).Once()

			if tc.balanceError == nil {
				s.aboveThresholdView.On("GetByKey", tc.walletId).Return(tc.aboveThreshold, tc.aboveThresholdError).Once()
			}

			walletDetails, err := s.walletService.GetWalletDetails(tc.walletId)

			if tc.expectedWalletDetails == nil {
				require.Error(s.T(), err)
				require.Nil(s.T(), walletDetails)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), walletDetails)
				require.EqualValues(s.T(), tc.expectedWalletDetails.WalletId, walletDetails.WalletId)
				require.EqualValues(s.T(), tc.expectedWalletDetails.Balance, walletDetails.Balance)
				require.EqualValues(s.T(), tc.expectedWalletDetails.AboveThreshold, walletDetails.AboveThreshold)
			}
		})
	}
}
