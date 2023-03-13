package mocks

import (
	"coinbit-wallet/generated/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockBalanceView struct {
	mock.Mock
}

func (m *MockBalanceView) Run(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBalanceView) GetByKey(key string) (*model.Balance, error) {
	args := m.Called(key)
	return args.Get(0).(*model.Balance), args.Error(1)
}
