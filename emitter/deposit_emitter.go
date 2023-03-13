package emitter

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"

	"github.com/lovoo/goka"
)

type IDepositEmitter interface {
	EmitSync(deposit *model.Deposit) error
}

type DepositEmitter struct {
	emitter *goka.Emitter
}

func NewDepositEmitter(brokers []string) (*DepositEmitter, error) {
	logger.Info("Creating new deposit emitter..")
	depositEmitter, err := goka.NewEmitter(brokers, config.TopicDeposit, new(util.DepositCodec))
	if err != nil {
		logger.Error("error creating emitter: %v", err)
		return nil, err
	}
	return &DepositEmitter{emitter: depositEmitter}, nil
}

func (e *DepositEmitter) EmitSync(deposit *model.Deposit) error {
	logger.Info("emitting deposit request")
	err := e.emitter.EmitSync(deposit.WalletId, deposit)
	if err != nil {
		logger.Error("error emitting message: %v", err)
		return err
	}
	logger.Info("deposit request emitted")
	return nil
}
