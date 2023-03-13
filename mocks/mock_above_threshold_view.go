package mocks

import (
	"coinbit-wallet/generated/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAboveThresholdView struct {
	mock.Mock
}

func (m *MockAboveThresholdView) Run(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAboveThresholdView) GetByKey(key string) (*model.AboveThreshold, error) {
	args := m.Called(key)
	return args.Get(0).(*model.AboveThreshold), args.Error(1)
}
