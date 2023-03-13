package mocks

import (
	"coinbit-wallet/dto/request"
	"coinbit-wallet/dto/response"

	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) DepositWallet(request request.WalletDepositRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockWalletService) GetWalletDetails(walletId string) (*response.GetWalletDetailsResponse, error) {
	args := m.Called(walletId)
	return args.Get(0).(*response.GetWalletDetailsResponse), args.Error(1)
}
