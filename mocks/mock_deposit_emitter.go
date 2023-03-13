package mocks

import (
	"coinbit-wallet/generated/model"

	"github.com/stretchr/testify/mock"
)

type MockDepositEmitter struct {
	mock.Mock
}

func (m *MockDepositEmitter) EmitSync(deposit *model.Deposit) error {
	args := m.Called(deposit)
	return args.Error(0)
}
